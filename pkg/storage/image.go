package storage

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/containerd/btrfs"
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

// /var/btrfs/img/
func (s ImageStore) RootDir() string {
	return s.rootDir
}

// /var/btrfs/img/abc/
func (s ImageStore) ImageDir(id string) string {
	return filepath.Join(s.RootDir(), id)
}

func (s ImageStore) CreateImage(reader io.ReadCloser, id string) (*ImageHandle, error) {
	if err := os.MkdirAll(s.rootDir, 0700); err != nil {
		return nil, err
	}

	iDir := s.ImageDir(id)
	prs, err := util.Exist(iDir)
	if err != nil {
		return nil, err
	}

	h := NewImageHandle(id, iDir)

	if prs {
		log.Printf("image <%s> already exists", iDir)
		return h, nil
	}

	if err := btrfs.SubvolCreate(iDir); err != nil {
		return h, err
	}

	if err := Untar(h.RootfsDir(), reader); err != nil {
		return h, err
	}

	log.Printf("image <%s> stored in <%s>", h.id, s.rootDir)
	return h, nil
}

func (s ImageStore) GetImage(id string) *ImageHandle {
	d := filepath.Join(s.RootDir(), id)
	ok, err := util.Exist(d)
	if err != nil {
		log.Println("Warning: ", err)
		return nil
	}
	if !ok {
		return nil
	}
	return NewImageHandle(id, d)
}

func (s ImageStore) RemoveImage(id string) error {
	return os.RemoveAll(s.ImageDir(id))
}

func (s ImageStore) ListImages() ([]*ImageHandle, error) {
	items, err := os.ReadDir(s.RootDir())
	if err != nil {
		return nil, err
	}
	images := make([]*ImageHandle, 0)
	for _, item := range items {
		if item.IsDir() {
			imageId := item.Name()
			images = append(images, NewImageHandle(imageId, filepath.Join(s.RootDir(), imageId)))
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

// /var/btrfs/img/abc/
func (h ImageHandle) ImageDir() string {
	return h.imageDir
}

// /var/btrfs/img/abc/rootfs
func (h ImageHandle) RootfsDir() string {
	return filepath.Join(h.ImageDir(), "rootfs")
}

// /var/btrfs/img/abc/source
func (h ImageHandle) SourceFile() string {
	return filepath.Join(h.ImageDir(), "source")
}

func (h ImageHandle) SetSource(content any) error {
	f, err := os.OpenFile(h.SourceFile(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(f)
	if err := encoder.Encode(content); err != nil {
		return err
	}
	return nil
}
