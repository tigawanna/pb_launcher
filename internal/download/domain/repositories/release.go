package repositories

import (
	"context"
	"pb_launcher/internal/download/domain/dtos"
)

type ReleaseVersionRepository interface {
	FindVersions(ctx context.Context) ([]dtos.Release, error)
	InsertVersions(ctx context.Context, releases []dtos.Release) error
}
