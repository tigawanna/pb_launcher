package proxy

import (
	"pb_launcher/configs"
	"testing"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/stretchr/testify/assert"
)

type config struct {
	configs.Config
	domain string
}

func (c *config) GetPublicApiDomain() string { return c.domain }
func mockConfig(domain string) configs.Config {
	return &config{domain: domain}
}

func TestExtractServiceID(t *testing.T) {
	rp := NewDynamicReverseProxy(nil, mockConfig("pb.labenv.test"), &apis.ServeConfig{})

	tests := []struct {
		host      string
		expected  string
		expectErr bool
	}{
		{"service1.pb.labenv.test", "service1", false},
		{"sub.service1.pb.labenv.test", "", true},
		{"service-123.pb.labenv.test:8080", "", true},
		{"pb.labenv.test", "", false},
		{"..pb.labenv.test", "", true},
		{".pb.labenv.test", "", true},
		{"invalid.domain.com", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		id, err := rp.extractServiceID(tt.host)
		if tt.expectErr {
			assert.Error(t, err, "expected error for host: %s", tt.host)
		} else {
			assert.NoError(t, err, "unexpected error for host: %s", tt.host)
			assert.Equal(t, tt.expected, id, "incorrect service ID for host: %s", tt.host)
		}
	}
}
