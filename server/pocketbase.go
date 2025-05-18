package server

import (
	"context"
	"log/slog"
	"pb_launcher/configs"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"go.uber.org/fx"
)

func NewPocketbaseServer() *pocketbase.PocketBase {
	return pocketbase.New()
}
func StartPocketbase(lc fx.Lifecycle, pb *pocketbase.PocketBase, config *configs.Configs) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				slog.Info("Starting API server", "address", config.HttpAddr)
				if err := apis.Serve(pb.App, *config.ServeConfig); err != nil {
					slog.Error("API server encountered an error", "error", err)
					return
				}
				slog.Info("API server stopped gracefully")
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			slog.Info("Initiating shutdown sequence for PocketBase")
			event := new(core.TerminateEvent)
			event.App = pb

			if err := pb.OnTerminate().Trigger(event, func(e *core.TerminateEvent) error {
				slog.Info("Resetting bootstrap state for PocketBase")
				if err := e.App.ResetBootstrapState(); err != nil {
					slog.Error("Failed to reset bootstrap state", "error", err)
					return err
				}
				slog.Info("Bootstrap state successfully reset")
				return nil
			}); err != nil {
				slog.Error("Error during termination event handling", "error", err)
				return err
			}

			slog.Info("PocketBase shutdown completed successfully")
			return nil
		},
	})
}
