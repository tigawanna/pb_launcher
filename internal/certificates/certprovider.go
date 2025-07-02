package certificates

import (
	"fmt"
	"pb_launcher/configs"
	"pb_launcher/internal/certificates/certstore"
	"pb_launcher/internal/certificates/providers/cloudflare"
	"pb_launcher/internal/certificates/providers/mkcert"
	"pb_launcher/internal/certificates/providers/selfsigned"
	"pb_launcher/internal/certificates/tlscommon"

	"go.uber.org/fx"
)

func NewProvider(c configs.Config) (tlscommon.Provider, error) {
	switch provider := c.GetTlsConfig().GetProvider(); provider {
	case "selfsigned":
		return selfsigned.NewSelfSignedProvider(), nil
	case "mkcert":
		return mkcert.NewMkcertProvider(), nil
	case "cloudflare":
		return cloudflare.NewCloudflareProvider(c)
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
)
