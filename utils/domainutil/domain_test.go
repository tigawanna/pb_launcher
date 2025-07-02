package domainutil_test

import (
	"pb_launcher/utils/domainutil"
	"testing"
)

func TestIsWildcardDomain(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"*.example.com", true},
		{"example.com", false},
		{"sub.example.com", false},
	}

	for _, tt := range tests {
		result := domainutil.IsWildcardDomain(tt.input)
		if result != tt.expected {
			t.Errorf("IsWildcardDomain(%q) = %v; want %v", tt.input, result, tt.expected)
		}
	}
}

func TestSubdomainMatchesWildcard(t *testing.T) {
	tests := []struct {
		subdomain string
		wildcard  string
		expected  bool
	}{
		{"test.example.com", "*.example.com", true},
		{"example.com", "*.example.com", true},
		{"other.test.example.com", "*.example.com", true},
		{"example.org", "*.example.com", false},
		{"test.example.com", "example.com", false},
	}

	for _, tt := range tests {
		result := domainutil.SubdomainMatchesWildcard(tt.subdomain, tt.wildcard)
		if result != tt.expected {
			t.Errorf("SubdomainMatchesWildcard(%q, %q) = %v; want %v", tt.subdomain, tt.wildcard, result, tt.expected)
		}
	}
}

func TestToWildcardDomain(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"example.com", "*.example.com"},
		{"*.example.com", "*.example.com"},
		{"sub.example.com", "*.sub.example.com"},
		{"*.sub.example.com", "*.sub.example.com"},
	}

	for _, tt := range tests {
		result := domainutil.ToWildcardDomain(tt.input)
		if result != tt.expected {
			t.Errorf("ToWildcardDomain(%q) = %q; want %q", tt.input, result, tt.expected)
		}
	}
}
