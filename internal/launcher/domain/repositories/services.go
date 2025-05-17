package repositories

import (
	"context"
	"pb_luncher/internal/launcher/domain/models"
)

type ServiceRepository interface {
	FindAll(ctx context.Context) ([]models.Service, error)
	UpdateServiceStatus(ctx context.Context, id string, status models.ServiceStatus) error
	SetServiceError(ctx context.Context, id string, status models.ServiceStatus, errorMessage string) error
}
