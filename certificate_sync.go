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

func RegisterCertificateSync(
	provider tlscommon.Provider,
	storer tlscommon.Store,
	cfg configs.Config,
	executor *serialexecutor.SequentialExecutor,
) error {
	if !cfg.UseHttps() {
		return nil
	}

	domain := utils.NormalizeWildcardDomain(cfg.GetDomain())
	var firstExecutionDone atomic.Bool

	certSyncTask := serialexecutor.NewTask(
		func(ctx context.Context) {
			defer firstExecutionDone.Store(true)

			_, err := storer.Resolve(domain)
			if err == nil {
				return
			}

			if !errors.Is(err, tlscommon.ErrCertificateNotFound) &&
				!errors.Is(err, tlscommon.ErrInvalidPEM) &&
				!errors.Is(err, tlscommon.ErrCertificateExpired) {
				slog.Error("unexpected error resolving certificate", "domain", domain, "error", err)
				if !firstExecutionDone.Load() {
					os.Exit(1)
				}
				return
			}

			cert, err := provider.RequestCertificate(domain)
			if err != nil {
				slog.Error("failed to request certificate", "domain", domain, "error", err)
				if !firstExecutionDone.Load() {
					os.Exit(1)
				}
				return
			}

			if err := storer.Store(domain, *cert); err != nil {
				slog.Error("failed to store certificate", "domain", domain, "error", err)
				if !firstExecutionDone.Load() {
					os.Exit(1)
				}
				return
			}
		},
		cfg.GetCertificateCheckInterval(),
		math.MaxInt,
	)
	return executor.Add(certSyncTask)
}
