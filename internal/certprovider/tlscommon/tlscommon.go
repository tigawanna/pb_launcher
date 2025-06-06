package tlscommon

import "errors"

type Certificate struct {
	CertPEM []byte
	KeyPEM  []byte
}

var ErrUnsupportedProvider = errors.New("unsupported certificate provider")

type Provider interface {
	RequestCertificate(domain string) (Certificate, error)
}
