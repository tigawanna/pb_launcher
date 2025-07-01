package mkcert

import (
	"os"
	"os/exec"
	"path/filepath"
	"pb_launcher/internal/certificates/tlscommon"

	"strings"
)

type MkcertProvider struct{}

var _ tlscommon.Provider = (*MkcertProvider)(nil)

func NewMkcertProvider() *MkcertProvider {
	return &MkcertProvider{}
}

func (s *MkcertProvider) RequestCertificate(domain string) (*tlscommon.Certificate, error) {
	tmpDir, err := os.MkdirTemp("", "mkcert")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	certPath := filepath.Join(tmpDir, "cert.pem")
	keyPath := filepath.Join(tmpDir, "key.pem")

	args := []string{"-cert-file", certPath, "-key-file", keyPath}
	if after, ok := strings.CutPrefix(domain, "*."); ok {
		args = append(args, domain, after)
	} else {
		args = append(args, domain)
	}

	cmd := exec.Command("mkcert", args...)
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	keyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	return &tlscommon.Certificate{CertPEM: certPEM, KeyPEM: keyPEM}, nil
}
