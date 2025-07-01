package process

import (
	"testing"
	"time"
)

// TestProcess_CloseChanIsClosed verifies that after stopping a process,
// the internal closeChan is properly closed, ensuring the process cleanup completes
func TestProcess_CloseChanIsClosed(t *testing.T) {
	service := New("test-service", "sleep", []string{"1"})

	if err := service.Start(); err != nil {
		t.Fatalf("failed to start process: %v", err)
	}

	if service.closeChan == nil {
		t.Fatalf("closeChan should be initialized after Start()")
	}

	if err := service.Stop(); err != nil {
		t.Fatalf("failed to stop process: %v", err)
	}

	select {
	case <-service.closeChan:
		// OK, channel closed as expected
	default:
		t.Fatalf("closeChan should be closed after Stop()")
	}
}

// TestProcess_ForceKillOnTimeout simulates an unresponsive process by blocking closeChan
// and ensures that after the 10-second timeout, the process is forcefully terminated with SIGKILL

func TestProcess_ForceKillOnTimeout(t *testing.T) {
	service := New("test-service", "sleep", []string{"60"})

	if err := service.Start(); err != nil {
		t.Fatalf("failed to start process: %v", err)
	}

	cmd := service.h.currentCommand()
	if cmd == nil {
		t.Fatalf("expected active command after Start()")
	}

	// Replace closeChan to block indefinitely, simulating unresponsive process
	blockChan := make(chan struct{})
	service.closeChan = blockChan

	done := make(chan struct{})
	go func() {
		if err := service.Stop(); err != nil {
			t.Errorf("stop returned error: %v", err)
		}
		close(done)
	}()

	select {
	case <-done:
		// OK, Stop() completed after force kill
	case <-time.After(12 * time.Second):
		t.Fatalf("Stop() did not complete within expected timeout")
	}
}
