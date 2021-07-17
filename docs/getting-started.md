## Installation

Syscall Monkey is written in Go, and there are a few different ways you can get started.

### Binary

You can build the binary from the source:

```sh
git clone https://github.com/syscallmonkey/monkey.git
cd monkey
make bin/monkey
./bin/monkey -h
```

Note, that if you're running on MacOS, you can build for Linux, but unfortunately MacOS's version of `ptrace` doesn't allow any of this magic to happen.

### Pull Docker container

When using with Kubernetes, you can use our official docker images.

```sh
docker pull ghcr.io/syscallmonkey/monkey:0.0.1rc1
```

Check [the latest available versions here](https://github.com/syscallmonkey/monkey/pkgs/container/monkey).

Note, that these container builds are minimal, and don't include things like bash.

```sh
$ docker run --rm -ti --cap-add SYS_PTRACE ghcr.io/syscallmonkey/monkey:0.0.1 -h
Usage:
  monkey [OPTIONS]

Application Options:
  -p, --attach=  Attach to the specified pid
  -t, --target=  Attach to process matching this name
  -c, --config=  Configuration file with desired scenario
  -o, --output=  Write the tracing output to the file (instead of stdout)
  -C, --summary  Show verbose debug information
  -s, --silent   Don't display tracing info

Help Options:
  -h, --help     Show this help message
```


### Building Docker container

If you're like to build the container locally from the source, it's easy:

```sh
git clone https://github.com/syscallmonkey/monkey.git
cd monkey
make build
make run
```

`make run` is a shortcut that will start the newly built image and run a bash session inside:

```sh
root@3e14fcd5843c:/# monkey -h
Usage:
  monkey [OPTIONS]

Application Options:
  -p, --attach=  Attach to the specified pid
  -t, --target=  Attach to process matching this name
  -c, --config=  Configuration file with desired scenario
  -o, --output=  Write the tracing output to the file (instead of stdout)
  -C, --summary  Show verbose debug information
  -s, --silent   Don't display tracing info

Help Options:
  -h, --help     Show this help message
```

### Compatibility

Currently, only `Linux` on `x86_64` is supported. If you need arm support, file an issue.


## Running Syscall Monkey

There are two ways of running Syscall Monkey:

- start a process (append the command at the end)
- attach to a running process (`-p` flag)


### Start a new process

To start a new process and manipulate it, just append the command at the end of the `monkey` command.

What you're going to see is the list of all `syscalls` made by the program, printed in the following format:

```sh
SYSCALL_NAME(ARG_NAME=ARG_VALUE) = RETURN_CODE (OPTIONALLY ERROR DESCRIPTION)
```

For example, let's see what `syscalls` a simple `sleep 1` makes:

```sh
root@c69219c773ff:/# monkey sleep 1
Version v0.0.1, build Sat Jul 17 10:30:03 UTC 2021
Started new process pid 21
execve(filename=NULL, argv=0, envp=0) = 0
brk(brk=0) = 94331848404992
arch_prctl(task=12289, code=140722511137424, addr=140631102649024) = -1 (errno 22: invalid argument)
access(filename=/etc/ld.so.preload, mode=4) = -1 (errno 2: no such file or directory)
openat(dfd=4294967196, filename=/etc/ld.so.cache, flags=524288, mode=0) = 3
fstat(fd=3, statbuf=140722511133840) = 0
mmap(addr=0, len=6530, prot=1, flags=2, fd=3, off=0) = 140631102525440
close(fd=3) = 0
openat(dfd=4294967196, filename=/lib/x86_64-linux-gnu/libc.so.6, flags=524288, mode=0) = 3
read(fd=3, buf=, count=832) = 832
pread64(fd=3, buf=, count=784, pos=64) = 784
pread64(fd=3, buf=, count=32, pos=848) = 32
pread64(fd=3, buf=, count=68, pos=880) = 68
fstat(fd=3, statbuf=140722511133920) = 0
mmap(addr=0, len=8192, prot=3, flags=34, fd=4294967295, off=0) = 140631102517248
pread64(fd=3, buf=, count=784, pos=64) = 784
pread64(fd=3, buf=, count=32, pos=848) = 32
pread64(fd=3, buf=, count=68, pos=880) = 68
mmap(addr=0, len=2036952, prot=1, flags=2050, fd=3, off=0) = 140631100477440
mprotect(start=140631100628992, len=1847296, prot=0) = 0
mmap(addr=140631100628992, len=1540096, prot=5, flags=2066, fd=3, off=151552) = 140631100628992
mmap(addr=140631102169088, len=303104, prot=1, flags=2066, fd=3, off=1691648) = 140631102169088
mmap(addr=140631102476288, len=24576, prot=3, flags=2066, fd=3, off=1994752) = 140631102476288
mmap(addr=140631102500864, len=13528, prot=3, flags=50, fd=4294967295, off=0) = 140631102500864
close(fd=3) = 0
arch_prctl(task=4098, code=140631102522752, addr=18446603442607026496) = 0
mprotect(start=140631102476288, len=12288, prot=1) = 0
mprotect(start=94331820990464, len=4096, prot=1) = 0
mprotect(start=140631102717952, len=4096, prot=1) = 0
munmap(addr=140631102525440, len=6530) = 0
brk(brk=0) = 94331848404992
brk(brk=94331848540160) = 94331848540160
clock_nanosleep(which_clock=0, flags=0, rqtp=140722511137168, rmtp=0) = 0
close(fd=1) = 0
close(fd=2) = 0
exit_group(error_code=0)
--- program exited ---
```

Let's take a look at an interesting lines.

#### access syscall

[`access`](https://linux.die.net/man/2/access) syscall check real user's permissions for a file. Its signature is:

```c
int access(const char *pathname, int mode);
```

In the output, Syscall Monkey intercepted the syscall, and we can see that the program check the file `/etc/ld.so.preload`, and it gets an error number 2 "no such file or directory":

```python
access(filename=/etc/ld.so.preload, mode=4) = -1 (errno 2: no such file or directory)
```


### Attaching to a running process

To attach to a process that's already running, you can specify the PID of the target via `-p` flag. For example:

```sh
root@c69219c773ff:/# sleep 5&
[1] 10
root@c69219c773ff:/# monkey -p 10
...
restart_syscall() = 0
close(fd=1) = 0
close(fd=2) = 0
exit_group(error_code=0)
--- program exited ---
[1]+  Done                    sleep 5
```


## Quick start

You now have a Syscall Monkey available. Here's a quick example to get you started:

### Change the return value of geteuid

Use this scenario to sometimes change the user returned by [`getuid`](https://linux.die.net/man/2/geteuid):

```yaml

# cat /etc/passwd
# root:x:0:0:root:/root:/bin/bash
# daemon:x:1:1:daemon:/usr/sbin:/usr/s
# bin:x:2:2:bin:/bin:/usr/sbin/nologin

# cat /examples/getuid-random.yml
rules:
  - name: probably daemon
    probability: 0.66
    match:
      name: geteuid
    modify:
      return: 1
  - name: but maybe bin
    probability: 0.5
    match:
      name: geteuid
    modify:
      return: 2
```

And this should be enough to confuse a lot of people:

```sh
root@f34cc94a6b6d:/# whoami
root
root@f34cc94a6b6d:/# monkey -s -c /examples/getuid-random.yml whoami
daemon
root@f34cc94a6b6d:/# monkey -s -c /examples/getuid-random.yml whoami
root
root@f34cc94a6b6d:/# monkey -s -c /examples/getuid-random.yml whoami
bin
```

## Time for mischief!

This should be enough to get you started. Have fun!
