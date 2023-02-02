package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/ninroot/gocker/config"
	"github.com/ninroot/gocker/pkg/cgroups"
	"github.com/ninroot/gocker/pkg/container"
	"github.com/ninroot/gocker/pkg/image"
	"github.com/ninroot/gocker/pkg/storage"
	"github.com/ninroot/gocker/pkg/util"
	"github.com/sirupsen/logrus"
)

type RunRequest struct {
	ImageName        string
	ImageTag         string
	ContainerName    string
	ContainerCommand string
	ContainerID      string
	ContainerArgs    []string
	ContainerLimits
}

type ContainerLimits struct {
	MemoryLimit int
	PidsLimit   int
}

type runtimeService struct {
	imgStore storage.ImageStore
	conStore storage.ContainerStore
	cgroup   cgroups.CGroup
}

func NewRuntimeService() *runtimeService {
	return &runtimeService{
		imgStore: storage.NewImageStore(util.EnsureDir(config.DefaultImageStoreRootDir), storage.Btrfs{}),
		conStore: *storage.NewContainerStore(util.EnsureDir(config.DefaultContainerStoreRootDir), storage.Btrfs{}),
		cgroup:   cgroups.New(util.EnsureDir(config.DefaultCGroupDir)),
	}
}

func (r runtimeService) Run(req RunRequest) error {
	req.ContainerID = container.RandID()
	args := append([]string{"internal"},
		"--ContainerName", req.ContainerName,
		"--ImageName", req.ImageName,
		"--ImageTag", req.ImageTag,
		"--ContainerCommand", req.ContainerCommand,
		"--ContainerID", req.ContainerID,
		"--", strings.Join(req.ContainerArgs, " "),
	)
	cmd := exec.Command("/proc/self/exe", args...)

	logrus.WithField("args", args).Debug("Internal Run")

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS,
		// make the mounting point no longer visible to the host
		Unshareflags: syscall.CLONE_NEWNS,
	}

	// cmd.SysProcAttr = &syscall.SysProcAttr{
	// 	Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWUSER | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET,
	// }

	// if not set, we end up having uid=65534(nobody) gid=65534(nogroup) groups=65534(nogroup)
	cmd.SysProcAttr.Credential = &syscall.Credential{
		Uid: uint32(0),
		Gid: uint32(0),
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Failed to start container: %s", err)
	}

	g := r.cgroup.NewGroup(req.ContainerID)
	defer func() {
		if err := g.Delete(); err != nil {
			logrus.Fatal(err)
		}
	}()

	err := applyCGroup(g, cmd.Process.Pid, req.ContainerLimits)
	if err != nil {
		return fmt.Errorf("Failed to apply cgroup: %s", err)
	}

	return cmd.Wait()
}

func applyCGroup(g cgroups.Group, pid int, l ContainerLimits) error {
	logrus.WithField("pid", pid).Debug("Set cgroup to process")

	if err := g.SetPidMax(l.PidsLimit); err != nil {
		return fmt.Errorf("Failed to set pids limit: %s", err)
	}
	if err := g.SetMemoryLimit(l.MemoryLimit); err != nil {
		return fmt.Errorf("Failed to memory limit: %s", err)
	}
	if err := g.SetNotifyOnRelease(true); err != nil {
		return fmt.Errorf("Failed set notify on release: %s", err)
	}
	return g.AddProc(pid)
}

