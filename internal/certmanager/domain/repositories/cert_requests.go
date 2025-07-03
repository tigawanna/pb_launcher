package repositories

import (
	"context"
	"errors"
	"pb_launcher/internal/certmanager/domain/models"
)

var ErrCertRequestNotFound = errors.New("cert request not found")

type CertRequestRepository interface {
	DomainsWithHttpsEnabled(ctx context.Context) ([]string, error)
	CreatePending(ctx context.Context, domain string, attempt int) error
	MarkAsApproved(ctx context.Context, id string) error
	MarkAsFailed(ctx context.Context, id, message string) error
	PendingByDomain(ctx context.Context, domain string) ([]models.CertRequest, error)
	LastByDomain(ctx context.Context, domain string) (*models.CertRequest, error)
}
