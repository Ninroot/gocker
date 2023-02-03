# Gocker

<img src="https://user-images.githubusercontent.com/11426226/212131970-c8f78c2c-3441-44d9-bffb-07793f145e87.png" height="200" alt="gocker logo">

A Docker implementation written in Golang designed for educational purposes. We __do not recommend using it in production environments__, and suggest running it inside a virtual machine instead.

## Prerequisites

Gocker runs only on Linux-based system with version 3.10 or higher of the Linux kernel.

Required packages:
- `libcgroup-tools`

Required configuration:
- A btrfs filesystem mounted under `/var/gocker` (configurable)
- A cgroup filesystem mounted under `/sys/fs/cgroup/` (configurable) if not already the case

## Install

```bash
git clone https://github.com/Ninroot/gocker
cd gocker/
make
cd build/
./gocker --help
```

## Example of use

When using Gocker, you need to be specific when pulling or running an image, as it does not have the same magic as Docker does. For example, use `pull amd64/alpine` instead of just `pull alpine`.

```shell
# gocker requires root privileges
sudo su

# image for ARM-based system (like mac running on Apple silicon)
./gocker pull arm64v8/alpine

./gocker run arm64v8/alpine:latest /bin/sh

./gocker image rm arm64v8/alpine
```

## Isolation

- file system: via chroot, gives the illusion the container can navigate inside a different distro.
- PID: via namespace, isolates the processes running inside a container with processes of the host / other containers.
- Mount: via namespace, isolates the mounts to the container. This prevents the mounts from being visible from the host.
- UTS: via namespace, allows to set a new hostname inside the container without affecting the hostname of the host.

## Learn

Great souces to learn Docker:
- [How Docker Works - Intro to Namespaces](https://youtu.be/-YnMr1lj4Z8)
- [Containers the hard way: Gocker: A mini Docker written in Go](https://unixism.net/2020/06/containers-the-hard-way-gocker-a-mini-docker-written-in-go/) - Detailed article explaining another Docker implementation in Golang (not related to this project despite a similar name)
- [Bocker](https://github.com/p8952/bocker) - Docker implemented in around 100 lines of bash
- [Containers From Scratch with Golang
](https://medium.com/@ssttehrani/containers-from-scratch-with-golang-5276576f9909)
