package proxy

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"
	"pb_launcher/configs"
	"pb_launcher/internal/certificates/tlscommon"
	"strings"

	"go.uber.org/fx"
)

func RunHTTPSProxy(
	lc fx.Lifecycle,
	proxyHandler *DynamicReverseProxy,
	certStore tlscommon.Store,
	cfg configs.Config,
) {
	mux := http.NewServeMux()
	mux.Handle("/", proxyHandler)

	baseDomain := cfg.GetDomain()

	addr := fmt.Sprintf("%s:%s", cfg.GetBindAddress(), cfg.GetBindHttpsPort())
	server := &http.Server{
		Addr: addr,
		TLSConfig: &tls.Config{
			GetCertificate: func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
				domain := hello.ServerName
				if strings.HasSuffix(domain, baseDomain) {
					domain = "*." + baseDomain
				}
				cert, err := certStore.Resolve(domain)
				if err != nil {
					slog.Error("failed to resolve certificate", "domain", domain, "error", err)
					return nil, err
				}
				tlsCert, err := tls.X509KeyPair(cert.CertPEM, cert.KeyPEM)
				return &tlsCert, err
			},
		},
		Handler: mux,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				err := server.ListenAndServeTLS("", "")
				if err != nil && err != http.ErrServerClosed {
					slog.Error("HTTPS proxy encountered an error", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			slog.Info("shutting down HTTPS proxy", "addr", server.Addr)
			return server.Shutdown(ctx)
		},
	})
}
