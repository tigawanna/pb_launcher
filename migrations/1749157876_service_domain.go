package migrations

import (
	"pb_launcher/collections"
	"pb_launcher/utils"

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
			&core.SelectField{
				Name:      "use_https",
				Values:    []string{"no", "yes"},
				MaxSelect: 1,
				System:    true,
				Required:  true,
			},
			&core.TextField{
				Name:   "status",
				System: true,
			},
		)
		servicesDomain.Indexes = append(
			servicesDomain.Indexes,
			"CREATE UNIQUE INDEX idx_services_domains_domain ON services_domains (domain)",
		)

		servicesDomain.ListRule = utils.StrPointer(`@request.auth.id != ""`)
		servicesDomain.ViewRule = utils.StrPointer(`@request.auth.id != ""`)
		servicesDomain.CreateRule = utils.StrPointer(`@request.auth.id != ""`)
		servicesDomain.UpdateRule = utils.StrPointer(`@request.auth.id != ""`)
		servicesDomain.DeleteRule = utils.StrPointer(`@request.auth.id != ""`)

		return app.Save(servicesDomain)
	}, func(app core.App) error {
		sd, err := app.FindCollectionByNameOrId(collections.ServicesDomains)
		if err != nil {
			return err
		}
		return app.Delete(sd)
	})
}
