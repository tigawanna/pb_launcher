package main

import (
	"context"
	"errors"
	"log/slog"
	"math"
	"os"
	"pb_launcher/configs"
	"pb_launcher/helpers/serialexecutor"
	"pb_launcher/internal/certificates/tlscommon"
	certmanager "pb_launcher/internal/certmanager/domain"
	"pb_launcher/utils/domainutil"
	"sync/atomic"
)

func RegisterCertificateAutoRenewal(
	provider tlscommon.Provider,
	store tlscommon.Store,
	cfg configs.Config,
	executor *serialexecutor.SequentialExecutor,
) error {
	if !cfg.IsHttpsEnabled() {
		return nil
	}

	wildcardDomain := domainutil.ToWildcardDomain(cfg.GetDomain())
	var initialExecutionComplete atomic.Bool

	certificateTask := serialexecutor.NewTask(
		func(ctx context.Context) {
			defer initialExecutionComplete.Store(true)

			currentCert, err := store.Resolve(wildcardDomain)
			if err != nil &&
				!errors.Is(err, tlscommon.ErrCertificateNotFound) &&
				!errors.Is(err, tlscommon.ErrInvalidPEM) &&
				!errors.Is(err, tlscommon.ErrCertificateExpired) {
				slog.Error("unexpected error resolving certificate", "domain", wildcardDomain, "error", err)
				if !initialExecutionComplete.Load() {
					os.Exit(1)
				}
				return
			}

			if err == nil && currentCert.Ttl > cfg.GetMinCertificateTtl() {
				return
			}

			newCert, err := provider.RequestCertificate(wildcardDomain)
			if err != nil {
				slog.Error("failed to request certificate", "domain", wildcardDomain, "error", err)
				if !initialExecutionComplete.Load() {
					os.Exit(1)
				}
				return
			}

			if err := store.Store(wildcardDomain, *newCert); err != nil {
				slog.Error("failed to store certificate", "domain", wildcardDomain, "error", err)
				if !initialExecutionComplete.Load() {
					os.Exit(1)
				}
				return
			}

			slog.Info("certificate successfully requested and stored", "domain", wildcardDomain)
		},
		cfg.GetCertificateCheckInterval(),
		math.MaxInt,
	)

	return executor.Add(certificateTask)
}

func RegisterCertRequestPlanner(
	executor *serialexecutor.SequentialExecutor,
	planner *certmanager.CertRequestPlannerUsecase,
	cfg configs.Config) error {

	plannerTask := serialexecutor.NewTask(
		func(ctx context.Context) {
			domains, err := planner.Domains(ctx)
			if err != nil {
				slog.Error("failed to fetch domains", "err", err)
				return
			}
			for _, domain := range domains {
				if err := planner.PostSSLDomainRequest(ctx, domain, true); err != nil {
					slog.Error("failed to schedule cert request",
						"domain", domain,
						"err", err,
					)
				}
			}
		},
		cfg.GetCertRequestPlannerInterval(),
		0,
	)
	return executor.Add(plannerTask)
}
