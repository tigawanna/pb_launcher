package migrations

import (
	"pb_launcher/collections"
	"pb_launcher/utils"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		certRequests := core.NewBaseCollection(collections.CertRequests)
		certRequests.Fields.Add(
			&core.TextField{
				Name:     "domain",
				System:   true,
				Required: true,
				Pattern:  DomainPattern,
			},
			&core.SelectField{
				Name:      "status",
				Values:    []string{"pending", "approved", "failed"},
				MaxSelect: 1,
				System:    true,
				Required:  true,
			},
			&core.NumberField{
				Name:     "attempt",
				Min:      utils.Ptr[float64](1),
				System:   true,
				OnlyInt:  true,
				Required: true,
			},
			&core.TextField{
				Name:   "message",
				System: true,
			},
			&core.DateField{
				Name:   "requested",
				System: true,
			},
			&core.AutodateField{
				Name:     "created",
				System:   true,
				OnCreate: true,
			},
		)

		return app.Save(certRequests)
	}, func(app core.App) error {
		certRequests, err := app.FindCollectionByNameOrId(collections.CertRequests)
		if err != nil {
			return err
		}
		return app.Delete(certRequests)
	})
}
