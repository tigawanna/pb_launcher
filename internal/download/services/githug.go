package services

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"pb_launcher/configs"
	"pb_launcher/internal/download/domain/dtos"
	"pb_launcher/internal/download/domain/services"
	"regexp"

	"github.com/hashicorp/go-version"
	"github.com/tidwall/gjson"
)

type ReleaseVersionsGithub struct{}

var _ services.ReleaseVersionsService = (*ReleaseVersionsGithub)(nil)

func NewReleaseVersionsGithub(c *configs.Configs) *ReleaseVersionsGithub {
	return &ReleaseVersionsGithub{}
}

func (rv *ReleaseVersionsGithub) buildUrl(repo dtos.Repository) string {
	return fmt.Sprintf(
		"https://api.github.com/repos/%s/releases?per_page=%d",
		repo.Repo,
		repo.Retention,
	)
}

func (rv *ReleaseVersionsGithub) buildDwnUrl(repo dtos.Repository, asset dtos.ReleaseAsset) string {
	return fmt.Sprintf(
		"https://api.github.com/repos/%s/releases/assets/%s",
		repo.Repo,
		asset.AssetID,
	)
}

func (rv *ReleaseVersionsGithub) searchAsset(assets []gjson.Result, releaseFilePattern *regexp.Regexp) (dtos.ReleaseAsset, bool) {
	for _, asset := range assets {
		downloadUrl := asset.Get("browser_download_url").String()
		id := asset.Get("id").Int()
		if !releaseFilePattern.MatchString(downloadUrl) {
			continue
		}
		return dtos.ReleaseAsset{
			AssetID:       fmt.Sprint(id),
			DownloadURL:   downloadUrl,
			AssetFileName: asset.Get("name").String(),
			AssetSize:     asset.Get("size").Int(),
		}, true
	}
	return dtos.ReleaseAsset{}, false
}

func (rv *ReleaseVersionsGithub) parseReleases(data []byte, releaseFilePattern *regexp.Regexp) []dtos.Release {
	var releases []dtos.Release
	results := gjson.ParseBytes(data).Array()
	// prerelease
	for _, obj := range results {
		if obj.Get("prerelease").Bool() {
			continue
		}
		releaseName := obj.Get("name").String()
		publishedAt := obj.Get("published_at").Time()
		tagName := obj.Get("tag_name").String()

		releaseVersion, err := version.NewVersion(tagName)
		if err != nil {
			slog.Error("invalid version format", "error", err, "tag_name", tagName)
			continue
		}
		if releaseVersion == nil {
			continue
		}

		asset, ok := rv.searchAsset(obj.Get("assets").Array(), releaseFilePattern)
		if !ok {
			slog.Info("no matching asset found", "release_name", releaseName)
			continue
		}

		releases = append(releases, dtos.Release{
			Version:      releaseVersion,
			ReleaseName:  releaseName,
			PublishedAt:  publishedAt,
			ReleaseAsset: asset,
		})
	}

	return releases
}

func (rv *ReleaseVersionsGithub) FetchReleases(ctx context.Context, repo dtos.Repository) ([]dtos.Release, error) {
	repositoryURL := rv.buildUrl(repo)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, repositoryURL, nil)
	if err != nil {
		slog.Error("failed to create GitHub releases request", "error", err, "url", repositoryURL)
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	if repo.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", repo.Token))
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("failed to fetch GitHub releases", "error", err, "url", repositoryURL)
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		slog.Error("unexpected GitHub response status", "status_code", res.StatusCode, "url", repositoryURL)
		return nil, fmt.Errorf("unexpected GitHub response status: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("error reading GitHub response body", "error", err, "url", repositoryURL)
		return nil, err
	}
	releases := rv.parseReleases(data, repo.ReleaseFilePattern)
	if len(releases) > 3 {
		releases = releases[:3]
	}
	for i := range releases {
		releases[i].RepositoryID = repo.ID
	}
	return releases, nil
}

func (rv *ReleaseVersionsGithub) Download(ctx context.Context, repo dtos.Repository, asset dtos.ReleaseAsset) (string, error) {
	downloadUrl := rv.buildDwnUrl(repo, asset)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadUrl, nil)
	if err != nil {
		slog.Error("create HTTP request for GitHub release download", "error", err, "method", http.MethodGet, "url", downloadUrl)
		return "", err
	}
	if repo.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", repo.Token))
	}
	req.Header.Set("Accept", "application/octet-stream")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("execute HTTP request for GitHub release download", "error", err, "method", req.Method, "url", req.URL.String())
		return "", err
	}
	defer res.Body.Close()

	tempFile, err := os.CreateTemp("", "release-*.zip")
	if err != nil {
		slog.Error("Failed to create temporary file for GitHub release", "error", err)
		return "", err
	}

	if _, err := io.Copy(tempFile, res.Body); err != nil {
		slog.Error("failed to write release to temp file", "error", err, "path", tempFile.Name())
		tempFile.Close()
		os.Remove(tempFile.Name())
		return "", err
	}

	if err := tempFile.Close(); err != nil {
		slog.Error("failed to close temp file", "error", err, "path", tempFile.Name())
		os.Remove(tempFile.Name())
		return "", err
	}
	return tempFile.Name(), nil
}
