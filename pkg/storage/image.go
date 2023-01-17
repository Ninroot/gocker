package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/containerd/btrfs"
	"github.com/ninroot/gocker/pkg/container"
	"github.com/ninroot/gocker/pkg/util"
)

type ImageStore struct {
	rootDir string
}

func NewImageStore(rootDir string) ImageStore {
	return ImageStore{
		rootDir: rootDir,
	}
}

func (s ImageStore) RootDir() string {
	return s.rootDir
}

func (s ImageStore) ImageDir(id string) string {
	return filepath.Join("img", s.RootDir(), id)
}

func (s ImageStore) CreateImage(reader io.ReadCloser, id string) (*ImageHandle, error) {
	if err := os.MkdirAll(s.rootDir, 0700); err != nil {
		return nil, err
	}

	imageDir := s.ImageDir(id)
	prs, err := util.Exist(imageDir)
	if err != nil {
		return nil, err
	}

	h := NewImageHandle(id, imageDir)

	if prs {
		log.Printf("image <%s> already exists", imageDir)
		return h, nil
	}

	if err := btrfs.SubvolCreate(imageDir); err != nil {
		return h, err
	}

	if err := Untar(h.RootfsDir(), reader); err != nil {
		return h, err
	}

	source, err := os.Create(h.SourceFile())
	if err != nil {
		return h, err
	}
	encoder := json.NewEncoder(source)
	if err := encoder.Encode(h); err != nil {
		return h, err
	}

	log.Printf("image <%s> stored in <%s>", h.id, s.rootDir)
	return h, nil
}

func (s ImageStore) CreateContainer(id string) (string, error) {
	h := s.GetImage(id)
	if h == nil {
		return "", fmt.Errorf("image <%s> not found", id)
	}

	contDir := h.ContDir()
	if err := os.MkdirAll(contDir, 0700); err != nil {
		return "", err
	}

	uuid := container.RandID()

	// TODO move to store.container
	// filepath.Join(s.imageDir, "img", img.Digest
	if err := btrfs.SubvolSnapshot(filepath.Join(contDir, uuid), h.ImageDir(), false); err != nil {
		return "", err
	}
	return uuid, nil
}

func (s ImageStore) GetImage(id string) *ImageHandle {
	p := filepath.Join(s.RootDir(), id)
	ok, err := util.Exist(p)
	if err == nil {
		log.Println("Warning: ", err)
		return nil
	}
	if !ok {
		return nil
	}
	return NewImageHandle(id, s.rootDir)
}

func (s ImageStore) RemoveImage(id string) error {
	return os.RemoveAll(s.ImageDir(id))
}

// func (s ImageStore) ListImages() ([]*ImageHandle, error) {
// 	return s.findImages(nil)
// }

// FindImage returns the image found in the store or nil if not found.
// func (s ImageStore) FindImage(image *ImageHandle) (*ImageHandle, error) {
// 	found, err := s.findImages(image)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(found) == 0 {
// 		return nil, nil
// 	}
// 	return found[0], nil
// }

func (s ImageStore) ListImages() ([]*ImageHandle, error) {
	items, err := os.ReadDir(s.RootDir())
	if err != nil {
		return nil, err
	}
	images := make([]*ImageHandle, 0)
	for _, item := range items {
		if item.IsDir() {
			images = append(images, NewImageHandle(item.Name(), s.RootDir()))
		}
	}
	return images, nil
}

type ImageHandle struct {
	id       string
	imageDir string
}

func NewImageHandle(id string, imageDir string) *ImageHandle {
	return &ImageHandle{
		id:       id,
		imageDir: imageDir,
	}
}

func (h ImageHandle) ImageDir() string {
	return h.imageDir
}

func (h ImageHandle) RootfsDir() string {
	return filepath.Join(h.ImageDir(), "rootfs")
}

func (h ImageHandle) SourceFile() string {
	return filepath.Join(h.ImageDir(), "source")
}

func (h ImageHandle) ContDir() string {
	return filepath.Join(h.ImageDir(), "con")
}
