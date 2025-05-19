package dtos

import (
	"regexp"
	"time"

	"github.com/hashicorp/go-version"
)

type Repository struct {
	ID                 string
	Repo               string
	Token              string
	ReleaseFilePattern *regexp.Regexp
	ExecFilePattern    *regexp.Regexp
	Retention          int
}

type ReleaseAsset struct {
	AssetID       string
	DownloadURL   string
	AssetFileName string
	AssetSize     int64
}

type Release struct {
	RepositoryID string
	Version      *version.Version
	ReleaseName  string
	PublishedAt  time.Time
	ReleaseAsset
}
