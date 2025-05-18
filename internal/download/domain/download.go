package domain

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"pb_launcher/helpers/unzip"
	"pb_launcher/internal/download/domain/dtos"
	"pb_launcher/internal/download/domain/repositories"
	"pb_launcher/internal/download/domain/services"
)

type DownloadUsecase struct {
	service         services.ReleaseVersionsService
	repository      repositories.ReleaseVersionRepository
	artifactStorage services.ArtifactStorage
	unzip           *unzip.Unzip
}

func NewDownloadUsecase(
	service services.ReleaseVersionsService,
	repository repositories.ReleaseVersionRepository,
	artifactStorage services.ArtifactStorage,
	unzip *unzip.Unzip,
) *DownloadUsecase {
	return &DownloadUsecase{
		service:         service,
		repository:      repository,
		artifactStorage: artifactStorage,
		unzip:           unzip,
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

func (uc *DownloadUsecase) processDownload(ctx context.Context, release dtos.Release) error {
	zipPath, err := uc.service.Download(ctx, release.DownloadURL)
	if err != nil {
		slog.Error("failed to download release", "error", err, "url", release.DownloadURL)
		return err
	}
	defer os.Remove(zipPath)

	tempDir, err := os.MkdirTemp("", "release-extract-*")
	if err != nil {
		slog.Error("failed to create temp directory for extraction", "error", err)
		return err
	}
	defer os.RemoveAll(tempDir)

	extractedPaths, err := uc.unzip.Extract(zipPath, tempDir)
	if err != nil {
		slog.Error("failed to extract release", "error", err, "zipPath", zipPath)
		return err
	}

	for _, relativePath := range extractedPaths {
		fullPath := filepath.Join(tempDir, relativePath)
		file, err := os.Open(fullPath)
		if err != nil {
			slog.Error("failed to open extracted file", "error", err, "path", fullPath)
			return err
		}
		defer file.Close()
		outFilePath := filepath.Join(release.Version.String(), relativePath)
		if _, err := uc.artifactStorage.Save(ctx, outFilePath, file); err != nil {
			file.Close()
			slog.Error("failed to save extracted file", "error", err, "path", relativePath)
			return err
		}

		file.Close()
	}
	slog.Info("release downloaded successfully", "version", release.Version.String())
	return nil
}

func (uc *DownloadUsecase) resolveMissingReleases(ctx context.Context) error {
	releases, err := uc.repository.FindVersions(ctx)
	if err != nil {
		slog.Error("failed to retrieve available releases", "error", err)
		return err
	}

	downloadedVersions, err := uc.artifactStorage.Versions(ctx)
	if err != nil {
		slog.Error("failed to retrieve downloaded versions", "error", err)
		return err
	}

	downloadedSet := make(map[string]struct{}, len(downloadedVersions))
	for _, version := range downloadedVersions {
		downloadedSet[version.String()] = struct{}{}
	}

	var pendingDownloads []dtos.Release
	for _, release := range releases {
		if _, exists := downloadedSet[release.Version.String()]; !exists {
			pendingDownloads = append(pendingDownloads, release)
		}
	}

	for _, pending := range pendingDownloads {
		if err := uc.processDownload(ctx, pending); err != nil {
			slog.Error("failed to process download", "version", pending.Version.String(), "error", err)
			return err
		}
	}

	return nil
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

	if err := uc.resolveMissingReleases(ctx); err != nil {
		slog.Error("failed to resolve missing releases", "error", err)
		return err
	}
	return nil
}
