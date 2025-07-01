package serialexecutor

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"
)

type SequentialExecutor struct {
	running bool
	tasks   []*Task
	queue   chan *Task
	ctx     context.Context
	cancel  context.CancelFunc
	mu      sync.Mutex
}

func NewSequentialExecutor() *SequentialExecutor {
	return &SequentialExecutor{
		tasks: make([]*Task, 0),
	}
}

func (s *SequentialExecutor) Add(task *Task) error {
	if task == nil {
		return fmt.Errorf("task must not be nil")
	}
	if task.action == nil {
		return fmt.Errorf("task action must not be nil")
	}
	if task.interval <= 0 {
		return fmt.Errorf("task interval must be greater than zero")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.running {
		return errors.New("cannot add tasks while running")
	}

	s.tasks = append(s.tasks, task)
	return nil
}

func (s *SequentialExecutor) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return nil
	}
	slices.SortFunc(s.tasks, func(a, b *Task) int {
		return b.priority - a.priority
	})
	s.queue = make(chan *Task, len(s.tasks))
	for _, t := range s.tasks {
		s.queue <- t
	}

	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.running = true

	go func() {
		for {
			select {
			case <-s.ctx.Done():
				return
			case task, ok := <-s.queue:
				if !ok {
					return
				}
				task.Exec(s.ctx, s.queue)
			}
		}
	}()

	return nil
}

func (s *SequentialExecutor) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return nil
	}

	s.running = false
	s.cancel()
	close(s.queue)
	return nil
}
