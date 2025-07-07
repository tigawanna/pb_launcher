package domain

import (
	"context"
	"fmt"
	"log/slog"
	http01 "pb_launcher/internal/certificates/http_01"
	"pb_launcher/internal/certificates/tlscommon"
	"pb_launcher/internal/certmanager/domain/models"
	"pb_launcher/internal/certmanager/domain/repositories"
	"time"
)

type CertRequestExecutorUsecase struct {
	service    *http01.HTTP01TLSCertificateRequestService
	repository repositories.CertRequestRepository
	store      tlscommon.Store
}

func NewCertRequestExecutorUsecase(
	service *http01.HTTP01TLSCertificateRequestService,
	repository repositories.CertRequestRepository,
	store tlscommon.Store,
) *CertRequestExecutorUsecase {
	return &CertRequestExecutorUsecase{
		service:    service,
		repository: repository,
		store:      store,
	}
}

func (c *CertRequestExecutorUsecase) pendingToExecute(ctx context.Context) ([]models.CertRequest, error) {
	requests, err := c.repository.Pending(ctx)
	if err != nil {
		return nil, err
	}

	pending := []models.CertRequest{}
	for _, r := range requests {
		if r.NotBefore == nil || !r.NotBefore.After(time.Now()) {
			pending = append(pending, r)
		}
	}
	return pending, nil
}

func (c *CertRequestExecutorUsecase) processRequest(ctx context.Context, req models.CertRequest) error {
	cert, err := c.service.RequestCertificate(req.Domain)
	if err != nil {
		if markErr := c.repository.MarkAsFailed(ctx, req.ID, err.Error()); markErr != nil {
			slog.Warn("failed to mark certificate request as failed", "error", markErr, "request_id", req.ID)
		}
		return fmt.Errorf("certificate request failed for domain %s: %w", req.Domain, err)
	}

	if err := c.store.Store(req.Domain, *cert); err != nil {
		if markErr := c.repository.MarkAsFailed(ctx, req.ID, err.Error()); markErr != nil {
			slog.Warn("failed to mark certificate request as failed", "error", markErr, "request_id", req.ID)
		}
		slog.Error("failed to store certificate", "domain", req.Domain, "error", err)
		return fmt.Errorf("failed to store certificate for domain %s: %w", req.Domain, err)
	}

	if err := c.repository.MarkAsApproved(ctx, req.ID); err != nil {
		slog.Warn("failed to mark certificate request as approved", "request_id", req.ID, "error", err)
	}

	return nil
}

func (c *CertRequestExecutorUsecase) ExecutePlan(ctx context.Context) error {
	requests, err := c.pendingToExecute(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve pending certificate requests: %w", err)
	}

	for _, req := range requests {
		if err := c.processRequest(ctx, req); err != nil {
			slog.Error("failed to process certificate request", "request_id", req.ID, "domain", req.Domain, "error", err)
		}
	}

	return nil
}
