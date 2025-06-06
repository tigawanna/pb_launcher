package selfsigned

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"pb_launcher/internal/certificates/tlscommon"

	"strings"
	"time"
)

type SelfSignedProvider struct{}

var _ tlscommon.Provider = (*SelfSignedProvider)(nil)

func NewSelfSignedProvider() *SelfSignedProvider {
	return &SelfSignedProvider{}
}

func (s *SelfSignedProvider) RequestCertificate(domain string) (tlscommon.Certificate, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tlscommon.Certificate{}, err
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return tlscommon.Certificate{}, err
	}

	dnsNames := []string{domain}
	if after, ok := strings.CutPrefix(domain, "*."); ok {
		base := after
		dnsNames = append(dnsNames, base)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      pkix.Name{CommonName: domain},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(1, 0, 0),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     dnsNames,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return tlscommon.Certificate{}, err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return tlscommon.Certificate{}, err
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes})

	return tlscommon.Certificate{CertPEM: certPEM, KeyPEM: keyPEM}, nil
}
