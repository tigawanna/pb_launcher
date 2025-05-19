package migrations

import (
	"pb_launcher/collections"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		releases, err := app.FindCollectionByNameOrId(collections.Releases)
		if err != nil {
			return err
		}
		services := core.NewBaseCollection(collections.Services)
		services.Fields.Add(
			&core.TextField{
				Name:        "name",
				Presentable: true,
				System:      true,
			},
			&core.RelationField{
				Name:         "release",
				CollectionId: releases.Id,
				System:       true,
				Required:     true,
			},
			&core.TextField{
				Name:    "ip",
				System:  true,
				Pattern: `^\d{1,3}(?:\.\d{1,3}){3}$`,
			},
			&core.NumberField{
				Name:   "port",
				System: true,
			},
			&core.SelectField{
				Name:   "status",
				System: true,
				Values: []string{"idle", "running", "stopped"},
			},
			&core.SelectField{
				Name:     "restart_policy",
				System:   true,
				Required: true,
				Values:   []string{"no", "on-failure"},
			},
			&core.TextField{
				Name:   "error_message",
				System: true,
			},
			&core.DateField{
				Name:   "last_started_at",
				System: true,
			},
			&core.DateField{
				Name:   "created_at",
				System: true,
			},
			&core.DateField{
				Name:   "deleted_at",
				System: true,
			},
		)
		return app.Save(services)
	}, func(app core.App) error {
		services, err := app.FindCollectionByNameOrId(collections.Services)
		if err != nil {
			return err
		}
		return app.Delete(services)
	})
}
