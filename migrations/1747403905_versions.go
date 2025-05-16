package migrations

import (
	"pb_luncher/collections"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		releases := core.NewBaseCollection(collections.Releases)
		releases.System = true
		releases.Fields.Add(
			&core.TextField{
				Name:     "version",
				Required: true,
				Pattern:  "^[0-9]+(\\.[0-9]+)*$",
				System:   true,
			},
			&core.TextField{
				Name:     "release_name",
				Required: true,
				System:   true,
			},
			&core.DateField{
				Name:     "published_at",
				System:   true,
				Required: true,
			},
			&core.TextField{
				Name:     "asset_file_name",
				System:   true,
				Required: true,
			},
			&core.URLField{
				Name:        "download_url",
				System:      true,
				OnlyDomains: []string{"github.com"},
				Required:    true,
			},
			&core.NumberField{
				Name:   "asset_size",
				System: true,
			},
		)
		releases.Indexes = append(
			releases.Indexes,
			"CREATE UNIQUE INDEX idx_releases ON releases (version)",
		)
		return app.Save(releases)
	}, nil)
}
