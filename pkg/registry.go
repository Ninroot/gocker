package pkg

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/heroku/docker-registry-client/registry"
	"github.com/ninroot/gocker/config"
	"golang.org/x/term"
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
	registry string
	imgStore ImageStore
}

func NewRegistryService(imgStore ImageStore) RegistryService {
	return RegistryService{
		registry: config.DefaultRegistry,
		imgStore: imgStore,
	}
}

func (reg *RegistryService) Pull(image string) error {
	username, password, err := login(reg.registry)
	if err != nil {
		return err
	}

	hub, err := registry.New(reg.registry, username, password)
	if err != nil {
		return err
	}

	repo := "arm64v8/alpine"
	tags, err := hub.Tags(repo)
	if err != nil {
		return err
	}
	fmt.Println(tags)

	manifest, err := hub.ManifestV2(repo, "latest")
	if err != nil {
		return err
	}
	fmt.Println(manifest)

	digest := manifest.Layers[0].Digest
	reader, err := hub.DownloadBlob(repo, digest)
	if err != nil {
		return err
	}

	name := filepath.Base(repo) + ".tar"
	file, err := os.Create(name)
	if err != nil {
		return err
	}

	defer file.Close()
	io.Copy(file, reader)

	return nil
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
