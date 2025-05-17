package process

import (
	"os/exec"
	"sync"
)

type ProcessState int

const (
	Stopped ProcessState = iota
	Starting
	Running
	Stopping
)

type handler struct {
	sync.RWMutex
	status ProcessState
	cmd    *exec.Cmd
}

func (_c *handler) updateStatus(status ProcessState) {
	_c.Lock()
	defer _c.Unlock()
	_c.status = status
}

func (_c *handler) currentState() ProcessState {
	_c.RLock()
	defer _c.RUnlock()
	return _c.status
}

func (_c *handler) replaceCommand(c *exec.Cmd) {
	_c.Lock()
	defer _c.Unlock()
	_c.cmd = nil
	_c.cmd = c
}

func (_c *handler) currentCommand() *exec.Cmd {
	_c.RLock()
	defer _c.RUnlock()
	return _c.cmd
}
