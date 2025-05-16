package download

import (
	"context"
	"log/slog"
	"pb_luncher/configs"
	"pb_luncher/helpers/taskrunner"
	"pb_luncher/internal/download/domain"
	"pb_luncher/internal/download/domain/repositories"
	"pb_luncher/internal/download/domain/services"
	"pb_luncher/internal/download/repos"
	infra_services "pb_luncher/internal/download/services"

	"go.uber.org/fx"
)

var Module = fx.Module("download",
	fx.Provide(
		fx.Annotate(
			repos.NewReleaseVersionRepository,
			fx.As(new(repositories.ReleaseVersionRepository)),
		),
	),
	fx.Provide(
		fx.Annotate(
			infra_services.NewReleaseVersionsGithub,
			fx.As(new(services.ReleaseVersionsService)),
		),
	),
	fx.Provide(
		fx.Annotate(
			infra_services.NewArtifactStorage,
			fx.As(new(services.ArtifactStorage)),
		),
	),
	fx.Provide(domain.NewDownloadUsecase),
	fx.Invoke(releaseSyncWorker),
)

func releaseSyncWorker(lc fx.Lifecycle, releaseUsecase *domain.DownloadUsecase, cfg *configs.Configs) {
	var runner = taskrunner.NewTaskRunner(
		func(ctx context.Context) {
			if err := releaseUsecase.Run(ctx); err != nil {
				slog.Error("error processing GitHub release sync task", "error", err)
			}
		},
		cfg.ReleaseSyncInterval,
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
			return nil
		},
	})
}
