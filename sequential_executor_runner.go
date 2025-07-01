package main

import (
	"context"
	"log/slog"
	"pb_launcher/helpers/serialexecutor"

	"go.uber.org/fx"
)

func RunSequentialExecutor(lc fx.Lifecycle, executor *serialexecutor.SequentialExecutor) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := executor.Start()
			if err != nil {
				slog.Error("failed to start SequentialExecutor", slog.Any("error", err))
				return err
			}
			slog.Info("sequentialExecutor started successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			err := executor.Stop()
			if err != nil {
				slog.Error("failed to stop SequentialExecutor", slog.Any("error", err))
				return err
			}
			slog.Info("sequentialExecutor stopped successfully")
			return nil
		},
	})
}
