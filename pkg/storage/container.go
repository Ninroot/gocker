package storage

import (
	"log"
	"os"
	"path/filepath"

	"github.com/containerd/btrfs"
	"github.com/ninroot/gocker/pkg/util"
)

type ContainerHandle struct {
	id      string
	contDir string
}

func NewContainerHandle(id string, contDir string) *ContainerHandle {
	return &ContainerHandle{
		id:      id,
		contDir: contDir,
	}
}

// /var/btrfs/cont/abc/rootfs/
func (h *ContainerHandle) RootfsDir() string {
	return filepath.Join(h.contDir, "rootfs")
}

// /var/btrfs/cont/abc/
func (h *ContainerHandle) ContDir() string {
	return filepath.Join(h.contDir)
}

type ContainerStore struct {
	rootDir string
}

func NewContainerStore(rootDir string) *ContainerStore {
	return &ContainerStore{
		rootDir: rootDir,
	}
}

func (s *ContainerStore) RemoveContainer(id string) error {
	return os.RemoveAll(s.ContainerDir(id))
}

// /var/btrfs/cont/
func (s *ContainerStore) RootDir() string {
	return s.rootDir
}

// /var/btrfs/cont/abc/
func (s *ContainerStore) ContainerDir(id string) string {
	return filepath.Join(s.RootDir(), id)
}

func (s *ContainerStore) GetContainer(id string) *ContainerHandle {
	d := filepath.Join(s.RootDir(), id)
	ok, err := util.Exist(d)
	if err != nil {
		log.Println("Warning: ", err)
		return nil
	}
	if !ok {
		return nil
	}
	return NewContainerHandle(id, d)
}

func (s *ContainerStore) ListContainers() ([]*ContainerHandle, error) {
	items, err := os.ReadDir(s.RootDir())
	if err != nil {
		return nil, err
	}
	conts := make([]*ContainerHandle, 0)
	for _, item := range items {
		if item.IsDir() {
			contId := item.Name()
			conts = append(conts, NewContainerHandle(contId, filepath.Join(s.RootDir(), contId)))
		}
	}
	return conts, nil
}

func (s *ContainerStore) CreateContainer(id string, imagePath string) (*ContainerHandle, error) {
	contDir := filepath.Join(s.RootDir(), id)
	if err := btrfs.SubvolSnapshot(contDir, imagePath, false); err != nil {
		return nil, err
	}
	return NewContainerHandle(id, contDir), nil
}
