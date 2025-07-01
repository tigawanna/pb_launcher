package serialexecutor

import (
	"context"
	"time"
)

type TaskFunc func(ctx context.Context)

type Task struct {
	action   TaskFunc
	interval time.Duration
	priority int
}

func NewTask(task TaskFunc, interval time.Duration, priority int) *Task {
	return &Task{task, interval, priority}
}

func (t *Task) Exec(ctx context.Context, queue chan<- *Task) {
	t.action(ctx)

	go func() {
		timer := time.NewTimer(t.interval)
		defer timer.Stop()

		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			select {
			case queue <- t:
			case <-ctx.Done():
			}
		}
	}()
}
