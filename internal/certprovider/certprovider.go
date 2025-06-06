package certprovider

import (
	"fmt"
	"pb_launcher/configs"
	"pb_launcher/internal/certprovider/mkcert"
	"pb_launcher/internal/certprovider/selfsigned"
	"pb_launcher/internal/certprovider/tlscommon"
)

func NewProvider(c configs.Config) (tlscommon.Provider, error) {
	switch provider := c.GetTlsConfig().GetProvider(); provider {
	case "", "selfsigned":
		return selfsigned.NewSelfSignedProvider(), nil
	case "mkcert":
		return mkcert.NewMkcertProvider(), nil
	default:
		return nil, fmt.Errorf("%w: %s", tlscommon.ErrUnsupportedProvider, provider)
	}
}
