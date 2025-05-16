package repos

import (
	"context"
	"log/slog"
	"pb_luncher/collections"
	"pb_luncher/internal/download/domain/dtos"
	"pb_luncher/internal/download/domain/repositories"

	"github.com/hashicorp/go-version"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

type ReleaseVersionRepository struct {
	app *pocketbase.PocketBase
}

var _ repositories.ReleaseVersionRepository = (*ReleaseVersionRepository)(nil)

func NewReleaseVersionRepository(app *pocketbase.PocketBase) *ReleaseVersionRepository {
	return &ReleaseVersionRepository{app: app}
}

func (r *ReleaseVersionRepository) FindVersions(ctx context.Context) ([]dtos.Release, error) {
	records, err := r.app.FindAllRecords(collections.Releases)
	if err != nil {
		slog.Error("failed to fetch releases from database", "error", err)
		return nil, err
	}

	releases := make([]dtos.Release, 0, len(records))

	for _, record := range records {
		versionString := record.GetString("version")
		v, err := version.NewVersion(versionString)
		if err != nil {
			slog.Warn("invalid version format", "version", versionString, "error", err)
			continue
		}

		releases = append(releases, dtos.Release{
			Version:     v,
			ReleaseName: record.GetString("release_name"),
			PublishedAt: record.GetDateTime("published_at").Time(),
			ReleaseAsset: dtos.ReleaseAsset{
				AssetFileName: record.GetString("asset_file_name"),
				DownloadURL:   record.GetString("download_url"),
				AssetSize:     int64(record.GetInt("asset_size")),
			},
		})
	}
	return releases, nil
}

func (r *ReleaseVersionRepository) InsertVersions(ctx context.Context, releases []dtos.Release) error {
	if len(releases) == 0 {
		slog.Info("no new releases to insert")
		return nil
	}

	collection, err := r.app.FindCollectionByNameOrId(collections.Releases)
	if err != nil {
		slog.Error("failed to find releases collection", "error", err)
		return err
	}

	for _, release := range releases {
		record := core.NewRecord(collection)
		record.Set("version", release.Version.String())
		record.Set("release_name", release.ReleaseName)
		record.Set("published_at", release.PublishedAt)
		record.Set("asset_file_name", release.AssetFileName)
		record.Set("download_url", release.DownloadURL)
		record.Set("asset_size", release.AssetSize)

		if err := r.app.Save(record); err != nil {
			slog.Error("failed to save release record", "version", release.Version.String(), "error", err)
			return err
		}
	}
	return nil
}
