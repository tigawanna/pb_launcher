package services

import (
	"context"
	"pb_luncher/internal/download/domain/dtos"
)

type ReleaseVersionsService interface {
	FetchReleases(ctx context.Context) ([]dtos.Release, error)
}
