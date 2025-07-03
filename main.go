package main

import (
	"log/slog"
	"os"
	"path"
	"pb_launcher/configs"
	"pb_launcher/helpers/logstore"
	"pb_launcher/helpers/serialexecutor"
	"pb_launcher/helpers/unzip"
	"pb_launcher/internal"

	"pb_launcher/internal/certificates"
	"pb_launcher/internal/certmanager"
	"pb_launcher/internal/download"
	"pb_launcher/internal/launcher"
	"pb_launcher/internal/proxy"
	_ "pb_launcher/migrations"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func initializeServer() *pocketbase.PocketBase {
	app := pocketbase.New()
	if err := app.Bootstrap(); err != nil {
		slog.Error("Failed to bootstrap PocketBase", "error", err)
		os.Exit(1)
	}
	return app
}

func main() {
	app := initializeServer()
	rootCmd := createRootCommand(app)
	registerCommands(rootCmd, app)
	executeRootCommand(rootCmd)
}

func createRootCommand(app core.App) *cobra.Command {
	var configFile string
	comand := &cobra.Command{
		Use: path.Base(os.Args[0]),
		Run: func(cmd *cobra.Command, args []string) {
			fx.New(
				fx.Provide(func() (configs.Config, error) {
					return configs.LoadConfigs(configFile)
				}),
				certificates.Module,
				fx.Provide(configs.NewPBServeConfig),
				fx.Provide(unzip.NewUnzip),
				fx.Provide(logstore.NewServiceLogDB),
				fx.Provide(serialexecutor.NewSequentialExecutor),
				fx.Supply(app),
				download.Module,
				launcher.Module,
				proxy.Module,
				certmanager.Module,
				internal.Module, // hooks
				fx.Invoke(
					StartApiServer,
					ServeEmbeddedUI,
					// Tasks
					RegisterCertificateAutoRenewal,
					RegisterCertRequestPlanner,

					RegisterBinaryReleaseSync,
					RegisterLauncherRunner,
					RunSequentialExecutor, // Start Stask Runner
				),
			).Run()
		},
	}
	comand.Flags().StringVarP(&configFile, "config", "c", "", "Path to the config file (yml)")
	return comand
}

func buildUpgradeCommand(migrationsRunner *core.MigrationsRunner) *cobra.Command {
	return &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade the database schema to the latest version",
		Run: func(cmd *cobra.Command, args []string) {
			applied, err := migrationsRunner.Up()
			if err != nil {
				slog.Error("Database upgrade failed", "error", err)
				os.Exit(1)
			}
			if len(applied) == 0 {
				slog.Info("No new migrations to apply")
				return
			}
			for _, file := range applied {
				slog.Info("Migration applied", "file", file)
			}
		},
	}
}

func buildDowngradeCommand(migrationsRunner *core.MigrationsRunner) *cobra.Command {
	return &cobra.Command{
		Use:   "downgrade",
		Short: "Downgrade the database schema to a previous version",
		Run: func(cmd *cobra.Command, args []string) {
			reverted, err := migrationsRunner.Down(1)
			if err != nil {
				slog.Error("Database downgrade failed", "error", err)
				os.Exit(1)
			}
			if len(reverted) == 0 {
				slog.Info("No migration to revert")
				return
			}
			for _, file := range reverted {
				slog.Debug("Reverted migration", "file", file)
			}
		},
	}
}

func registerCommands(rootCmd *cobra.Command, app core.App) {
	migrationsRunner := core.NewMigrationsRunner(app, core.AppMigrations)
	rootCmd.AddCommand(buildUpgradeCommand(migrationsRunner))
	rootCmd.AddCommand(buildDowngradeCommand(migrationsRunner))
}

func executeRootCommand(rootCmd *cobra.Command) {
	if err := rootCmd.Execute(); err != nil {
		slog.Error("Command execution failed", "error", err)
		os.Exit(1)
	}
}
