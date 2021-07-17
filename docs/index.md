# Welcome to Syscall Monkey

## TL;DR

`Syscall Monkey` is like next-gen, cloud native [`strace`](https://man7.org/linux/man-pages/man1/strace.1.html):

- attach and detach processes using [`ptrace`](https://man7.org/linux/man-pages/man2/ptrace.2.html) (Linux only)
- intercept and manipulate their [`syscalls`](https://man7.org/linux/man-pages/man2/syscalls.2.html) (block, change arguments, return value)
- prepare scenarios in a simple `yaml` format
- write advanced scenarios using `syscallmonkey` as an SDK

## Teaser

### Change the return value

Here's how you can trick `whoami` into thinking it runs as `daemon` user (1), intead of `root` (0)

```sh
root@f34cc94a6b6d:/# whoami
root
```

Write this to `scenario.yml` to always return 1 for the [`getuid`](https://linux.die.net/man/2/geteuid) syscall:

```yaml
# cat /examples/getuid-user1.yml
rules:
  - name: switch geteuid to return a different user ID
    match:
      name: geteuid
    modify:
      return: 1
```

```sh
root@f34cc94a6b6d:/# monkey -s -c /examples/getuid-user1.yml whoami
daemon
```

This is because the user number 1 happens to be daemon on my system:

```sh
root@02a8cb7164ef:/# head -n2 /etc/passwd 
root:x:0:0:root:/root:/bin/bash
daemon:x:1:1:daemon:/usr/sbin:/usr/sbin/nologin
```

### Change an argument of the call

How about tricking the process to [`openat`](https://linux.die.net/man/2/openat) a different file instead? Easy:

```yaml
# cat /examples/openat-etc-passwd.yml
rules:
  - name: trick the program to read a different file, instead of /etc/passwd
    match:
      name: openat
      args:
        - number: 1
          string: "/etc/passwd"
    modify:
      args:
        - number: 1
          string: "/tmp/passwd"
```

```sh
root@f34cc94a6b6d:/# whoami
root
root@bc2f54570070:/# echo "LOL-HACKED:x:0:0:root:/root:/bin/bash" > /tmp/passwd
root@bc2f54570070:/# monkey -s -c /examples/openat-etc-passwd.yml whoami
LOL-HACKED
```