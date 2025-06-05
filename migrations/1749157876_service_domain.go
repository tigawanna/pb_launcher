package migrations

import (
	"pb_launcher/collections"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		services, err := app.FindCollectionByNameOrId(collections.Services)
		if err != nil {
			return err
		}
		servicesDomain := core.NewBaseCollection(collections.ServicesDomains)
		servicesDomain.Fields.Add(
			&core.TextField{
				Name:     "domain",
				System:   true,
				Required: true,
				Pattern:  `^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]`,
			},
			&core.RelationField{
				Name:         "service",
				CollectionId: services.Id,
				System:       true,
				Required:     true,
			},
		)
		servicesDomain.Indexes = append(
			servicesDomain.Indexes,
			"CREATE UNIQUE INDEX idx_services_domains_domain ON services_domains (domain)",
		)

		return app.Save(servicesDomain)
	}, func(app core.App) error {
		sd, err := app.FindCollectionByNameOrId(collections.ServicesDomains)
		if err != nil {
			return err
		}
		return app.Delete(sd)
	})
}
