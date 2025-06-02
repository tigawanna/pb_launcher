package repos

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"pb_launcher/collections"
	"pb_launcher/internal/launcher/domain/models"
	"pb_launcher/internal/launcher/domain/repositories"
	"regexp"
	"strings"
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

func (s *ServiceRepository) services(ids ...string) ([]models.Service, error) {
	qry := `
		select 
			s.id, 
			s.status, 
			s.restart_policy, 
			r.version, 
			r.repository, 
			rpo.exec_file_pattern,
			s.boot_completed,
			s.boot_user_email,
			s.boot_user_password,
			s.deleted
		from services s
		inner join releases r on s."release" = r.id
		inner join repositories rpo on rpo.id = r.repository`

	var quoted []string
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		quoted = append(quoted, fmt.Sprintf("'%s'", id))
	}
	if len(quoted) > 0 {
		qry += fmt.Sprintf(" and s.id in (%s)", strings.Join(quoted, ","))
	}
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
		bootCompleted, _ := row["boot_completed"]
		bootUserEmail, _ := row["boot_user_email"]
		bootUserPassword, _ := row["boot_user_password"]
		deleted, _ := row["deleted"]

		ExecFilePattern, err := regexp.Compile(execPattern.String)
		if err != nil {
			slog.Warn("invalid exec file pattern", "error", err, "pattern", execPattern)
			continue
		}

		services = append(services, models.Service{
			ID:               id.String,
			Status:           models.ServiceStatus(status.String),
			RestartPolicy:    models.RestartPolicy(restartPolicy.String),
			Version:          version.String,
			RepositoryID:     repository.String,
			ExecFilePattern:  ExecFilePattern,
			BootUserEmail:    bootUserEmail.String,
			BootUserPassword: bootUserPassword.String,
			BootCompleted:    bootCompleted.String == "yes",
			Deleted:          deleted.String,
		})
	}

	return services, nil
}

// Services implements repositories.ServiceRepository.
func (s *ServiceRepository) Services(ctx context.Context) ([]models.Service, error) {
	return s.services()
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

// FindService implements repositories.ServiceRepository.
func (s *ServiceRepository) FindService(ctx context.Context, id string) (*models.Service, error) {
	services, err := s.services(id)
	if err != nil {
		return nil, err
	}
	if len(services) == 0 {
		return nil, fmt.Errorf("service not found: %s", id)
	}
	return &services[0], nil
}

// MarkServiceStoped implements repositories.ServiceRepository.
func (s *ServiceRepository) MarkServiceStoped(ctx context.Context, id string) error {

	record, err := s.app.FindRecordById(collections.Services, id)
	if err != nil {
		return err
	}

	record.Set("status", string(models.Stopped))
	record.Set("error_message", nil)

	if err := s.app.Save(record); err != nil {
		return err
	}

	return nil
}

// MarkServiceFailure implements repositories.ServiceRepository.
func (s *ServiceRepository) MarkServiceFailure(ctx context.Context, id string, errorMessage string) error {

	record, err := s.app.FindRecordById(collections.Services, id)
	if err != nil {
		return err
	}

	record.Set("status", string(models.Failure))
	record.Set("error_message", errorMessage)

	if err := s.app.Save(record); err != nil {
		return err
	}

	return nil
}

// MarkServiceRunning implements repositories.ServiceRepository.
func (s *ServiceRepository) MarkServiceRunning(ctx context.Context, id, listenIp, port string) error {

	record, err := s.app.FindRecordById(collections.Services, id)
	if err != nil {
		return err
	}

	record.Set("status", string(models.Running))
	record.Set("last_started", time.Now())
	record.Set("error_message", nil)
	record.Set("ip", listenIp)
	record.Set("port", port)

	if err := s.app.Save(record); err != nil {
		return err
	}

	return nil
}

// BootCompleted implements repositories.ServiceRepository.
func (s *ServiceRepository) BootCompleted(ctx context.Context, id string) error {
	record, err := s.app.FindRecordById(collections.Services, id)
	if err != nil {
		return err
	}
	record.Set("boot_completed", "yes")
	if err := s.app.Save(record); err != nil {
		return err
	}
	return nil
}
