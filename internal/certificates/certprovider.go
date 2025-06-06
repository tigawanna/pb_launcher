package certificates

import (
	"errors"
	"fmt"
	"pb_launcher/configs"
	"pb_launcher/internal/certificates/certstore"
	"pb_launcher/internal/certificates/providers/mkcert"
	"pb_launcher/internal/certificates/providers/selfsigned"
	"pb_launcher/internal/certificates/tlscommon"
	"strings"

	"go.uber.org/fx"
)

func NewProvider(c configs.Config) (tlscommon.Provider, error) {
	switch provider := c.GetTlsConfig().GetProvider(); provider {
	case "selfsigned":
		return selfsigned.NewSelfSignedProvider(), nil
	case "mkcert":
		return mkcert.NewMkcertProvider(), nil
	default:
		return nil, fmt.Errorf("%w: %s", tlscommon.ErrUnsupportedProvider, provider)
	}
}

var Module = fx.Module("tls_provider",
	fx.Provide(
		fx.Private,
		certstore.NewTlsStorer,
	),
	fx.Provide(fx.Annotate(
		certstore.NewTlsStorerCache,
		fx.As(new(tlscommon.Store)),
	)),
	fx.Provide(NewProvider),
	fx.Invoke(PrepareCertificates),
)

func PrepareCertificates(lc fx.Lifecycle, cfg configs.Config, provider tlscommon.Provider, storer tlscommon.Store) error {
	if !cfg.UseHttps() {
		return nil
	}

	domain := cfg.GetDomain()
	if !strings.HasPrefix(domain, "*.") {
		domain = "*." + domain
	}

	_, err := storer.Resolve(domain)
	if err == nil {
		return nil
	}

	if !errors.Is(err, tlscommon.ErrCertificateNotFound) &&
		!errors.Is(err, tlscommon.ErrInvalidPEM) &&
		!errors.Is(err, tlscommon.ErrCertificateExpired) {
		return err
	}

	requestAndStoreCertificate := func() error {
		cert, err := provider.RequestCertificate(domain)
		if err != nil {
			return err
		}
		return storer.Store(domain, cert)
	}

	// TODO

	return requestAndStoreCertificate()
}
