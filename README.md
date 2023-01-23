# Gocker

![gocker](https://user-images.githubusercontent.com/11426226/212131970-c8f78c2c-3441-44d9-bffb-07793f145e87.png)

A Docker implementation written in Golang designed for educational purposes. We __do not recommend using it in production environments__, and suggest running it inside a virtual machine instead.

## Prerequisites

Gocker runs only on Linux-based system with version 3.10 or higher of the Linux kernel.

Required packages:
- libcgroup-tools

Required configuration:
- A btrfs filesystem mounted under /var/gocker

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
# image for ARM-based system (like mac running on Apple silicon)
./gocker pull arm64v8/alpine

# requires sudo
sudo ./gocker run arm64v8/alpine

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
- https://youtu.be/-YnMr1lj4Z8
- https://medium.com/@ssttehrani/containers-from-scratch-with-golang-5276576f9909

## TODO
- [x] change hostname inside the container using the `/proc/self/exec` trick described in the medium article
- [x] basic chroot
- [x] basic pull with API
- [x] docker image list
- [x] docker image rm
- [ ] copy on write fs for the image 

Bugs
- [ ] conts/ is created by root resulting in the necessity of using sudo to rm. Happens also with img/ when it hasn't be created beforhand and run is invoked with sudo resulting in the incapacity to pull the image without root 