package serialexecutor

import (
	"context"
	"testing"
	"time"
)

func TestTaskExec(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	executed := make(chan bool, 1)
	queue := make(chan *Task, 1)

	task := &Task{
		action: func(ctx context.Context) {
			executed <- true
		},
		interval: 100 * time.Millisecond,
	}

	task.Exec(ctx, queue)

	select {
	case <-executed:
		// ok
	case <-time.After(500 * time.Millisecond):
		t.Fatal("action was not executed")
	}

	select {
	case <-queue:
		// ok
	case <-time.After(500 * time.Millisecond):
		t.Fatal("task was not re-queued after interval")
	}
}

func TestTaskExecWithInterval(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	queue := make(chan *Task, 1)
	executed := make(chan time.Time, 2)

	task := &Task{
		action: func(ctx context.Context) {
			executed <- time.Now()
		},
		interval: 200 * time.Millisecond,
	}

	start := time.Now()
	task.Exec(ctx, queue)

	select {
	case t1 := <-executed:
		elapsed1 := t1.Sub(start)
		if elapsed1 > 50*time.Millisecond {
			t.Fatalf("first execution should be immediate, got delay %v", elapsed1)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("first execution did not happen")
	}

	select {
	case <-queue:
		// task was re-queued correctly
	case <-time.After(500 * time.Millisecond):
		t.Fatal("task was not re-queued")
	}

	// Re-execute the task manually to test interval behavior
	task.Exec(ctx, queue)

	select {
	case <-executed:
		// second execution completed
	case <-time.After(500 * time.Millisecond):
		t.Fatal("second execution did not happen")
	}

	// Wait slightly longer than the interval to confirm the task gets re-queued
	time.Sleep(220 * time.Millisecond)

	select {
	case <-queue:
		// task was queued again after interval
	default:
		t.Fatal("task should have been queued again after interval")
	}
}
