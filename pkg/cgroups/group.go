package cgroups

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/sirupsen/logrus"
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

func (g Group) getCpuDir() string {
	return filepath.Join(g.cgroupDir, "cpu", "gocker", g.name)
}

func (g Group) getCfsPeriodFile() string {
	return filepath.Join(g.getCpuDir(), "cpu.cfs_period_us")
}

func (g Group) getCfsQuotaFile() string {
	return filepath.Join(g.getCpuDir(), "cpu.cfs_quota_us")
}

// setCpuLimit sets the limit in number of CPU
func (g Group) SetCpuLimit(cpu int) error {
	if cpu > runtime.NumCPU() {
		logrus.Info("CPU limit has been set to max")
		cpu = runtime.NumCPU()
	}
	period := 1000000
	// amount of time, in microseconds, that a CPU can run in a single scheduling period
	if err := write(g.getCfsPeriodFile(), period); err != nil {
		return fmt.Errorf("could not set CFS Period: %v", err)
	}
	// how much cpu resource can be used in every period
	if err := write(g.getCfsQuotaFile(), cpu*period); err != nil {
		return fmt.Errorf("could not set CFS Quota: %v", err)
	}
	return nil
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
