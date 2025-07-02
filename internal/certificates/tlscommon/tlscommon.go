package tlscommon

import (
	"errors"
	"time"
)

type Certificate struct {
	CertPEM []byte
	KeyPEM  []byte
	// TTL indicates how much time remains before the certificate expires.
	Ttl time.Duration
}

var ErrUnsupportedProvider = errors.New("unsupported certificate provider")

type Provider interface {
	RequestCertificate(domain string) (*Certificate, error)
}
