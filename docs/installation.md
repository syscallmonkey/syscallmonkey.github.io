## Installation

### Binary

You can build the binary from the source:

```sh
git clone https://github.com/syscallmonkey/monkey.git
cd monkey
make bin/monkey
./bin/monkey -h
```

### Use Docker container

```sh
docker pull ghcr.io/syscallmonkey/monkey:0.0.1rc1
```

Check [the latest available versions here](https://github.com/syscallmonkey/monkey/pkgs/container/monkey).


### Building Docker container

```sh
git clone https://github.com/syscallmonkey/monkey.git
cd monkey
make build
make run

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