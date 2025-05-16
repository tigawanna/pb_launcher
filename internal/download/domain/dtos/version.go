package dtos

import (
	"time"

	"github.com/hashicorp/go-version"
)

type ReleaseAsset struct {
	DownloadURL   string
	AssetFileName string
	AssetSize     int64
}

type Release struct {
	Version     *version.Version
	ReleaseName string
	PublishedAt time.Time
	ReleaseAsset
}
