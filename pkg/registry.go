package pkg

import (
	"fmt"
	"log"
	"os"

	"github.com/heroku/docker-registry-client/registry"
	"github.com/ninroot/gocker/config"
	"golang.org/x/term"
)

type ImageId struct {
	// Full path name of the image. Which means it can include the respository. E.g: `library/alpine`
	Name   string `json:"name"`
	Tag    string `json:"tag"`
	Digest string `json:"digest"`
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
	image, err := Parse(imageName)
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

	if err := reg.imgStore.CreateImage(reader, image); err != nil {
		return err
	}

	return nil
}

func login(registry string) (string, string, error) {
	username := os.Getenv("GOCKER_REGISTRY_USERNAME")
	password := os.Getenv("GOCKER_REGISTRY_PASSWORD")

	if username == "" && password == "" {
		fmt.Printf("To interact with the registry <%s>, credentials are required.\n", registry)
	}
	if username == "" {
		fmt.Printf("Username: ")
		if _, err := fmt.Scanf("%s", &username); err != nil {
			return "", "", err
		}
	}
	if password == "" {
		fmt.Printf("Password: ")
		p, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return "", "", err
		}
		password = string(p)
	}

	return username, password, nil
}
