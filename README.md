# Gocker

A Docker implementation written in Golang designed for educational purposes. We __do not recommend using it in production environments__, and suggest running it inside a virtual machine instead.

## Example of use

```shell
# for arm image
./build/gocker pull arm64v8/alpine
./build/gocker run arm64v8/alpine
```

## Dev comments

`CLONE_NEWUTS` create the process in a new UTS namespace.
Allow the set of a new hostname inside the container (using for example `hostname <name>`) without affecting the hostname of the host.

`CLONE_NEWUSER` creates the process in a new user namespace.
if set without `SysProcAttr.Credential`, `id` returns uid=65534(nobody) gid=65534(nogroup) groups=65534(nogroup)

`CLONE_NEWPID` creates the process in a new PID namespace resulting in having a PID equal to 1.
`echo $$` returns 1
Only a privileged process (CAP_SYS_ADMIN) can employ `CLONE_NEWPID`. In other words, sudo is required to run the binary.
if set without `SysProcAttr.Credential`, `id` returns uid=0(root) gid=0(root) groups=0(root)

Great souces to learn Docker:
- https://youtu.be/-YnMr1lj4Z8
- https://medium.com/@ssttehrani/containers-from-scratch-with-golang-5276576f9909

## TODO
- [x] change hostname inside the container using the `/proc/self/exec` trick described in the medium article
- [x] basic chroot
- [x] basic pull with API
- [x] docker image list
- [ ] copy on write fs for the image 