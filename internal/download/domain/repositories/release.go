package repositories

import (
	"context"
	"pb_launcher/internal/download/domain/dtos"
)

type ReleaseRepository interface {
	ListRepositories(ctx context.Context) ([]dtos.Repository, error)
	ListReleases(ctx context.Context, repositoryID string) ([]dtos.Release, error)
	SaveReleases(ctx context.Context, releases []dtos.Release) error
}
