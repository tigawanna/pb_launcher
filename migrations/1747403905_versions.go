package migrations

import (
	"pb_launcher/collections"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		repo := core.NewBaseCollection(collections.Repositories)
		retentionMin := 1.0
		retentionMax := 6.0
		repo.Fields.Add(
			&core.TextField{
				Name:        "name",
				Required:    true,
				System:      true,
				Presentable: true,
			},
			&core.TextField{
				Name:     "repository",
				System:   true,
				Required: true,
				Pattern:  `^.+\/.+$`,
			},
			&core.TextField{
				Name:   "token",
				System: true,
			},
			&core.NumberField{
				Name:    "retention",
				OnlyInt: true,
				Min:     &retentionMin,
				Max:     &retentionMax,
			},
			&core.TextField{
				Name:     "release_file_pattern",
				System:   true,
				Required: true,
			},
			&core.TextField{
				Name:     "exec_file_pattern",
				System:   true,
				Required: true,
			},
			&core.BoolField{
				Name:   "disabled",
				System: true,
			},
		)
		if err := app.Save(repo); err != nil {
			return err
		}

		// default repository
		pb := core.NewRecord(repo)
		pb.Set("id", "pb91u2l315h29a5")
		pb.Set("name", "PocketBase")
		pb.Set("retention", 3)
		pb.Set("repository", "pocketbase/pocketbase")
		pb.Set("release_file_pattern", `pocketbase_.+_linux_amd64\.zip`)
		pb.Set("exec_file_pattern", `^pocketbase`)

		if err := app.Save(pb); err != nil {
			return err
		}

		releases := core.NewBaseCollection(collections.Releases)
		releases.System = true
		releases.Fields.Add(
			&core.RelationField{
				Name:         "repository",
				System:       true,
				CollectionId: repo.Id,
				Presentable:  true,
				MinSelect:    1,
				MaxSelect:    1,
			},
			&core.TextField{
				Name:        "version",
				Presentable: true,
				Required:    true,
				Pattern:     "^[0-9]+(\\.[0-9]+)*$",
				System:      true,
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
			&core.TextField{
				Name:     "asset_id",
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
			"CREATE UNIQUE INDEX idx_releases ON releases (repository,version)",
		)
		return app.Save(releases)
	}, nil)
}
