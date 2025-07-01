package serialexecutor

import (
	"context"
	"testing"
	"time"
)

func TestSequentialExecutorBasic(t *testing.T) {
	exec := NewSequentialExecutor()

	executions := make(chan time.Time, 2)
	task := &Task{
		action: func(ctx context.Context) {
			executions <- time.Now()
		},
		interval: 100 * time.Millisecond,
	}

	if err := exec.Add(task); err != nil {
		t.Fatalf("failed to add task: %v", err)
	}

	if err := exec.Start(); err != nil {
		t.Fatalf("failed to start executor: %v", err)
	}

	select {
	case <-executions:
		// first execution ok
	case <-time.After(500 * time.Millisecond):
		t.Fatal("first execution timeout")
	}

	select {
	case <-executions:
		// re-executed after interval
	case <-time.After(500 * time.Millisecond):
		t.Fatal("second execution timeout")
	}

	exec.Stop()
}

func TestSequentialExecutorStopCancels(t *testing.T) {
	exec := NewSequentialExecutor()
	taskExecuted := make(chan bool, 1)

	task := &Task{
		action: func(ctx context.Context) {
			taskExecuted <- true
		},
		interval: 100 * time.Millisecond,
	}

	exec.Add(task)
	exec.Start()

	select {
	case <-taskExecuted:
		// executed ok
	case <-time.After(500 * time.Millisecond):
		t.Fatal("task did not execute")
	}

	exec.Stop()

	time.Sleep(150 * time.Millisecond)

	select {
	case <-taskExecuted:
		t.Fatal("task re-executed after stop")
	default:
		// ok, no further execution
	}
}

func TestSequentialExecutorOrdenEjecucion(t *testing.T) {
	exec := NewSequentialExecutor()
	taskExecuted := make(chan int, 3)

	task1 := &Task{
		action: func(ctx context.Context) {
			taskExecuted <- 3
		},
		interval: 100 * time.Millisecond,
		priority: 1,
	}
	task2 := &Task{
		action: func(ctx context.Context) {
			taskExecuted <- 2
		},
		interval: 100 * time.Millisecond,
		priority: 2,
	}
	task3 := &Task{
		action: func(ctx context.Context) {
			taskExecuted <- 1
		},
		interval: 100 * time.Millisecond,
		priority: 3,
	}

	exec.Add(task2)
	exec.Add(task1)
	exec.Add(task3)

	if err := exec.Start(); err != nil {
		t.Error(err)
	}

	for i := 1; i <= 3; i++ {
		select {
		case num := <-taskExecuted:
			if num != i {
				t.Errorf("expected %d, got %d", i, num)
			}
		case <-time.After(500 * time.Millisecond):
			t.Errorf("timeout waiting for task %d", i)
		}
	}
}

func TestSequentialExecutorCheckSimpleIntervals(t *testing.T) {
	exec := NewSequentialExecutor()

	executions := make([]time.Time, 0)
	task := &Task{
		action: func(ctx context.Context) {
			start := time.Now()
			time.Sleep(50 * time.Millisecond)
			end := time.Now()
			go func() {
				executions = append(executions, start, end)
			}()
		},
		interval: 100 * time.Millisecond,
	}

	if err := exec.Add(task); err != nil {
		t.Fatalf("failed to add task: %v", err)
	}
	startTime := time.Now()
	if err := exec.Start(); err != nil {
		t.Fatalf("failed to start executor: %v", err)
	}
	time.Sleep(450 * time.Millisecond)
	exec.Stop()

	points := make([]int64, len(executions))
	for i, e := range executions {
		points[i] = e.Sub(startTime).Milliseconds() / 10
	}
	expected := []int64{0, 5, 15, 20, 30, 35}

	for i := range points {
		if len(points) != len(expected) {
			t.Fatalf("array lengths differ: expected %d, got %d\nexpected: %v\ngot: %v", len(expected), len(points), expected, points)
		}
		if points[i] != expected[i] {
			t.Fatalf("arrays differ\nexpected: %v\ngot: %v", expected, points)
		}
	}
}

