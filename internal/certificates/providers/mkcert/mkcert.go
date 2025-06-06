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

func (s *MkcertProvider) RequestCertificate(domain string) (tlscommon.Certificate, error) {
	tmpDir, err := os.MkdirTemp("", "mkcert")
	if err != nil {
		return tlscommon.Certificate{}, err
	}
	defer os.RemoveAll(tmpDir)

	args := []string{}
	if after, ok := strings.CutPrefix(domain, "*."); ok {
		args = append(args, domain, after)
	} else {
		args = append(args, domain)
	}

	cmd := exec.Command("mkcert", append(args, "-cert-file", "cert.pem", "-key-file", "key.pem")...)
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		return tlscommon.Certificate{}, err
	}

	certPEM, err := os.ReadFile(filepath.Join(tmpDir, "cert.pem"))
	if err != nil {
		return tlscommon.Certificate{}, err
	}

	keyPEM, err := os.ReadFile(filepath.Join(tmpDir, "key.pem"))
	if err != nil {
		return tlscommon.Certificate{}, err
	}

	return tlscommon.Certificate{CertPEM: certPEM, KeyPEM: keyPEM}, nil
}
