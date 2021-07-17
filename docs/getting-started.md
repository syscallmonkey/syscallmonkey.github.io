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



## Quick start

You now have a Syscall Monkey available. Time for some mischief!

Here's a quick example to get you started:

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

## Good to go!

You now have a Syscall Monkey available. Time for some mischief!
