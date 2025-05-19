package services

import (
	"context"
	"pb_launcher/internal/download/domain/dtos"
)

type ReleaseVersionsService interface {
	FetchReleases(ctx context.Context, repo dtos.Repository) ([]dtos.Release, error)
	Download(ctx context.Context, repo dtos.Repository, asset dtos.ReleaseAsset) (string, error)
}
