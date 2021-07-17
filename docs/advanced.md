## Custom logic in Go

Sometimes, you need to build something custom. Syscall Monkey has got you covered!


### Implementing a custom manipulator

It's easy to build a custom binary with a manipulator written in Go.

In order to do that, you need to implement the [`SyscallManipulator` interface](https://github.com/syscallmonkey/monkey/blob/main/pkg/syscall/manipulator.go), and pass an instance to the `RunTracer` function.

It consists of two functions:

- `HandleEntry` that's called before a syscall is about to be executed - this is where you can modify the arguments, the syscall code, block the call entirely etc.
- `HandleExit` which is where you can modify the return code

For every sycall that Syscall Monkey tracer, it will call the two functions in order, and apply any modifications that you request.

For [example](/examples/example-sdk-usage/example-sdk.go), this silly manipulator will mess the program up by changing all syscalls to [`getpid`](https://man7.org/linux/man-pages/man2/getpid.2.html) and also chaging the return value on every other call:

```go
package main

import (
	"os"

	"github.com/syscallmonkey/monkey/pkg/config"
	"github.com/syscallmonkey/monkey/pkg/run"
	"github.com/syscallmonkey/monkey/pkg/syscall"
)

// ExampleManipulator does some random stuff, to illustrate what you can do
type ExampleManipulator struct {
	Count int
}

func (sm *ExampleManipulator) HandleEntry(state syscall.SyscallState) syscall.SyscallState {
	// change syscall to always be getpid
	state.SyscallCode = 102
	// and also count the entries
	sm.Count++
	return state
}

func (sm *ExampleManipulator) HandleExit(returnValue uint64) uint64 {
	// change the syscall return value on every other call
	if sm.Count%2 == 0 {
		return 0
	}
	return returnValue
}

func main() {
	// parse the config (or hardcode them, if you'd like)
	config := config.ParseCommandLineFlags(os.Args[1:])
	// implement your manipulator
	m := ExampleManipulator{}
	// run the tracer
	run.RunTracer(config, &m)
}
```

You can hardcode the config, or you can use `config.ParseCommandLineFlags(os.Args[1:])` to inherit all the other flags that the regular Syscall Monkey supports.