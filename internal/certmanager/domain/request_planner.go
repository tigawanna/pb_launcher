package domain

import (
	"context"
	"errors"
	"pb_launcher/configs"
	"pb_launcher/internal/certificates/tlscommon"
	"pb_launcher/internal/certmanager/domain/models"
	"pb_launcher/internal/certmanager/domain/repositories"
	"time"
)

type CertRequestPlannerUsecase struct {
	repository  repositories.CertRequestRepository
	store       tlscommon.Store
	minTTL      time.Duration
	maxAttempts int
}

func NewCertRequestPlannerUsecase(
	repository repositories.CertRequestRepository,
	store tlscommon.Store,
	conf configs.Config,
) *CertRequestPlannerUsecase {
	return &CertRequestPlannerUsecase{
		repository:  repository,
		store:       store,
		minTTL:      conf.GetMinCertificateTtl(),
		maxAttempts: conf.GetMaxDomainCertAttempts(),
	}
}

func (uc *CertRequestPlannerUsecase) Domains(ctx context.Context) ([]string, error) {
	return uc.repository.DomainsWithHttpsEnabled(ctx)
}

// PostSSLDomainRequest schedules a task to request an SSL certificate for the given domain.
// If checkMaxAttempts is true, the function validates whether the maximum number of attempts has been exceeded before scheduling
func (uc *CertRequestPlannerUsecase) PostSSLDomainRequest(ctx context.Context, domain string, checkMaxAttempts bool) error {
	pending, err := uc.repository.PendingByDomain(ctx, domain)
	if err != nil {
		return err
	}
	if len(pending) > 0 {
		return nil // already has pending request
	}

	attempt := 1
	if checkMaxAttempts {
		last, err := uc.repository.LastByDomain(ctx, domain)
		if err != nil && !errors.Is(err, repositories.ErrCertRequestNotFound) {
			return err
		}
		if last != nil && last.Status == models.CertStateFailed {
			if last.Attempt >= uc.maxAttempts {
				return nil // exceeded max attempts
			}
			attempt = last.Attempt + 1
		}
	}

	currentCert, err := uc.store.Resolve(domain)
	if err != nil &&
		!errors.Is(err, tlscommon.ErrCertificateNotFound) &&
		!errors.Is(err, tlscommon.ErrInvalidPEM) &&
		!errors.Is(err, tlscommon.ErrCertificateExpired) {
		return err
	}

	if err == nil && currentCert.GetTTL() > uc.minTTL {
		return nil // valid cert, no need to renew
	}

	return uc.repository.CreatePending(ctx, domain, attempt)
}
