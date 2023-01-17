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

// /var/btrfs/cont/
func (s *ContainerStore) RootDir() string {
	return s.rootDir
}

func (s *ContainerStore) GetContainer(id string) *ContainerHandle {
	d := filepath.Join(s.RootDir(), id)
	ok, err := util.Exist(d)
	if err == nil {
		log.Println("Warning: ", err)
		return nil
	}
	if !ok {
		return nil
	}
	return NewContainerHandle(id, d)
}

func (s *ContainerStore) CreateContainer(id string, imagePath string) (*ContainerHandle, error) {
	contDir := filepath.Join(s.RootDir(), id)
	if err := os.MkdirAll(contDir, 0700); err != nil {
		return nil, err
	}

	if err := btrfs.SubvolSnapshot(contDir, imagePath, false); err != nil {
		return nil, err
	}
	return NewContainerHandle(id, contDir), nil
}
