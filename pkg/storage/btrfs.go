package storage

import (
	"github.com/containerd/btrfs"
)

type COWFS interface {
	SubvolCreate(path string) error
	SubvolSnapshot(dst string, src string) error
}

type Btrfs struct {
}

// SubvolCreate creates a subvolume at the provided path.
func (b Btrfs) SubvolCreate(path string) error {
	return btrfs.SubvolCreate(path)
}

// SubvolSnapshot creates a snapshot in dst from src.
func (b Btrfs) SubvolSnapshot(dst string, src string) error {
	return btrfs.SubvolSnapshot(dst, src, false)
}
