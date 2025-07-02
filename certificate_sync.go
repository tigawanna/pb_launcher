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
	"pb_launcher/utils"
	"sync/atomic"
)

func RegisterCertificateAutoRenewal(
	provider tlscommon.Provider,
	store tlscommon.Store,
	cfg configs.Config,
	executor *serialexecutor.SequentialExecutor,
) error {
	if !cfg.UseHttps() {
		return nil
	}

	domain := utils.NormalizeWildcardDomain(cfg.GetDomain())
	var initialExecutionComplete atomic.Bool

	certificateTask := serialexecutor.NewTask(
		func(ctx context.Context) {
			defer initialExecutionComplete.Store(true)

			currentCert, err := store.Resolve(domain)
			if err != nil &&
				!errors.Is(err, tlscommon.ErrCertificateNotFound) &&
				!errors.Is(err, tlscommon.ErrInvalidPEM) &&
				!errors.Is(err, tlscommon.ErrCertificateExpired) {
				slog.Error("unexpected error resolving certificate", "domain", domain, "error", err)
				if !initialExecutionComplete.Load() {
					os.Exit(1)
				}
				return
			}

			if err == nil && currentCert.Ttl > cfg.GetMinCertificateTtl() {
				return
			}

			newCert, err := provider.RequestCertificate(domain)
			if err != nil {
				slog.Error("failed to request certificate", "domain", domain, "error", err)
				if !initialExecutionComplete.Load() {
					os.Exit(1)
				}
				return
			}

			if err := store.Store(domain, *newCert); err != nil {
				slog.Error("failed to store certificate", "domain", domain, "error", err)
				if !initialExecutionComplete.Load() {
					os.Exit(1)
				}
				return
			}

			slog.Info("certificate successfully requested and stored", "domain", domain)
		},
		cfg.GetCertificateCheckInterval(),
		math.MaxInt,
	)

	return executor.Add(certificateTask)
}
