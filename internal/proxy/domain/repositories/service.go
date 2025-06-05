package repositories

import (
	"context"
	"errors"
	"pb_launcher/internal/proxy/domain/dtos"
)

var ErrNotFound = errors.New("not found")

type ServiceRepository interface {
	FindServiceIDByDomain(ctx context.Context, doamin string) (*string, error)
	FindRunningServiceByID(ctx context.Context, id string) (*dtos.RunningServiceDto, error)
}
