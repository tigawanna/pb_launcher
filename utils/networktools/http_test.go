package networktools_test

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"net/url"
	"pb_launcher/utils/networktools"
	"testing"
)

func TestPrepareProxyHeaders(t *testing.T) {
	tests := []struct {
		name              string
		remoteAddr        string
		initialHeaders    map[string]string
		isSecure          bool
		expectedRealIP    string
		expectedForwarded string
		expectedProto     string
		expectedHost      string
	}{
		{
			name:              "Without existing headers, insecure request",
			remoteAddr:        "192.0.2.1:1234",
			isSecure:          false,
			expectedRealIP:    "192.0.2.1",
			expectedForwarded: "192.0.2.1",
			expectedProto:     "http",
			expectedHost:      "target.example.com",
		},
		{
			name:              "Without existing headers, secure request",
			remoteAddr:        "198.51.100.2:5678",
			isSecure:          true,
			expectedRealIP:    "198.51.100.2",
			expectedForwarded: "198.51.100.2",
			expectedProto:     "https",
			expectedHost:      "target.example.com",
		},
		{
			name:              "With existing X-Forwarded-For",
			remoteAddr:        "203.0.113.5:9999",
			initialHeaders:    map[string]string{"X-Forwarded-For": "1.1.1.1, 2.2.2.2"},
			isSecure:          false,
			expectedRealIP:    "203.0.113.5",
			expectedForwarded: "1.1.1.1, 2.2.2.2, 203.0.113.5",
			expectedProto:     "http",
			expectedHost:      "target.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://original.example.com", nil)
			req.RemoteAddr = tt.remoteAddr

			if tt.isSecure {
				req.TLS = &tls.ConnectionState{}
			}

			for k, v := range tt.initialHeaders {
				req.Header.Set(k, v)
			}

			targetURL, _ := url.Parse("http://target.example.com")

			networktools.PrepareProxyHeaders(req, targetURL)

			if req.Header.Get("X-Real-IP") != tt.expectedRealIP {
				t.Errorf("Expected X-Real-IP %s, got %s", tt.expectedRealIP, req.Header.Get("X-Real-IP"))
			}

			if req.Header.Get("X-Forwarded-For") != tt.expectedForwarded {
				t.Errorf("Expected X-Forwarded-For %s, got %s", tt.expectedForwarded, req.Header.Get("X-Forwarded-For"))
			}

			if req.Header.Get("X-Forwarded-Proto") != tt.expectedProto {
				t.Errorf("Expected X-Forwarded-Proto %s, got %s", tt.expectedProto, req.Header.Get("X-Forwarded-Proto"))
			}

			if req.Header.Get("X-Forwarded-Host") != tt.expectedHost {
				t.Errorf("Expected X-Forwarded-Host %s, got %s", tt.expectedHost, req.Header.Get("X-Forwarded-Host"))
			}

			if req.Host != tt.expectedHost {
				t.Errorf("Expected req.Host %s, got %s", tt.expectedHost, req.Host)
			}
		})
	}
}

func TestIsRequestSecure(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		forwardedProto string
		expectSecure   bool
	}{
		{
			name:         "Request with TLS",
			url:          "https://example.com",
			expectSecure: true,
		},
		{
			name:         "Request without TLS",
			url:          "http://example.com",
			expectSecure: false,
		},
		{
			name:           "Request with X-Forwarded-Proto https",
			url:            "http://example.com",
			forwardedProto: "https",
			expectSecure:   true,
		},
		{
			name:           "Request with X-Forwarded-Proto http",
			url:            "http://example.com",
			forwardedProto: "http",
			expectSecure:   false,
		},
		{
			name:         "Request without TLS or X-Forwarded-Proto",
			url:          "http://example.com",
			expectSecure: false,
		},
		{
			name:           "X-Forwarded-Proto with multiple values",
			url:            "http://example.com",
			forwardedProto: "https, http",
			expectSecure:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			if tt.forwardedProto != "" {
				req.Header.Set("X-Forwarded-Proto", tt.forwardedProto)
			}

			secure := networktools.IsRequestSecure(req)
			if secure != tt.expectSecure {
				t.Errorf("Expected secure=%v, got %v", tt.expectSecure, secure)
			}
		})
	}
}
func TestBuildHostURL(t *testing.T) {
	tests := []struct {
		name     string
		schema   string
		host     string
		port     string
		paths    []string
		expected string
	}{
		{
			name:     "Simple host without port and without paths",
			schema:   "http",
			host:     "example.com",
			port:     "",
			paths:    nil,
			expected: "http://example.com",
		},
		{
			name:     "Simple host with port 8080 and without paths",
			schema:   "http",
			host:     "example.com",
			port:     "8080",
			paths:    nil,
			expected: "http://example.com:8080",
		},
		{
			name:     "Host with port included, override port",
			schema:   "http",
			host:     "example.com:3000",
			port:     "9090",
			paths:    nil,
			expected: "http://example.com:9090",
		},
		{
			name:     "Host with port included, port 80 omitted",
			schema:   "http",
			host:     "example.com:3000",
			port:     "80",
			paths:    nil,
			expected: "http://example.com",
		},
		{
			name:     "Full URL as host, extract only host",
			schema:   "https",
			host:     "http://example.com:4000/some/path",
			port:     "443",
			paths:    nil,
			expected: "https://example.com",
		},
		{
			name:     "Full URL as host, override port and add paths",
			schema:   "https",
			host:     "http://example.com:4000/some/path",
			port:     "8443",
			paths:    []string{"api", "v1"},
			expected: "https://example.com:8443/api/v1",
		},
		{
			name:     "Paths with mixed slashes",
			schema:   "http",
			host:     "example.com",
			port:     "",
			paths:    []string{"/one/", "two", "/three/"},
			expected: "http://example.com/one/two/three",
		},
		{
			name:     "Host and paths with multiple slashes",
			schema:   "https",
			host:     "example.com///",
			port:     "",
			paths:    []string{"///api///", "///v1///", "///test///"},
			expected: "https://example.com/api/v1/test",
		},
		{
			name:     "Host with port, multiple slashes, and custom port",
			schema:   "https",
			host:     "example.com:5000///",
			port:     "8443",
			paths:    []string{"///api///", "///v1///", "///test///"},
			expected: "https://example.com:8443/api/v1/test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := networktools.BuildHostURL(tt.schema, tt.host, tt.port, tt.paths...)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
