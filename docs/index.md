# Welcome to Syscall Monkey

## TL;DR

`Syscall Monkey` is like [`strace`](https://man7.org/linux/man-pages/man1/strace.1.html) for fiddling:

- attach and detach processes using [`ptrace`](https://man7.org/linux/man-pages/man2/ptrace.2.html) (Linux only)
- trace their [`syscalls`](https://man7.org/linux/man-pages/man2/syscalls.2.html) - names, arguments, return values
- manipulate [`syscalls`](https://man7.org/linux/man-pages/man2/syscalls.2.html) (block, change arguments, return value) to simulate failure
- prepare scenarios in a simple `yaml` format
- write custom scenarios using `syscallmonkey` as an SDK

