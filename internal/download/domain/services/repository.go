package services

import (
	"context"
	"pb_launcher/internal/download/domain/dtos"
)

type ReleaseVersionsService interface {
	FetchReleases(ctx context.Context) ([]dtos.Release, error)
	Download(ctx context.Context, weburl string) (string, error)
}
