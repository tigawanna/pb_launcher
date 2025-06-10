package proxy

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"pb_launcher/configs"
	"strings"

	"go.uber.org/fx"
)

func RunHttpProxy(lc fx.Lifecycle, handler *DynamicReverseProxy, cfg configs.Config) {
	mux := http.NewServeMux()
	if cfg.UseHttps() {
		mux.Handle("/", RedirectHTTPS(handler, cfg))
	} else {
		mux.Handle("/", handler)
	}

	addr := fmt.Sprintf("%s:%s", cfg.GetBindAddress(), cfg.GetBindPort())
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					slog.Error("proxy server error", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})

}

func RedirectHTTPS(next http.Handler, cfg configs.Config) http.Handler {
	httpsPort := cfg.GetBindHttpsPort()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.TLS != nil || strings.HasPrefix(r.Header.Get("X-Forwarded-Proto"), "https") {
			next.ServeHTTP(w, r)
			return
		}
		host := r.Host
		if httpsPort != "" && httpsPort != "443" {
			if h, _, err := net.SplitHostPort(r.Host); err == nil {
				host = h
			}
			host = net.JoinHostPort(host, httpsPort)
		}

		http.Redirect(w, r, "https://"+host+r.URL.RequestURI(), http.StatusPermanentRedirect)
	})
}
