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
		comands := core.NewBaseCollection(collections.ServicesComands)
		comands.System = true
		comands.Fields.Add(
			&core.RelationField{
				Name:         "service",
				CollectionId: services.Id,
				System:       true,
				Required:     true,
				MinSelect:    1,
				MaxSelect:    1,
			},
			&core.SelectField{
				Name:      "action",
				Values:    []string{"stop", "start", "restart"},
				System:    true,
				Required:  true,
				MaxSelect: 1,
			},
			&core.SelectField{
				Name:   "status",
				System: true,
				Values: []string{"pending", "success", "error"},
			},
			&core.TextField{
				Name:   "error_message",
				System: true,
			},
			&core.DateField{
				Name:   "executed",
				System: true,
			},
			&core.AutodateField{
				Name:     "created",
				System:   true,
				OnCreate: true,
			},
		)
		comands.Indexes = append(comands.Indexes,
			`CREATE INDEX idx_comands_service ON comands(service)`,
			`CREATE INDEX idx_comands_executed ON comands(executed)`,
			`CREATE INDEX idx_comands_created ON comands(created)`,
			`CREATE INDEX idx_comands_status ON comands(status)`,
		)
		return app.Save(comands)
	}, nil)
}