func (r runtimeService) InitContainer(req RunRequest) error {
	img, err := r.FindImageByNameAndId(req.ImageName, req.ImageTag)
	if err != nil {
		return err
	}
	if img == nil {
		return fmt.Errorf("image not found: %s:%s", req.ImageName, req.ImageTag)
	}

	imgH := r.imgStore.GetImage(img.Digest)
	contH, err := r.conStore.CreateContainer(req.ContainerID, imgH.ImageDir())
	if err != nil {
		return nil
	}

	c := container.Container{
		ID:        req.ContainerID,
		Name:      req.ContainerName,
		Image:     *img,
		CreatedAt: time.Now(),
		Command:   req.ContainerCommand,
		Args:      req.ContainerArgs,
	}

	contH.SetSpec(c)

	unbind := bindDevices(contH.RootfsDir())
	defer unbind()

	if err := syscall.Chroot(contH.RootfsDir()); err != nil {
		return err
	}
	if err := syscall.Chdir("/"); err != nil {
		return err
	}

	// hostname will be affected if this function runs in a process that hasn't been with CLONE_NEWUTS
	// happens typically when debugging
	syscall.Sethostname([]byte(req.ContainerID))

	// mount /proc to make commands such `ps` working
	if err := syscall.Mount("proc", "proc", "proc", 0, ""); err != nil {
		return err
	}
	defer syscall.Unmount("/proc", 0)

	cmd := exec.Command(c.Command, c.Args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func bindDevices(rootDir string) func() {
	devices := []string{
		"/dev/zero",
		"/dev/null",
	}
	unbind := []func(){}

	for _, d := range devices {
		u, err := bindDevice(d, filepath.Join(rootDir, d))
		if err != nil {
			logrus.WithField("device", d).WithError(err).Warn("Failed to bind device")
		}
		if u != nil {
			unbind = append(unbind, u)
		}
	}

	return func() {
		for _, u := range unbind {
			u()
		}
	}
}

func bindDevice(source, target string) (unmount func(), err error) {
	f, err := os.Create(target)
	if err != nil {
		return nil, fmt.Errorf("Failed to create target file: %v", err)
	}
	defer f.Close()

	if err := syscall.Mount(source, target, "bind", syscall.MS_RDONLY|syscall.MS_BIND, ""); err != nil {
		return nil, fmt.Errorf("Failed to mount: %v", err)
	}
	return func() { syscall.Unmount(target, 0) }, nil
}

func (r runtimeService) ListImages() (*[]image.Image, error) {
	imgs, err := r.imgStore.ListImages()
	if err != nil {
		return nil, err
	}

	images := make([]image.Image, 0)
	for _, img := range imgs {
		f, err := os.Open(img.SourceFile())
		if err != nil {
			return nil, err
		}
		defer f.Close()
		content, err := io.ReadAll(f)
		if err != nil {
			return nil, err
		}
		var j image.Image
		if err := json.Unmarshal(content, &j); err != nil {
			return nil, err
		}
		images = append(images, j)
	}
	return &images, nil
}

func (r runtimeService) FindImageByNameAndId(name string, tag string) (*image.Image, error) {
	if name == "" || tag == "" {
		return nil, nil
	}

	// Can be optimized: no need list all image first
	imgs, err := r.ListImages()
	if err != nil {
		return nil, err
	}
	if imgs == nil {
		return nil, nil
	}

	for _, img := range *imgs {
		if img.Name == name && img.Tag == tag {
			return &img, nil
		}
	}
	return nil, nil
}

func (r runtimeService) RemoveImage(name string, tag string) error {
	img, err := r.FindImageByNameAndId(name, tag)
	if err != nil {
		return err
	}
	if img == nil {
		return fmt.Errorf("image not found")
	}
	return r.imgStore.RemoveImage(img.Digest)
}

func (r runtimeService) ListContainers() (*[]container.Container, error) {
	conts, err := r.conStore.ListContainers()
	if err != nil {
		return nil, err
	}

	containers := make([]container.Container, 0)
	for _, c := range conts {
		f, err := os.Open(c.SpecFile())
		if err != nil {
			return nil, err
		}
		defer f.Close()
		content, err := io.ReadAll(f)
		if err != nil {
			return nil, err
		}
		var j container.Container
		if err := json.Unmarshal(content, &j); err != nil {
			return nil, err
		}
		containers = append(containers, j)
	}
	return &containers, nil
}

func (r runtimeService) RemoveContainer(id string) error {
	if id == "" {
		return fmt.Errorf("container id required")
	}
	return r.conStore.RemoveContainer(id)
}
