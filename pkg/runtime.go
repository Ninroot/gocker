package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"

	"github.com/ninroot/gocker/config"
	"github.com/ninroot/gocker/pkg/cgroups"
	"github.com/ninroot/gocker/pkg/container"
	"github.com/ninroot/gocker/pkg/image"
	"github.com/ninroot/gocker/pkg/storage"
	"github.com/ninroot/gocker/pkg/util"
)

type runtimeService struct {
	imgStore storage.ImageStore
	conStore storage.ContainerStore
	cgroup   cgroups.CGroup
}

func NewRuntimeService() *runtimeService {
	return &runtimeService{
		imgStore: storage.NewImageStore(util.EnsureDir(config.DefaultImageStoreRootDir)),
		conStore: *storage.NewContainerStore(util.EnsureDir(config.DefaultContainerStoreRootDir)),
		cgroup:   cgroups.New(util.EnsureDir(config.DefaultCGroupDir)),
	}
}

func Run() error {
	args := append([]string{"tech"}, os.Args[2:]...)
	cmd := exec.Command("/proc/self/exe", args...)

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

	r.applyCGroup(cmd.Process.Pid)

	return cmd.Wait()
}

func (r runtimeService) applyCGroup(pid int) error {
	log.Println("Setting cgroup for pid", pid)
	g := r.cgroup.NewGroup(strconv.Itoa(pid))
	if err := g.SetPidMax(10); err != nil {
		return err
	}
	if err := g.SetNotifyOnRelease(true); err != nil {
		return err
	}
	return g.AddProc(pid)
}

func getCommand(cmd []string) (command string, args []string) {
	if len(cmd) == 0 {
		return "", []string{}
	}
	return cmd[0], cmd[1:]
}

func (r runtimeService) InitContainer(imageName, imageTag, containerName string, containerCmd []string) error {
	img, err := r.FindImageByNameAndId(imageName, imageTag)
	if err != nil {
		return err
	}
	if img == nil {
		return fmt.Errorf("image not found: %s:%s", imageName, imageTag)
	}

	imgH := r.imgStore.GetImage(img.Digest)

	uuid := container.RandID()
	contH, err := r.conStore.CreateContainer(uuid, imgH.ImageDir())
	if err != nil {
		return nil
	}

	cmdName, cmdArgs := getCommand(containerCmd)

	c := container.Container{
		ID:        uuid,
		Name:      containerName,
		Image:     *img,
		CreatedAt: time.Now(),
		Command:   cmdName,
		Args:      cmdArgs,
	}

	contH.SetSpec(c)

	if err := syscall.Chroot(contH.RootfsDir()); err != nil {
		return err
	}
	if err := syscall.Chdir("/"); err != nil {
		return err
	}

	// hostname will be affected if this function runs in a process that hasn't been with CLONE_NEWUTS
	// happens typically when debugging
	syscall.Sethostname([]byte(uuid))

	// mount /proc to make commands such `ps` working
	syscall.Mount("proc", "proc", "proc", 0, "")
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
