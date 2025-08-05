package proxy

import (
	"testing"
)

func TestExtractID(t *testing.T) {
	rp := &DynamicReverseProxyDiscovery{apiDomain: "example.com"}

	tests := []struct {
		host     string
		expected string
		wantErr  bool
		errMsg   string
	}{
		{"service1.example.com", "service1", false, ""},
		{"service-abc.example.com", "service-abc", false, ""},
		{"service1.sub.example.com", "", true, "invalid ID: prefix contains invalid character '.'"},
		{"example.com", "", true, "invalid ID: host is the base domain"},
		{"otherdomain.com", "", false, ""},
		{"", "", false, ""},
	}

	for _, tt := range tests {
		id, err := rp.extractID(tt.host)
		if tt.wantErr {
			if err == nil || err.Error() != tt.errMsg {
				t.Errorf("extractID(%q) error = %v, want %v", tt.host, err, tt.errMsg)
			}
		} else {
			if err != nil {
				t.Errorf("extractID(%q) unexpected error: %v", tt.host, err)
			}
			if id != tt.expected {
				t.Errorf("extractID(%q) = %v, want %v", tt.host, id, tt.expected)
			}
		}
	}
}
