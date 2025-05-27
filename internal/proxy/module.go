package proxy

import (
	"context"
	"log/slog"
	"net/http"
	"pb_launcher/internal/proxy/domain"
	"pb_launcher/internal/proxy/domain/repositories"
	"pb_launcher/internal/proxy/repos"

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
	fx.Provide(domain.NewServiceDiscovery),
	fx.Provide(NewDynamicReverseProxy),
	fx.Invoke(RunProxy),
)

func RunProxy(lc fx.Lifecycle, handler *DynamicReverseProxy) {
	mux := http.NewServeMux()
	mux.Handle("/", handler)

	server := &http.Server{
		Addr:    "127.0.0.10:7080",
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
