package certmanager

import (
	"pb_launcher/internal/certmanager/domain"
	"pb_launcher/internal/certmanager/domain/repositories"
	"pb_launcher/internal/certmanager/repos"

	"go.uber.org/fx"
)

var Module = fx.Module("certmanager",
	fx.Provide(
		fx.Annotate(
			repos.NewCertRequestRepository,
			fx.As(new(repositories.CertRequestRepository)),
		),
	),

	fx.Provide(domain.NewCertRequestPlannerUsecase),
	fx.Provide(domain.NewCertRequestExecutorUsecase),
)
