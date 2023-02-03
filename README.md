# Gocker

![gocker](https://user-images.githubusercontent.com/11426226/212131970-c8f78c2c-3441-44d9-bffb-07793f145e87.png)

A Docker implementation written in Golang designed for educational purposes. We __do not recommend using it in production environments__, and suggest running it inside a virtual machine instead.

## Prerequisites

Gocker runs only on Linux-based system with version 3.10 or higher of the Linux kernel.

Required packages:
- libcgroup-tools

Required configuration:
- A btrfs filesystem mounted under /var/gocker (configurable)
- A cgroup filesystem mounted under /sys/fs/cgroup/ (configurable) if not already the case

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

## Dev comments

`CLONE_NEWUTS` create the process in a new UTS namespace.
Allow the set of a new hostname inside the container (using for example `hostname <name>`) without affecting the hostname of the host.

`CLONE_NEWUSER` creates the process in a new user namespace.
if set without `SysProcAttr.Credential`, `id` returns uid=65534(nobody) gid=65534(nogroup) groups=65534(nogroup)

`CLONE_NEWPID` creates the process in a new PID namespace resulting in having a PID equal to 1.
`echo $$` returns 1
Use of CLONE_NEWPID requires the CAP_SYS_ADMIN capability. In other words, sudo is required to run the binary.
if set without `SysProcAttr.Credential`, `id` returns uid=0(root) gid=0(root) groups=0(root)

`CLONE_NEWNS` the child process have a different mount namespace to its.
Use of CLONE_NEWNS requires the CAP_SYS_ADMIN capability.

Great souces to learn Docker:
- [How Docker Works - Intro to Namespaces](https://youtu.be/-YnMr1lj4Z8)
- [Containers the hard way: Gocker: A mini Docker written in Go](https://unixism.net/2020/06/containers-the-hard-way-gocker-a-mini-docker-written-in-go/) - Detailed article explaining another Docker implementation in Golang (not related to this project despite a similar name)
- [Bocker](https://github.com/p8952/bocker) - Docker implemented in around 100 lines of bash
- [Containers From Scratch with Golang
](https://medium.com/@ssttehrani/containers-from-scratch-with-golang-5276576f9909)

## TODO
- [x] change hostname inside the container using the `/proc/self/exec` trick described in the medium article
- [x] basic chroot
- [x] basic pull with API
- [x] docker image list
- [x] docker image rm
- [x] copy on write fs for the image
- [x] enable resources limitation with cgroups
- [ ] return the exit code of the container

Bugs
- [ ] conts/ is created by root resulting in the necessity of using sudo to rm. Happens also with img/ when it hasn't be created beforhand and run is invoked with sudo resulting in the incapacity to pull the image without root.