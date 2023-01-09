package pkg

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"syscall"

	"github.com/ninroot/gocker/config"
)

type runtimeService struct {
	regSvc RegistryService
}

func NewRuntimeService() *runtimeService {
	return &runtimeService{
		regSvc: NewRegistryService(
			NewImageStore(EnsureDir(config.DefaultImageStoreRootDir)),
		),
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

func (runtime runtimeService) InitContainer(args []string) error {
	log.Printf("Init with args %v, PID: %v", args, os.Getpid())
	inputImage, err := Parse(args[0])
	if err != nil {
		return err
	}

	image, err := runtime.regSvc.imgStore.FindImage(inputImage)
	if err != nil {
		return err
	}
	if image == nil {
		return fmt.Errorf("image <%s:%s> not found", image.Name, image.Digest)
	}
	syscall.Sethostname([]byte(filepath.Base(image.Name)))

	p := path.Join(runtime.regSvc.imgStore.rootDir, image.Digest, "rootfs")
	log.Println("rootpath", p)
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
