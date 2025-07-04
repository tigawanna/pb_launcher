package repos

import (
	"context"
	"pb_launcher/collections"
	"pb_launcher/internal/certmanager/domain/models"
	"pb_launcher/internal/certmanager/domain/repositories"
	"pb_launcher/utils"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

type CertRequestRepository struct {
	app *pocketbase.PocketBase
}

var _ repositories.CertRequestRepository = (*CertRequestRepository)(nil)

func NewCertRequestRepository(app *pocketbase.PocketBase) *CertRequestRepository {
	return &CertRequestRepository{app: app}
}

func (r *CertRequestRepository) DomainsWithHttpsEnabled(ctx context.Context) ([]string, error) {
	query := r.app.RecordQuery(collections.ServicesDomains).
		WithContext(ctx).
		Select("domain").
		AndWhere(dbx.NewExp("use_https = 'yes'"))

	var records []*core.Record
	if err := query.All(&records); err != nil {
		return nil, err
	}

	domains := make([]string, 0, len(records))
	for _, rec := range records {
		domains = append(domains, rec.GetString("domain"))
	}
	return domains, nil
}

func (r *CertRequestRepository) CreatePending(ctx context.Context, domain string, attempt int) error {
	collection, err := r.app.FindCollectionByNameOrId(collections.CertRequests)
	if err != nil {
		return err
	}
	record := core.NewRecord(collection)
	record.Set("domain", domain)
	record.Set("status", string(models.CertStatePending))
	record.Set("attempt", attempt)
	return r.app.Save(record)
}

func (r *CertRequestRepository) MarkAsApproved(ctx context.Context, id string) error {
	record, err := r.app.FindRecordById(collections.CertRequests, id)
	if err != nil {
		return err
	}
	record.Set("status", string(models.CertStateApproved))
	record.Set("message", nil)
	record.Set("requested", time.Now())
	return r.app.Save(record)
}

func (r *CertRequestRepository) MarkAsFailed(ctx context.Context, id, message string) error {
	record, err := r.app.FindRecordById(collections.CertRequests, id)
	if err != nil {
		return err
	}
	record.Set("status", string(models.CertStateFailed))
	record.Set("message", message)
	record.Set("requested", time.Now())
	return r.app.Save(record)
}

func (r *CertRequestRepository) PendingByDomain(ctx context.Context, domain string) ([]models.CertRequest, error) {
	exp := dbx.NewExp(
		"domain={:domain} AND status='pending'",
		dbx.Params{"domain": domain},
	)

	query := r.app.RecordQuery(collections.CertRequests).
		WithContext(ctx).
		AndWhere(exp).
		OrderBy("created desc")

	var records []*core.Record
	if err := query.All(&records); err != nil {
		return nil, err
	}

	requests := make([]models.CertRequest, 0, len(records))
	for _, rec := range records {
		requests = append(requests, mapCertRequest(rec))
	}
	return requests, nil
}

func (r *CertRequestRepository) DeletePendingByDomain(ctx context.Context, domain string) error {
	const qry = "DELETE FROM cert_requests WHERE domain = {:domain} AND status = 'pending'"
	_, err := r.app.DB().NewQuery(qry).
		Bind(dbx.Params{"domain": domain}).
		WithContext(ctx).
		Execute()
	return err
}

func (r *CertRequestRepository) LastByDomain(ctx context.Context, domain string) (*models.CertRequest, error) {
	query := r.app.RecordQuery(collections.CertRequests).
		WithContext(ctx).
		AndWhere(dbx.NewExp("domain={:domain}", dbx.Params{"domain": domain})).
		OrderBy("created desc").
		Limit(1)

	var records []*core.Record
	if err := query.All(&records); err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, repositories.ErrCertRequestNotFound
	}

	req := mapCertRequest(records[0])
	return &req, nil
}

func mapCertRequest(rec *core.Record) models.CertRequest {
	return models.CertRequest{
		ID:        rec.Id,
		Domain:    rec.GetString("domain"),
		Status:    models.CertRequestState(rec.GetString("status")),
		Attempt:   rec.GetInt("attempt"),
		Message:   utils.Ptr(rec.GetString("message")),
		Created:   rec.GetDateTime("created").Time(),
		Requested: utils.Ptr(rec.GetDateTime("requested").Time()),
	}
}
