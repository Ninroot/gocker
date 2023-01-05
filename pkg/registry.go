package pkg

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/heroku/docker-registry-client/registry"
	"github.com/ninroot/gocker/config"
	"golang.org/x/term"
)

type ImageStore struct {
	rootDir string
}

type ImageId struct {
	name string
	tag  string
}

func NewImageStore(rootDir string) ImageStore {
	return ImageStore{
		rootDir: rootDir,
	}
}

type RegistryService struct {
	registry string
	imgStore ImageStore
}

func NewRegistryService(imgStore ImageStore) RegistryService {
	return RegistryService{
		registry: config.DefaultRegistry,
		imgStore: imgStore,
	}
}

func (reg *RegistryService) Pull(imageName string) error {
	imageId, err := parse(imageName)
	if err != nil {
		return err
	}

	username, password, err := login(reg.registry)
	if err != nil {
		return err
	}

	hub, err := registry.New(reg.registry, username, password)
	if err != nil {
		return err
	}

	if imageId.tag == "" {
		imageId.tag = "latest"
	}
	manifest, err := hub.ManifestV2(imageId.name, imageId.tag)
	if err != nil {
		return err
	}
	log.Printf("Found manifest for image <%s:%s>", imageId.name, imageId.tag)

	digest := manifest.Layers[0].Digest
	reader, err := hub.DownloadBlob(imageId.name, digest)
	if err != nil {
		return err
	}

	name := filepath.Base(imageId.name) + ".tar"
	file, err := os.Create(name)
	if err != nil {
		return err
	}

	defer file.Close()
	io.Copy(file, reader)

	return nil
}

// name[:TAG]
// return repository and tag
func parse(imageName string) (ImageId, error) {
	s := strings.Split(imageName, ":")
	if len(s) == 1 {
		return ImageId{name: s[0]}, nil
	}
	if len(s) == 2 {
		return ImageId{name: s[0], tag: s[1]}, nil
	}
	return ImageId{}, errors.New("image name has the wrong format")
}

func login(registry string) (string, string, error) {
	username := os.Getenv("GOCKER_REGISTRY_USERNAME")
	password := os.Getenv("GOCKER_REGISTRY_PASSWORD")

	if username == "" && password == "" {
		fmt.Printf("To interact with the registry <%s>, credentials are required.\n", registry)
	}
	if username == "" {
		if _, err := fmt.Scanf("%s", &username); err != nil {
			return "", "", err
		}
	}
	if password == "" {
		p, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return "", "", err
		}
		password = string(p)
	}

	return username, password, nil
}
