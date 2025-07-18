package http01

import (
	"errors"
	"log/slog"
	"pb_launcher/configs"
	"pb_launcher/internal/certificates/tlscommon"
	"pb_launcher/utils/domainutil"
	"pb_launcher/utils/networktools"
	"strconv"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/http01"
)

type HTTP01TLSCertificateRequestService struct {
	acmeEmail      string
	ipAddress      string
	publicher      *Http01ChallengeAddressPublisher
	clientProvider *tlscommon.LetsEncryptClientAccountProvider
}

var _ tlscommon.Provider = (*HTTP01TLSCertificateRequestService)(nil)

func NewHTTP01TLSCertificateRequestService(
	publicher *Http01ChallengeAddressPublisher,
	clientProvider *tlscommon.LetsEncryptClientAccountProvider,
	c configs.Config,
) *HTTP01TLSCertificateRequestService {
	return &HTTP01TLSCertificateRequestService{
		publicher:      publicher,
		clientProvider: clientProvider,
		ipAddress:      c.GetListenIPAddress(),
		acmeEmail:      c.GetAcmeEmail(),
	}
}

// RequestCertificate implements tlscommon.Provider.
func (h *HTTP01TLSCertificateRequestService) RequestCertificate(domain string) (*tlscommon.Certificate, error) {
	if domainutil.IsWildcardDomain(domain) {
		return nil, errors.New("wildcard domains are not supported with HTTP-01 challenge")
	}

	client, err := h.clientProvider.SetupClient(h.acmeEmail)
	if err != nil {
		return nil, err
	}

	ip, portInt, err := networktools.GetAvailablePort(h.ipAddress)
	if err != nil {
		slog.Error("failed to find free port", "domain", domain, "error", err)
		return nil, err
	}
	port := strconv.Itoa(portInt)

	http01Provider := http01.NewProviderServer(ip, port)
	if err := client.Challenge.SetHTTP01Provider(http01Provider); err != nil {
		return nil, err
	}

	request := certificate.ObtainRequest{
		Domains: []string{domain},
		Bundle:  true,
	}
	if err := h.publicher.Publish(ip, port); err != nil {
		return nil, err
	}
	cert, err := client.Certificate.Obtain(request)
	if err != nil {
		return nil, err
	}

	return &tlscommon.Certificate{
		CertPEM: cert.Certificate,
		KeyPEM:  cert.PrivateKey,
	}, nil
}
