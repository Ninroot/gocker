package storage

import (
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/ninroot/gocker/pkg/container"
)

type Testfs struct{}

func (f Testfs) SubvolCreate(path string) error {
	return os.Mkdir(path, 755)
}

func (f Testfs) SubvolSnapshot(dst string, src string) error {
	cmd := exec.Command("cp", "--recursive", src, dst)
	return cmd.Run()
}

func TestCreateContainer(t *testing.T) {
	img := MkdirTempTest(t, "img")
	defer os.RemoveAll(img)
	con := MkdirTempTest(t, "con")
	defer os.RemoveAll(con)

	s := NewContainerStore(img, Testfs{})
	_, err := s.CreateContainer(container.RandID(), con)
	if err != nil {
		log.Fatal("ContainerStore failed: CreateContainer")
	}
}

func MkdirTempTest(t *testing.T, pattern string) string {
	d, err := os.MkdirTemp("", t.Name()+pattern)
	if err != nil {
		log.Fatal(err)
	}
	return d
}
