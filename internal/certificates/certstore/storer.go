package certstore

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"pb_launcher/configs"
	"pb_launcher/internal/certificates/tlscommon"

	"time"
)

const (
	cert_file_name        = "certificate.pem"
	private_key_file_name = "private_key.pem"
)
const FolderDateFormat = "2006-01-02_15-04-05"

type TlsStorer struct {
	rootPath string
}

var _ tlscommon.Store = (*TlsStorer)(nil)

func NewTlsStorer(c configs.Config) *TlsStorer {
	return &TlsStorer{
		rootPath: c.GetCertificatesDir(),
	}
}

func (s *TlsStorer) Store(domain string, cert tlscommon.Certificate) error {
	timestamp := time.Now().Format(FolderDateFormat)
	outputDir := filepath.Join(s.rootPath, domain, timestamp)

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return err
	}

	certFile := filepath.Join(outputDir, cert_file_name)
	keyFile := filepath.Join(outputDir, private_key_file_name)

	if err := os.WriteFile(certFile, cert.CertPEM, 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(keyFile, cert.KeyPEM, 0o600); err != nil {
		return err
	}

	return nil
}

func (s *TlsStorer) IsCertificateValid(cert *tlscommon.Certificate) error {
	block, _ := pem.Decode(cert.CertPEM)
	if block == nil {
		return tlscommon.ErrInvalidPEM
	}

	parsedCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}

	now := time.Now()
	if now.Before(parsedCert.NotBefore) || now.After(parsedCert.NotAfter) {
		return tlscommon.ErrCertificateExpired
	}

	ttl := time.Until(parsedCert.NotAfter)
	if ttl < 0 {
		ttl = 0
	}
	cert.Ttl = ttl

	return nil
}

func (s *TlsStorer) Resolve(domain string) (*tlscommon.Certificate, error) {
	domainPath := filepath.Join(s.rootPath, domain)

	entries, err := os.ReadDir(domainPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w: %s",
				tlscommon.ErrCertificateNotFound,
				domain,
			)
		}
		return nil, err
	}

	var latestDir string
	for _, entry := range entries {
		if entry.IsDir() {
			name := entry.Name()
			if latestDir == "" || name > latestDir {
				latestDir = name
			}
		}
	}

	if latestDir == "" {
		return nil, fmt.Errorf("%w: %s",
			tlscommon.ErrCertificateNotFound,
			domain,
		)
	}

	certDir := filepath.Join(domainPath, latestDir)
	certFile := filepath.Join(certDir, cert_file_name)
	keyFile := filepath.Join(certDir, private_key_file_name)

	certPEM, err := os.ReadFile(certFile)
	if err != nil {
		return nil, err
	}

	keyPEM, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	cert := tlscommon.Certificate{
		CertPEM: certPEM,
		KeyPEM:  keyPEM,
	}
	if err := s.IsCertificateValid(&cert); err != nil {
		return nil, err
	}
	return &cert, nil
}
