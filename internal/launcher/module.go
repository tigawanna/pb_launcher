package launcher

import (
	"pb_launcher/internal/launcher/domain/repositories"
	"pb_launcher/internal/launcher/repos"

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

// func releaseSyncWorker(lc fx.Lifecycle, releaseUsecase *domain.DownloadUsecase, cfg *configs.Configs) {
// 	var runner = taskrunner.NewTaskRunner(
// 		func(ctx context.Context) {
// 			if err := releaseUsecase.Run(ctx); err != nil {
// 				slog.Error("error processing GitHub release sync task", "error", err)
// 			}
// 		},
// 		cfg.ReleaseSyncInterval,
// 	)

// 	lc.Append(fx.Hook{
// 		OnStart: func(ctx context.Context) error {
// 			slog.Info("starting GitHub release sync task runner", "interval", cfg.ReleaseSyncInterval)
// 			runner.Start()
// 			return nil
// 		},
// 		OnStop: func(ctx context.Context) error {
// 			slog.Info("stopping GitHub release sync task runner")
// 			runner.Stop()
// 			return nil
// 		},
// 	})
// }
