package repos

import (
	"context"
	"log"
	"log/slog"
	"pb_launcher/collections"
	"pb_launcher/internal/launcher/domain/models"
	"pb_launcher/internal/launcher/domain/repositories"
	"regexp"
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

// Services implements repositories.ServiceRepository.
func (s *ServiceRepository) Services(ctx context.Context) ([]models.Service, error) {
	const qry = `
		select 
			s.id, 
			s.status, 
			s.restart_policy, 
			r.version, 
			r.repository, 
			rpo.exec_file_pattern
		from services s
		inner join releases r on s."release" = r.id
		inner join repositories rpo on rpo.id = r.repository
	`

	db := s.app.DB()

	results := []dbx.NullStringMap{}
	if err := db.NewQuery(qry).All(&results); err != nil {
		log.Fatal(err)
	}

	services := make([]models.Service, 0, len(results))
	for _, row := range results {
		id, _ := row["id"]
		status, _ := row["status"]
		restartPolicy, _ := row["restart_policy"]
		version, _ := row["version"]
		repository, _ := row["repository"]
		execPattern, _ := row["exec_file_pattern"]

		ExecFilePattern, err := regexp.Compile(execPattern.String)
		if err != nil {
			slog.Warn("invalid exec file pattern", "error", err, "pattern", execPattern)
			continue
		}

		services = append(services, models.Service{
			ID:              id.String,
			Status:          models.ServiceStatus(status.String),
			RestartPolicy:   models.RestartPolicy(restartPolicy.String),
			Version:         version.String,
			RepositoryID:    repository.String,
			ExecFilePattern: ExecFilePattern,
		})
	}

	return services, nil
}

// RunningServices implements repositories.ServiceRepository.
func (s *ServiceRepository) RunningServices(ctx context.Context) ([]models.Service, error) {
	services, err := s.Services(ctx)
	if err != nil {
		return nil, err
	}
	results := []models.Service{}
	for _, service := range services {
		if service.Status == models.Running {
			results = append(results, service)
		}
	}
	return results, nil
}

// SetServiceError implements repositories.ServiceRepository.
func (s *ServiceRepository) SetServiceError(ctx context.Context, id string, errorMessage string) error {

	record, err := s.app.FindRecordById(collections.Services, id)
	if err != nil {
		return err
	}

	record.Set("status", string(models.Stopped))
	record.Set("error_message", errorMessage)

	if err := s.app.Save(record); err != nil {
		return err
	}

	return nil
}

// SetServiceRunning implements repositories.ServiceRepository.
func (s *ServiceRepository) SetServiceRunning(ctx context.Context, id, listenIp, port string) error {

	record, err := s.app.FindRecordById(collections.Services, id)
	if err != nil {
		return err
	}

	record.Set("status", string(models.Running))
	record.Set("last_started_at", time.Now())
	record.Set("error_message", nil)
	record.Set("ip", listenIp)
	record.Set("port", port)

	if err := s.app.Save(record); err != nil {
		return err
	}

	return nil
}
