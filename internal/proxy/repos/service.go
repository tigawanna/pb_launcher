package repos

import (
	"context"
	"database/sql"
	"errors"
	"pb_launcher/collections"
	"pb_launcher/internal/proxy/domain/dtos"
	"pb_launcher/internal/proxy/domain/repositories"

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

func (r *ServiceRepository) FindRunningServiceByID(ctx context.Context, id string) (*dtos.RunningServiceDto, error) {
	record, err := r.app.FindRecordById(collections.Services, id, func(q *dbx.SelectQuery) error {
		q.AndWhere(dbx.NewExp("(deleted IS NULL OR deleted = '') AND status = 'running'"))
		return nil
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repositories.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &dtos.RunningServiceDto{
		ID:   record.Id,
		IP:   record.GetString("ip"),
		Port: record.GetInt("port"),
	}, nil
}
