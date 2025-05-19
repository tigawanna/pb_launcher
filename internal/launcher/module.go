package launcher

import (
	"pb_launcher/internal/launcher/domain"
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
	fx.Provide(domain.NewLauncherManager),
)
