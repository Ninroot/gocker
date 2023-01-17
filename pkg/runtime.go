package pkg

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"syscall"

	"github.com/ninroot/gocker/config"
	"github.com/ninroot/gocker/pkg/image"
	"github.com/ninroot/gocker/pkg/storage"
	"github.com/ninroot/gocker/pkg/util"
)

type runtimeService struct {
	imgStore storage.ImageStore
}

func NewRuntimeService() *runtimeService {
	return &runtimeService{
		imgStore: storage.NewImageStore(util.EnsureDir(config.DefaultImageStoreRootDir)),
	}
}

func Run(args []string) {
	log.Println("run with args", args)
	cmd := exec.Command("/proc/self/exe", append([]string{"tech"}, args...)...)
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

func (r runtimeService) InitContainer(args []string) error {
	log.Printf("Init with args %v, PID: %v", args, os.Getpid())
	input, err := image.Parse(args[0])
	if err != nil {
		return err
	}

	img, err := r.FindImageByNameAndId(input.Name, input.Tag)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	syscall.Sethostname([]byte(filepath.Base(img.Name)))

	// TODO finish
	p := path.Join(r.imgStore.rootDir, image.Digest, "rootfs")
	if err := syscall.Chroot(p); err != nil {
		return err
	}
	if err := syscall.Chdir("/"); err != nil {
		return err
	}

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

func (r runtimeService) FindImageByNameAndId(name string, tag string) (*image.Image, error) {
	if name == "" || tag == "" {
		return nil, nil
	}

	imgs, err := r.imgStore.ListImages()
	if err != nil {
		return nil, err
	}

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
		if j.Name == name && j.Tag == tag {
			return &j, nil
		}
	}
	return nil, nil
}
