package tlscommon

import (
	"errors"
)

var (
	ErrCertificateNotFound = errors.New("certificate not found")
	ErrInvalidPEM          = errors.New("invalid certificate PEM format")
	ErrCertificateExpired  = errors.New("certificate is expired or not yet valid")
)

type Store interface {
	Store(domain string, cert Certificate) error
	Resolve(domain string) (*Certificate, error)
}
