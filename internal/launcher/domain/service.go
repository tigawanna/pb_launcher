package domain

import (
	"context"
	"fmt"
	"os/exec"
	"sync"
)

type Service struct {
	ID      string
	mu      sync.Mutex
	cmd     *exec.Cmd
	errChan chan<- ServiceError
}

type ServiceError struct {
	ID    string
	Error error
}

func NewService(id string, command string, args []string, errChan chan<- ServiceError) *Service {
	return &Service{
		ID:      id,
		cmd:     exec.Command(command, args...),
		errChan: errChan,
	}
}

func (s *Service) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cmd.Process != nil {
		return fmt.Errorf("service already running")
	}

	s.cmd = exec.CommandContext(ctx, s.cmd.Path, s.cmd.Args[1:]...)
	err := s.cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		err := s.cmd.Wait()
		s.mu.Lock()
		s.cmd = nil
		s.mu.Unlock()

		if err != nil {
			s.errChan <- ServiceError{ID: s.ID, Error: err}
		}
	}()

	return nil
}

func (s *Service) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cmd == nil || s.cmd.Process == nil {
		return fmt.Errorf("service not running")
	}

	err := s.cmd.Process.Kill()
	if err != nil {
		return err
	}

	s.cmd = nil
	return nil
}

func (s *Service) Restart(ctx context.Context) error {
	if err := s.Stop(); err != nil {
		return err
	}
	return s.Start(ctx)
}

func (s *Service) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.cmd != nil && s.cmd.Process != nil
}
