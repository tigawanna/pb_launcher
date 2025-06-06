package proxy

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"pb_launcher/configs"
	"pb_launcher/internal/proxy/domain"
	"pb_launcher/internal/proxy/domain/repositories"
	"pb_launcher/internal/proxy/repos"

	"github.com/fatih/color"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"proxy",
	fx.Provide(
		fx.Annotate(
			repos.NewServiceRepository,
			fx.As(new(repositories.ServiceRepository)),
		),
	),
	fx.Provide(
		domain.NewServiceDiscovery,
		domain.NewDomainServiceDiscovery,
	),
	fx.Provide(NewDynamicReverseProxy),
	fx.Invoke(RunProxy, PrintProxyInfo),
)

func RunProxy(lc fx.Lifecycle, handler *DynamicReverseProxy, c configs.Config) {
	mux := http.NewServeMux()
	mux.Handle("/", handler)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", c.GetBindAddress(), c.GetBindPort()),
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

func PrintProxyInfo(c configs.Config) {
	regular := color.New()
	port := c.GetBindPort()
	scheme := map[string]string{"80": "http", "443": "https"}[port]
	if scheme == "" {
		scheme = "http"
		if c.UseHttps() {
			scheme = "https"
		}
		addr := fmt.Sprintf("%s://%s:%s", scheme, c.GetBindAddress(), port)
		pub := fmt.Sprintf("%s://%s:%s", scheme, c.GetDomain(), port)
		regular.Printf("├─ Proxy:  %s\n", color.CyanString(addr))
		regular.Printf("├─ Public: %s\n", color.CyanString(pub))
		return
	}
	regular.Printf("├─ Proxy:  %s\n", color.CyanString("%s://%s", scheme, c.GetBindAddress()))
	regular.Printf("├─ Public: %s\n", color.CyanString("%s://%s", scheme, c.GetDomain()))
}
