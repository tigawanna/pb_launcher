package download

import (
	"pb_launcher/internal/download/domain"
	"pb_launcher/internal/download/domain/repositories"
	"pb_launcher/internal/download/domain/services"
	"pb_launcher/internal/download/repos"
	infra_services "pb_launcher/internal/download/services"

	"go.uber.org/fx"
)

var Module = fx.Module("download",
	fx.Provide(
		fx.Annotate(
			repos.NewReleaseRepository,
			fx.As(new(repositories.ReleaseRepository)),
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
			infra_services.NewRepositoryArtifactStorage,
			fx.As(new(services.RepositoryArtifactStorage)),
		),
	),
	fx.Provide(domain.NewDownloadUsecase),
)
