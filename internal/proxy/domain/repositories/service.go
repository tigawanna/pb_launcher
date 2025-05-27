package repositories

import (
	"context"
	"errors"
	"pb_launcher/internal/proxy/domain/dtos"
)

var ErrServiceNotFound = errors.New("service not found")

type ServiceRepository interface {
	FindRunningServiceByID(ctx context.Context, id string) (*dtos.RunningServiceDto, error)
}
