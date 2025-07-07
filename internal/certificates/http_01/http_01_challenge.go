package http01

import (
	"errors"
	"fmt"
	"log/slog"
	"pb_launcher/configs"
	"pb_launcher/internal/certificates/tlscommon"
	"pb_launcher/utils/domainutil"
	"pb_launcher/utils/networktools"
	"strconv"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"github.com/hashicorp/go-retryablehttp"
)

type HTTP01TLSCertificateRequestService struct {
	acmeEmail string
	ipAddress string
	publicher *Http01ChallengeAddressPublisher
}

var _ tlscommon.Provider = (*HTTP01TLSCertificateRequestService)(nil)

func NewHTTP01TLSCertificateRequestService(
	publicher *Http01ChallengeAddressPublisher,
	c configs.Config,
) *HTTP01TLSCertificateRequestService {
	return &HTTP01TLSCertificateRequestService{
		publicher: publicher,
		ipAddress: c.GetBindAddress(),
		acmeEmail: c.GetAcmeEmail(),
	}
}

// RequestCertificate implements tlscommon.Provider.
func (h *HTTP01TLSCertificateRequestService) RequestCertificate(domain string) (*tlscommon.Certificate, error) {
	if domainutil.IsWildcardDomain(domain) {
		return nil, errors.New("wildcard domains are not supported with HTTP-01 challenge")
	}
	privateKey, err := certcrypto.GeneratePrivateKey(certcrypto.RSA2048)
	if err != nil {
		return nil, err
	}
	email := h.acmeEmail
	if email == "" {
		slog.Warn("`acme_email` not configured, using default placeholder", "email", fmt.Sprintf("admin@%s", domain))
		email = fmt.Sprintf("admin@%s", domain)
	}
	user := &tlscommon.Account{
		Email: email,
		Key:   privateKey,
	}

	config := lego.NewConfig(user)

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5
	retryClient.HTTPClient = config.HTTPClient
	retryClient.Logger = nil
	config.HTTPClient = retryClient.StandardClient()

	client, err := lego.NewClient(config)
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
	resource, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return nil, err
	}
	user.Registration = resource

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
