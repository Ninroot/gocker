package config

const (
	DefaultImageStoreRootDir     = "/var/gocker/img"
	DefaultContainerStoreRootDir = "/var/gocker/cont"
	DefaultCGroupDir             = "/sys/fs/cgroup/"
	DefaultRegistry              = "https://registry-1.docker.io/"

	DefaultPidsLimit   = 32
	DefaultMemoryLimit = 64 * 1024 * 1024 // 64MB

	DefaultLogLevel = "debug"
)
