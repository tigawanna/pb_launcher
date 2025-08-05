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

type ProxyEntriesRepository struct {
	app *pocketbase.PocketBase
}

var _ repositories.ProxyEntriesRepository = (*ProxyEntriesRepository)(nil)

func NewProxyEntriesRepository(app *pocketbase.PocketBase) *ProxyEntriesRepository {
	return &ProxyEntriesRepository{app: app}
}

func (r *ProxyEntriesRepository) FindEnabledProxyEntryByID(ctx context.Context, id string) (*dtos.ProxyEntryDto, error) {
	record, err := r.app.FindRecordById(collections.ProxyEntries, id, func(q *dbx.SelectQuery) error {
		q.AndWhere(dbx.NewExp("(deleted IS NULL OR deleted = '') AND enabled = 'yes'"))
		return nil
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repositories.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &dtos.ProxyEntryDto{
		ID:        record.Id,
		TargetUrl: record.GetString("target_url"),
	}, nil
}
