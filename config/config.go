package config

import "github.com/ninroot/gocker/pkg/logging"

const (
	DefaultImageStoreRootDir     = "/var/gocker/img"
	DefaultContainerStoreRootDir = "/var/gocker/cont"
	DefaultCGroupDir             = "/sys/fs/cgroup/"
	DefaultRegistry              = "https://registry-1.docker.io/"

	DefaultLogLevel = logging.Debug
)
