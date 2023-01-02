package pkg

import (
	"io"
	"os/exec"
)

type ImageStore struct {
	rootDir string
}

func NewImageStore(rootDir string) ImageStore {
	return ImageStore{
		rootDir: rootDir,
	}
}

type RegistryService struct {
	imgStore ImageStore
}

func NewRegistryService(imgStore ImageStore) RegistryService {
	return RegistryService{
		imgStore: imgStore,
	}
}

func (reg *RegistryService) Pull(image string) error {
	// return errors.New("bad")
	export := exec.Command("docker", "export", image)
	untar := exec.Command("tar", "-C", reg.imgStore.rootDir, "-xvf", "-")
	r, w := io.Pipe()

	export.Stdout = w
	untar.Stdin = r

	if err := export.Start(); err != nil {
		return err
	}
	if err := untar.Start(); err != nil {
		return err
	}
	if err := export.Wait(); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	if err := untar.Wait(); err != nil {
		return err
	}

	return nil
}
