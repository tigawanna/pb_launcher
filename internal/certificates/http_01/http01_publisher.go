package http01

import (
	"fmt"
	"net"
	"net/url"
	"sync"
)

// Http01ChallengeAddressPublisher publishes the IP and port where the HTTP-01 challenge server runs,
// so the reverse proxy can route Let's Encrypt traffic to the correct service.
type Http01ChallengeAddressPublisher struct {
	ip   string
	port string
	mu   sync.RWMutex
}

func NewHttp01ChallengeAddressPublisher() *Http01ChallengeAddressPublisher {
	return &Http01ChallengeAddressPublisher{}
}

func (h *Http01ChallengeAddressPublisher) Publish(ip, port string) error {
	if net.ParseIP(ip) == nil {
		return fmt.Errorf("invalid IP address: %s", ip)
	}
	if port == "" {
		return fmt.Errorf("port cannot be empty")
	}
	h.mu.Lock()
	defer h.mu.Unlock()

	h.ip = ip
	h.port = port
	return nil
}

func (h *Http01ChallengeAddressPublisher) ResolveAddress() (*url.URL, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if h.ip == "" || h.port == "" {
		return nil, fmt.Errorf("IP or port not published")
	}
	u := &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%s", h.ip, h.port),
	}
	return u, nil
}
