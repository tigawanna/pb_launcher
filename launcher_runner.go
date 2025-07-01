package main

import (
	"context"
	"log/slog"
	"pb_launcher/configs"
	"pb_launcher/helpers/serialexecutor"
	launcher "pb_launcher/internal/launcher/domain"
	"sync/atomic"
)

func RegisterLauncherRunner(
	executor *serialexecutor.SequentialExecutor,
	launcherManager *launcher.LauncherManager,
	config configs.Config) error {

	var recoveryDone atomic.Bool

	launcherRunnerTask := serialexecutor.NewTask(
		func(ctx context.Context) {
			if !recoveryDone.Load() {
				if err := launcherManager.RecoveryLastState(ctx); err != nil {
					slog.Error("recovery process failed", "error", err, "task", "launcherRunner")
					return
				}
				recoveryDone.Store(true)
			}
			if err := launcherManager.Run(ctx); err != nil {
				slog.Error("service runner task failed", "error", err, "task", "launcherRunner")
			}
		},
		config.GetCommandCheckInterval(),
		9999,
	)

	return executor.Add(launcherRunnerTask)
}
