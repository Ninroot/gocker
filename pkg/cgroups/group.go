package cgroups

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

type CGroup struct {
	rootfs string
}

func New(rootfs string) CGroup {
	return CGroup{
		rootfs: rootfs,
	}
}

type Group struct {
	cgroupDir string
	name      string
}

func (c CGroup) NewGroup(name string) Group {
	return Group{
		cgroupDir: c.rootfs,
		name:      name,
	}
}

func (g Group) Delete() error {
	return os.Remove(g.getPidsDir())
}

// /sys/fs/cgroup/memory/gocker/abc
func (g Group) getMemoryDir() string {
	return filepath.Join(g.cgroupDir, "memory", "gocker", g.name)
}

// /sys/fs/cgroup/memory/gocker/abc/pids.max
func (g Group) getLimitMemoryFile() string {
	return filepath.Join(g.getMemoryDir(), "memory.limit_in_bytes")
}

// SetMemoryLimit limits the memory of the process in Byte.
func (g Group) SetMemoryLimit(max int) error {
	return write(g.getLimitMemoryFile(), max)
}

// /sys/fs/cgroup/pids/gocker/abc
func (g Group) getPidsDir() string {
	return filepath.Join(g.cgroupDir, "pids", "gocker", g.name)
}

// /sys/fs/cgroup/pids/gocker/abc/pids.max
func (g Group) getPidMaxFile() string {
	return filepath.Join(g.getPidsDir(), "pids.max")
}

// /sys/fs/cgroup/pids/gocker/abc/cgroup.procs
func (g Group) getProcsFile() string {
	return filepath.Join(g.getPidsDir(), "cgroup.procs")
}

// /sys/fs/cgroup/pids/gocker/abc/notify_on_release
func (g Group) getNotifyOnReleaseFile() string {
	return filepath.Join(g.getPidsDir(), "notify_on_release")
}

// SetPidMax limits the number of child processes to prevent crashes like forkbomb.
func (g Group) SetPidMax(max int) error {
	return write(g.getPidMaxFile(), max)
}

// Add Proc adds the process to the cgroup.
func (g Group) AddProc(pid int) error {
	return write(g.getProcsFile(), pid)
}

func (g Group) SetNotifyOnRelease(b bool) error {
	v := 0
	if b {
		v = 1
	}
	return write(g.getNotifyOnReleaseFile(), v)
}

func write(path string, value int) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(path, []byte(strconv.Itoa(value)), 0700)
}
