package process

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"syscall"
	"time"
)

type ProcessErrorMessage struct {
	ID    string
	Error error
}

type ProcessOptions struct {
	errChan chan<- ProcessErrorMessage
}

type ProcessOption = func(*ProcessOptions)

func WithErrorChan(errChan chan<- ProcessErrorMessage) ProcessOption {
	return func(options *ProcessOptions) { options.errChan = errChan }
}

type Process struct {
	id      string
	options *ProcessOptions

	command string
	args    []string

	h *handler

	closeChan chan struct{}
}

func New(ID string, command string, args []string, options ...ProcessOption) *Process {
	p := &Process{
		id:      ID,
		h:       &handler{status: Stopped},
		command: command,
		args:    args,
		options: &ProcessOptions{},
	}

	for _, applay := range options {
		applay(p.options)
	}

	return p
}

func (p *Process) Status() ProcessState { return p.h.currentState() }
func (p *Process) IsRunning() bool      { return p.Status() == Running }

func (p *Process) Start() error {
	currentState := p.Status()
	if currentState != Stopped {
		return nil
	}

	cmd := exec.Command(p.command, p.args...)
	p.h.updateStatus(Starting)

	if err := cmd.Start(); err != nil {
		p.h.updateStatus(Stopped)
		slog.Error("failed to start process", "error", err, "process_id", p.id)
		return err
	}

	p.closeChan = make(chan struct{})
	go p.waitForExit(cmd, p.closeChan)

	p.h.replaceCommand(cmd)
	p.h.updateStatus(Running)
	return nil
}

func (p *Process) waitForExit(cmd *exec.Cmd, doneChan chan struct{}) {
	if err := cmd.Wait(); err != nil {
		if err.Error() != "signal: terminated" {
			if p.options.errChan != nil {
				p.options.errChan <- ProcessErrorMessage{
					ID:    p.id,
					Error: fmt.Errorf("process exited with error: %w", err),
				}
			}
			slog.Error("process exited with error", "error", err, "process_id", p.id)
		}
	}
	p.h.updateStatus(Stopped)
	if doneChan != nil {
		close(doneChan)
	}
}

func (p *Process) Stop() error {
	currentState := p.Status()
	if currentState != Running {
		return nil
	}

	cmd := p.h.currentCommand()
	if cmd == nil {
		slog.Warn("stop ignored: no active command found", "process_id", p.id)
		return nil
	}

	p.h.updateStatus(Stopping)
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		p.h.updateStatus(Running)
		slog.Error("failed to stop process", "error", err, "process_id", p.id)
		return err
	}
	// brief pause to ensure waitForExit completes status update to Stopped
	// time.Sleep(100 * time.Millisecond)
	if p.closeChan != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		select {
		case <-p.closeChan:
		case <-ctx.Done():
		}
	}
	return nil
}

func (p *Process) Restart() error {
	currentState := p.Status()
	if currentState != Running && currentState != Stopped {
		return fmt.Errorf("restart failed: invalid process state (%v)", currentState)
	}
	if currentState == Running {
		if err := p.Stop(); err != nil {
			slog.Error("restart failed at stop", "error", err, "process_id", p.id)
			return fmt.Errorf("restart failed at stop: %w", err)
		}
	}
	if err := p.Start(); err != nil {
		slog.Error("restart failed at start", "error", err, "process_id", p.id)
		return fmt.Errorf("restart failed at start: %w", err)
	}
	return nil
}