func TestSequentialExecutorCheckSimpleIntervalsT2(t *testing.T) {
	exec := NewSequentialExecutor()

	executions := make([]time.Time, 0)
	task := &Task{
		action: func(ctx context.Context) {
			start := time.Now()
			time.Sleep(80 * time.Millisecond)
			end := time.Now()
			go func() {
				executions = append(executions, start, end)
			}()
		},
		interval: 200 * time.Millisecond,
	}

	if err := exec.Add(task); err != nil {
		t.Fatalf("failed to add task: %v", err)
	}
	startTime := time.Now()
	if err := exec.Start(); err != nil {
		t.Fatalf("failed to start executor: %v", err)
	}
	time.Sleep(660 * time.Millisecond)
	exec.Stop()

	points := make([]int64, len(executions))
	for i, e := range executions {
		points[i] = e.Sub(startTime).Milliseconds() / 10
	}

	expected := []int64{0, 8, 28, 36, 56, 64}

	for i := range points {
		if len(points) != len(expected) {
			t.Fatalf("array lengths differ: expected %d, got %d\nexpected: %v\ngot: %v", len(expected), len(points), expected, points)
		}
		if points[i] != expected[i] {
			t.Fatalf("arrays differ\nexpected: %v\ngot: %v", expected, points)
		}
	}
}

func TestSequentialExecutorCheckIntervalsT1_T2(t *testing.T) {
	exec := NewSequentialExecutor()
	executions01 := make([]time.Time, 0)
	task01 := &Task{
		action: func(ctx context.Context) {
			start := time.Now()
			time.Sleep(50 * time.Millisecond)
			end := time.Now()
			go func() {
				executions01 = append(executions01, start, end)
			}()
		},
		interval: 100 * time.Millisecond,
	}
	executions02 := make([]time.Time, 0)
	task02 := &Task{
		action: func(ctx context.Context) {
			start := time.Now()
			time.Sleep(80 * time.Millisecond)
			end := time.Now()
			go func() {
				executions02 = append(executions02, start, end)
			}()
		},
		interval: 200 * time.Millisecond,
	}

	if err := exec.Add(task01); err != nil {
		t.Fatalf("failed to add task: %v", err)
	}
	if err := exec.Add(task02); err != nil {
		t.Fatalf("failed to add task: %v", err)
	}

	startTime := time.Now()
	if err := exec.Start(); err != nil {
		t.Fatalf("failed to start executor: %v", err)
	}
	time.Sleep(750 * time.Millisecond)
	exec.Stop()

	points01 := make([]int64, len(executions01))
	for i, e := range executions01 {
		points01[i] = e.Sub(startTime).Milliseconds() / 10
	}

	points02 := make([]int64, len(executions02))
	for i, e := range executions02 {
		points02[i] = e.Sub(startTime).Milliseconds() / 10
	}
	expected01 := []int64{0, 5, 15, 20, 30, 35, 45, 50, 60, 65}
	expected02 := []int64{5, 13, 35, 43, 65, 73}

	if len(points01) != len(expected01) {
		t.Fatalf("array lengths differ: expected %d, got %d\nexpected: %v\ngot: %v", len(expected01), len(points01), expected01, points01)
	}

	for i := range points01 {
		if points01[i] != expected01[i] {
			t.Fatalf("arrays differ\nexpected: %v\ngot: %v", expected01, points01)
		}
	}
	if len(points02) != len(expected02) {
		t.Fatalf("array lengths differ: expected %d, got %d\nexpected: %v\ngot: %v", len(expected02), len(points02), expected02, points02)
	}

	for i := range points02 {
		if points02[i] != expected02[i] {
			t.Fatalf("arrays differ\nexpected: %v\ngot: %v", expected02, points02)
		}
	}
}
