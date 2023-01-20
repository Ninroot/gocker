package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/ninroot/gocker/config"
	"github.com/ninroot/gocker/pkg/container"
	"github.com/ninroot/gocker/pkg/image"
	"github.com/ninroot/gocker/pkg/storage"
	"github.com/ninroot/gocker/pkg/util"
)

type runtimeService struct {
	imgStore storage.ImageStore
	conStore storage.ContainerStore
}

func NewRuntimeService() *runtimeService {
	return &runtimeService{
		imgStore: storage.NewImageStore(util.EnsureDir(config.DefaultImageStoreRootDir)),
		conStore: *storage.NewContainerStore(util.EnsureDir(config.DefaultContainerStoreRootDir)),
	}
}

func Run(imageName string, imageTag string, containerName string) {
	cmd := exec.Command("/proc/self/exe", append([]string{"tech"}, imageName, imageTag, containerName)...)
	// cmd := exec.Command("/bin/sh")
	// cmd.SysProcAttr = &syscall.SysProcAttr{}
	// cmd.SysProcAttr = &syscall.SysProcAttr{Cloneflags: syscall.CLONE_NEWUTS}
	// cmd.SysProcAttr = &syscall.SysProcAttr{Cloneflags: syscall.CLONE_NEWPID}
	// cmd.SysProcAttr = &syscall.SysProcAttr{Cloneflags: syscall.CLONE_NEWUSER}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
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

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func (r runtimeService) InitContainer(imageName string, imageTag string, containerName string) error {
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

	contH.SetSpec(container.Container{
		ID:        uuid,
		Name:      containerName,
		Image:     *img,
		CreatedAt: time.Now(),
	})

	if err := syscall.Chroot(contH.RootfsDir()); err != nil {
		return err
	}
	if err := syscall.Chdir("/"); err != nil {
		return err
	}

	syscall.Sethostname([]byte(uuid))

	// mount /proc to make commands such `ps` working
	syscall.Mount("proc", "proc", "proc", 0, "")
	defer syscall.Unmount("/proc", 0)

	cmd := exec.Command("/bin/sh")

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
