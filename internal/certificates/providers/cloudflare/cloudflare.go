package cloudflare

import (
	"errors"
	"net/http"
	"pb_launcher/configs"
	"pb_launcher/internal/certificates/tlscommon"
	"time"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns/cloudflare"
	"github.com/go-acme/lego/v4/registration"
	"github.com/hashicorp/go-retryablehttp"
)

type CloudflareProvider struct {
	authToken string
	email     string
}

var _ tlscommon.Provider = (*CloudflareProvider)(nil)

func NewCloudflareProvider(c configs.Config) (*CloudflareProvider, error) {
	tlsConf := c.GetTlsConfig()

	authToken, ok := tlsConf.GetProp("auth_token")
	if !ok || authToken == "" {
		return nil, errors.New("missing or empty 'auth_token' in TLS config")
	}

	email, ok := tlsConf.GetProp("email")
	if !ok || email == "" {
		return nil, errors.New("missing or empty 'email' in TLS config")
	}

	return &CloudflareProvider{
		authToken: authToken,
		email:     email,
	}, nil
}

func (s *CloudflareProvider) RequestCertificate(domain string) (*tlscommon.Certificate, error) {
	privateKey, err := certcrypto.GeneratePrivateKey(certcrypto.RSA2048)
	if err != nil {
		return nil, err
	}

	user := &tlscommon.Account{
		Email: s.email,
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

	resource, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return nil, err
	}
	user.Registration = resource

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

	request := certificate.ObtainRequest{
		Domains: []string{domain, "*." + domain},
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
