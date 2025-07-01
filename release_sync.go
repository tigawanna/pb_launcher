package main

import (
	"context"
	"log/slog"
	"pb_launcher/configs"
	"pb_launcher/helpers/serialexecutor"
	download "pb_launcher/internal/download/domain"
)

func RegisterBinaryReleaseSync(
	executor *serialexecutor.SequentialExecutor,
	downloader *download.DownloadUsecase,
	config configs.Config) error {

	releaseSyncTask := serialexecutor.NewTask(
		func(ctx context.Context) {
			if err := downloader.Run(ctx); err != nil {
				slog.Error("error processing GitHub release sync task", "error", err)
			}
		},
		config.GetReleaseSyncInterval(),
		99999,
	)
	return executor.Add(releaseSyncTask)
}
