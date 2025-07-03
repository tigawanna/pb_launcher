package networktools_test

import (
	"fmt"
	"net"
	"pb_launcher/utils/networktools"
	"testing"
)

func TestGetAvailablePort(t *testing.T) {
	ip := "127.0.0.1"

	gotIP, port, err := networktools.GetAvailablePort(ip)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if gotIP != ip {
		t.Errorf("Expected IP %s, got %s", ip, gotIP)
	}

	if port <= 0 || port > 65535 {
		t.Errorf("Expected valid port, got %d", port)
	}

	// Optional: Try binding manually to confirm availability
	addr := net.JoinHostPort(ip, fmt.Sprintf("%d", port))
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		t.Errorf("Expected to bind to %s, but got error: %v", addr, err)
	} else {
		listener.Close()
	}
}

func TestGetAvailablePort_InvalidIP(t *testing.T) {
	_, _, err := networktools.GetAvailablePort("invalid-ip")
	if err == nil {
		t.Error("Expected error for invalid IP, got nil")
	}
}
