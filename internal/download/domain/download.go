package domain

import (
	"context"
	"log/slog"
	"pb_luncher/internal/download/domain/dtos"
	"pb_luncher/internal/download/domain/repositories"
	"pb_luncher/internal/download/domain/services"
)

type DownloadUsecase struct {
	service    services.ReleaseVersionsService
	repository repositories.ReleaseVersionRepository
}

func NewDownloadUsecase(
	service services.ReleaseVersionsService,
	repository repositories.ReleaseVersionRepository,
) *DownloadUsecase {
	return &DownloadUsecase{
		service:    service,
		repository: repository,
	}
}

// diffReleases returns the releases present in 'a' but not in 'b'.
// Example:
// a = [{Version: 1.0.0}, {Version: 1.2.0}, {Version: 2.0.0}]
// b = [{Version: 1.0.0}, {Version: 2.0.0}]
// Result: [{Version: 1.2.0}]
func (uc *DownloadUsecase) DiffReleases(a []dtos.Release, b []dtos.Release) []dtos.Release {
	bVersions := make(map[string]struct{})
	for _, release := range b {
		if release.Version != nil {
			bVersions[release.Version.String()] = struct{}{}
		}
	}

	var diff []dtos.Release
	for _, release := range a {
		if release.Version != nil {
			if _, exists := bVersions[release.Version.String()]; !exists {
				diff = append(diff, release)
			}
		}
	}

	return diff
}

func (uc *DownloadUsecase) Run(ctx context.Context) error {
	availableReleases, err := uc.service.FetchReleases(ctx)
	if err != nil {
		slog.Error("failed to fetch available releases", "error", err)
		return err
	}

	existingReleases, err := uc.repository.FindVersions(ctx)
	if err != nil {
		slog.Error("failed to find existing releases", "error", err)
		return err
	}

	diff := uc.DiffReleases(availableReleases, existingReleases)

	if err := uc.repository.InsertVersions(ctx, diff); err != nil {
		slog.Error("failed to insert new releases", "error", err)
		return err
	}

	// todo; sync releases with download dir

	return nil
}
