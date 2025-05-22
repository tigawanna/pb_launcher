package repositories

import (
	"context"
	"pb_launcher/internal/launcher/domain/models"
)

type ServiceRepository interface {
	Services(ctx context.Context) ([]models.Service, error)
	RunningServices(ctx context.Context) ([]models.Service, error)
	FindService(ctx context.Context, id string) (*models.Service, error)

	MarkServiceStoped(ctx context.Context, id string) error
	MarkServiceFailure(ctx context.Context, id string, errorMessage string) error
	MarkServiceRunning(ctx context.Context, id string, listenIplistenIp, port string) error
	BootCompleted(ctx context.Context, id string) error
}
