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

type ReleaseVersionsGithub struct {
	releaseFilePattern *regexp.Regexp
	repositoryURL      string
}

var _ services.ReleaseVersionsService = (*ReleaseVersionsGithub)(nil)

func NewReleaseVersionsGithub(c *configs.Configs) *ReleaseVersionsGithub {
	return &ReleaseVersionsGithub{
		releaseFilePattern: c.ReleaseFilePattern,
		repositoryURL: fmt.Sprintf(
			"https://api.github.com/repos/%s/releases?per_page=10",
			c.GithubRepository,
		),
	}
}

func (rv *ReleaseVersionsGithub) searchAsset(assets []gjson.Result) (dtos.ReleaseAsset, bool) {
	for _, asset := range assets {
		downloadUrl := asset.Get("browser_download_url").String()
		if !rv.releaseFilePattern.MatchString(downloadUrl) {
			continue
		}
		return dtos.ReleaseAsset{
			DownloadURL:   downloadUrl,
			AssetFileName: asset.Get("name").String(),
			AssetSize:     asset.Get("size").Int(),
		}, true
	}
	return dtos.ReleaseAsset{}, false
}

func (rv *ReleaseVersionsGithub) parseReleases(data []byte) []dtos.Release {
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

		asset, ok := rv.searchAsset(obj.Get("assets").Array())
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

func (rv *ReleaseVersionsGithub) FetchReleases(ctx context.Context) ([]dtos.Release, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rv.repositoryURL, nil)
	if err != nil {
		slog.Error("failed to create GitHub releases request", "error", err, "url", rv.repositoryURL)
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("failed to fetch GitHub releases", "error", err, "url", rv.repositoryURL)
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		slog.Error("unexpected GitHub response status", "status_code", res.StatusCode, "url", rv.repositoryURL)
		return nil, fmt.Errorf("unexpected GitHub response status: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("error reading GitHub response body", "error", err, "url", rv.repositoryURL)
		return nil, err
	}
	releases := rv.parseReleases(data)
	if len(releases) > 3 {
		releases = releases[:3]
	}
	return releases, nil
}

func (rv *ReleaseVersionsGithub) Download(ctx context.Context, weburl string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, weburl, nil)
	if err != nil {
		slog.Error("failed to create GitHub releases request", "error", err, "url", rv.repositoryURL)
		return "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("failed to fetch GitHub releases", "error", err, "url", rv.repositoryURL)
		return "", err
	}
	defer res.Body.Close()

	tempFile, err := os.CreateTemp("", "release-*.zip")
	if err != nil {
		slog.Error("failed to create temp file for release", "error", err)
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
