package launcher

import (
	"pb_launcher/internal/launcher/domain"
	"pb_launcher/internal/launcher/domain/repositories"
	"pb_launcher/internal/launcher/domain/services"
	"pb_launcher/internal/launcher/repos"
	launcher_services "pb_launcher/internal/launcher/services"

	"go.uber.org/fx"
)

var Module = fx.Module("launcher",
	fx.Provide(
		fx.Annotate(
			repos.NewServiceRepository,
			fx.As(new(repositories.ServiceRepository)),
		),
		fx.Annotate(
			repos.NewCommandsRepository,
			fx.As(new(repositories.CommandsRepository)),
		),
	),
	fx.Provide(
		fx.Annotate(
			launcher_services.NewBinaryFinder,
			fx.As(new(services.BinaryFinder)),
		),
	),
	fx.Provide(domain.NewCleanServiceInstallTokenUsecase),
	fx.Provide(domain.NewLauncherManager),
)
