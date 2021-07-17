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
