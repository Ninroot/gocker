package pkg

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/heroku/docker-registry-client/registry"
	"github.com/ninroot/gocker/config"
	"golang.org/x/term"
)

type ImageStore struct {
	rootDir string
}

type ImageId struct {
	Name   string `json:"name"`
	Tag    string `json:"tag"`
	Digest string `json:"digest"`
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
	image, err := parse(imageName)
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

	if image.Tag == "" {
		image.Tag = "latest"
	}
	manifest, err := hub.ManifestV2(image.Name, image.Tag)
	if err != nil {
		return err
	}
	log.Printf("Found manifest for image <%s:%s>", image.Name, image.Tag)

	digest := manifest.Layers[0].Digest
	reader, err := hub.DownloadBlob(image.Name, digest)
	if err != nil {
		return err
	}

	image.Digest = string(digest)

	if err := CreateImage(reader, image); err != nil {
		return err
	}

	return nil
}

// name[:TAG]
// return repository and tag
func parse(imageName string) (ImageId, error) {
	s := strings.Split(imageName, ":")
	if len(s) == 1 {
		return ImageId{Name: s[0]}, nil
	}
	if len(s) == 2 {
		return ImageId{Name: s[0], Tag: s[1]}, nil
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
