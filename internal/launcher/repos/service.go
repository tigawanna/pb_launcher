package repos

import (
	"context"
	"fmt"
	"log/slog"
	"pb_luncher/collections"
	"pb_luncher/internal/launcher/domain/models"
	"pb_luncher/internal/launcher/domain/repositories"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
)

type ServiceRepository struct {
	app *pocketbase.PocketBase
}

var _ repositories.ServiceRepository = (*ServiceRepository)(nil)

func NewServiceRepository(app *pocketbase.PocketBase) *ServiceRepository {
	return &ServiceRepository{app: app}
}

// FindAll implements repositories.ServiceRepository.
func (s *ServiceRepository) FindAll(ctx context.Context) ([]models.Service, error) {
	services, err := s.app.FindAllRecords(collections.Services, dbx.NewExp("deleted_at IS NULL"))
	if err != nil {
		slog.Error("failed to find all services", "error", err, "collection", collections.Services)
		return nil, err
	}
	var results = make([]models.Service, 0, len(services))
	for _, s := range services {
		results = append(results, models.Service{
			ID:            s.GetString("id"),
			Status:        models.ServiceStatus(s.GetString("status")),
			RestartPolicy: models.RestartPolicy(s.GetString("restart_policy")),
		})
	}
	return results, nil
}

// SetServiceError implements repositories.ServiceRepository.
func (s *ServiceRepository) SetServiceError(ctx context.Context, id string, status models.ServiceStatus, errorMessage string) error {
	validErrorStatuses := map[models.ServiceStatus]bool{
		models.StartFailed:    true,
		models.UnexpectedExit: true,
	}

	if !validErrorStatuses[status] {
		return fmt.Errorf("invalid service status for error: %s", status)
	}

	record, err := s.app.FindRecordById(collections.Services, id)
	if err != nil {
		return err
	}

	record.Set("status", status)
	record.Set("error_message", errorMessage)

	if err := s.app.Save(record); err != nil {
		return err
	}

	return nil
}

func (s *ServiceRepository) UpdateServiceStatus(ctx context.Context, id string, status models.ServiceStatus) error {
	validStatuses := map[models.ServiceStatus]bool{
		models.Idle:           true,
		models.Starting:       true,
		models.Running:        true,
		models.Stopping:       true,
		models.Stopped:        true,
		models.StartFailed:    true,
		models.UnexpectedExit: true,
	}

	if !validStatuses[status] {
		return fmt.Errorf("invalid service status: %s", status)
	}

	record, err := s.app.FindRecordById(collections.Services, id)
	if err != nil {
		return err
	}

	record.Set("status", status)

	if status == models.Running {
		record.Set("last_started_at", time.Now())
		record.Set("error_message", nil)
	}

	if status == models.StartFailed || status == models.UnexpectedExit {
		record.Set("last_started_at", nil)
	}

	if err := s.app.Save(record); err != nil {
		return err
	}

	return nil
}
