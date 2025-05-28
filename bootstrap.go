package main

import (
	"context"
	"log/slog"
	"pb_launcher/configs"
	"pb_launcher/helpers/taskrunner"
	download "pb_launcher/internal/download/domain"
	launcher "pb_launcher/internal/launcher/domain"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"go.uber.org/fx"
)

func Bootstrap(lc fx.Lifecycle,
	pb *pocketbase.PocketBase,
	downloader *download.DownloadUsecase,
	luncherManager *launcher.LauncherManager,
	config configs.Config,
	pbConfig *apis.ServeConfig,
) {
	var mu sync.Mutex
	// region Server
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			mu.Lock()
			defer mu.Unlock()
			doneChan := make(chan struct{})
			pb.OnServe().BindFunc(func(e *core.ServeEvent) error {
				close(doneChan)
				return e.Next()
			})
			go func() {
				slog.Info("Starting API server", "address", pbConfig.HttpAddr)
				if err := apis.Serve(pb.App, *pbConfig); err != nil {
					slog.Error("API server encountered an error", "error", err)
					return
				}
				slog.Info("API server stopped gracefully")
			}()
			<-doneChan
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

	// region Download
	var downloaderRunner = taskrunner.NewTaskRunner(
		func(ctx context.Context) {
			mu.Lock()
			defer mu.Unlock()
			if err := downloader.Run(ctx); err != nil {
				slog.Error("error processing GitHub release sync task", "error", err)
			}
		},
		config.GetReleaseSyncInterval(),
	)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			slog.Info("starting GitHub release sync task runner", "interval", config.GetReleaseSyncInterval())
			downloaderRunner.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			slog.Info("Stopping GitHub release sync task runner")
			downloaderRunner.Stop()
			return nil
		},
	})
	// region Luncher
	var recoveryDone atomic.Bool
	var luncherRunner = taskrunner.NewTaskRunner(
		func(ctx context.Context) {
			mu.Lock()
			defer mu.Unlock()
			if !recoveryDone.Load() {
				if err := luncherManager.RecoveryLastState(ctx); err != nil {
					slog.Error("recovery process failed", "error", err, "task", "luncherRunner")
					return
				}
				recoveryDone.Store(true)
			}
			if err := luncherManager.Run(ctx); err != nil {
				slog.Error("service runner task failed",
					"error", err,
					"task", "luncherRunner",
				)
			}
		},
		10*time.Second,
	)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			slog.Info("service runner starting",
				"interval", 10*time.Second,
				"task", "luncherRunner",
			)
			luncherRunner.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			slog.Info("service runner stopping",
				"task", "luncherRunner",
			)

			luncherRunner.Stop()
			if err := luncherManager.Dispose(); err != nil {
				slog.Error("failed to stop service runner",
					"error", err,
					"task", "luncherRunner",
				)
				return err
			}
			return nil
		},
	})
}
