package cloudflare

import (
	"errors"
	"net/http"
	"pb_launcher/configs"
	"pb_launcher/internal/certificates/tlscommon"
	"pb_launcher/utils/domainutil"
	"time"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/providers/dns/cloudflare"
)

type CloudflareProvider struct {
	clientProvider *tlscommon.LetsEncryptClientAccountProvider
	authToken      string
	acmeEmail      string
}

var _ tlscommon.Provider = (*CloudflareProvider)(nil)

func NewCloudflareProvider(c configs.Config,
	clientProvider *tlscommon.LetsEncryptClientAccountProvider) (*CloudflareProvider, error) {

	tlsConf := c.GetTlsConfig()

	authToken, ok := tlsConf.GetProp("auth_token")
	if !ok || authToken == "" {
		return nil, errors.New("missing or empty cloudflare 'auth_token' in TLS config")
	}

	return &CloudflareProvider{
		clientProvider: clientProvider,
		authToken:      authToken,
		acmeEmail:      c.GetAcmeEmail(),
	}, nil
}

func (s *CloudflareProvider) RequestCertificate(domain string) (*tlscommon.Certificate, error) {
	client, err := s.clientProvider.SetupClient(s.acmeEmail)
	if err != nil {
		return nil, err
	}

	cloudflareConfig := &cloudflare.Config{
		AuthToken:          s.authToken,
		PropagationTimeout: 2 * time.Minute,
		TTL:                120,
		PollingInterval:    dns01.DefaultPollingInterval,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	provider, err := cloudflare.NewDNSProviderConfig(cloudflareConfig)
	if err != nil {
		return nil, err
	}

	err = client.Challenge.SetDNS01Provider(provider)
	if err != nil {
		return nil, err
	}

	baseDomain := domainutil.BaseDomain(domain)
	wildcardDomain := domainutil.ToWildcardDomain(domain)

	request := certificate.ObtainRequest{
		Domains: []string{baseDomain, wildcardDomain},
		Bundle:  true,
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
