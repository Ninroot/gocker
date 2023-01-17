package storage

import "path/filepath"

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

func (h *ContainerHandle) RootfsDir() string {
	return filepath.Join(h.contDir, "rootfs")
}
