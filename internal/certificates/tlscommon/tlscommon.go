package tlscommon

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"
)

type Certificate struct {
	CertPEM []byte
	KeyPEM  []byte
	// TTL indicates how much time remains before the certificate expires.
	ttl time.Duration
}

func (c *Certificate) GetTTL() time.Duration {
	if c.ttl > 0 {
		return c.ttl
	}
	block, _ := pem.Decode(c.CertPEM)
	if block == nil || block.Type != "CERTIFICATE" {
		return 0
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return 0
	}
	c.ttl = time.Until(cert.NotAfter)
	return c.ttl
}

var ErrUnsupportedProvider = errors.New("unsupported certificate provider")

type Provider interface {
	RequestCertificate(domain string) (*Certificate, error)
}
