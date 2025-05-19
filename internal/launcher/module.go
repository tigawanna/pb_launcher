package launcher

import (
	"context"
	"log/slog"
	"pb_launcher/configs"
	"pb_launcher/helpers/taskrunner"
	"pb_launcher/internal/launcher/domain"
	"pb_launcher/internal/launcher/domain/repositories"
	"pb_launcher/internal/launcher/repos"
	"time"

	"go.uber.org/fx"
)

var Module = fx.Module("launcher",
	fx.Provide(
		fx.Annotate(
			repos.NewServiceRepository,
			fx.As(new(repositories.ServiceRepository)),
		),
	),

	// fx.Provide(domain.NewDownloadUsecase),
	// fx.Invoke(releaseSyncWorker),
)

func releaseSyncWorker(lc fx.Lifecycle, luncher *domain.LauncherManager, cfg *configs.Configs) {
	var runner = taskrunner.NewTaskRunner(
		func(ctx context.Context) {
			if err := luncher.Run(ctx); err != nil {
				slog.Error("error processing GitHub release sync task", "error", err)
			}
		},
		10*time.Second,
	)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			slog.Info("starting GitHub release sync task runner", "interval", cfg.ReleaseSyncInterval)
			runner.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			slog.Info("stopping GitHub release sync task runner")
			runner.Stop()
			if err := luncher.Stop(); err != nil {
				//
			}
			return nil
		},
	})
}
