package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/ninroot/gocker/pkg/util"
	"github.com/sirupsen/logrus"
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

func (h *ContainerHandle) SetSpec(content any) error {
	f, err := os.OpenFile(h.SpecFile(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(f)
	if err := encoder.Encode(content); err != nil {
		return err
	}
	return nil
}

// /var/btrfs/img/abc/source
func (h *ContainerHandle) SpecFile() string {
	return filepath.Join(h.contDir, "spec.json")
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
	fs      COWFS
}

func NewContainerStore(rootDir string, fs COWFS) *ContainerStore {
	return &ContainerStore{
		rootDir: rootDir,
		fs:      fs,
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
		logrus.Warn(err)
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

// TODO: rename to bundle to match the terminology?
// https://github.com/opencontainers/runtime-spec#application-bundle-builders
func (s *ContainerStore) CreateContainer(id string, imagePath string) (*ContainerHandle, error) {
	contDir := filepath.Join(s.RootDir(), id)
	if err := s.fs.SubvolSnapshot(contDir, imagePath); err != nil {
		return nil, err
	}
	return NewContainerHandle(id, contDir), nil
}
