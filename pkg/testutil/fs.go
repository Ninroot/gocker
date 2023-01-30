package testutil

import (
	"os"
	"os/exec"
)

type Testfs struct{}

func (f Testfs) SubvolCreate(path string) error {
	return os.Mkdir(path, 755)
}

func (f Testfs) SubvolSnapshot(dst string, src string) error {
	cmd := exec.Command("cp", "--recursive", src, dst)
	return cmd.Run()
}
