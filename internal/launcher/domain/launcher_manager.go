package domain

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os/exec"
	"path"
	"pb_launcher/configs"
	"pb_launcher/helpers/process"
	"pb_launcher/internal/launcher/domain/models"
	"pb_launcher/internal/launcher/domain/repositories"
	"pb_launcher/internal/launcher/domain/services"
	"strings"
	"sync"
)

type LauncherManager struct {
	sync.RWMutex
	dataDir    string
	repository repositories.ServiceRepository
	finder     services.BinaryFinder
	//
	processList map[string]*process.Process
	errChan     chan process.ProcessErrorMessage
}

func NewLauncherManager(
	repository repositories.ServiceRepository,
	finder services.BinaryFinder,
	c *configs.Configs,
) *LauncherManager {
	return &LauncherManager{
		repository:  repository,
		finder:      finder,
		dataDir:     c.DataDir,
		processList: make(map[string]*process.Process),
		errChan:     make(chan process.ProcessErrorMessage),
	}
}

func (lm *LauncherManager) findFreePort() (string, error) {
	listener, err := net.Listen("tcp", "127.0.0.2:0")
	if err != nil {
		return "", fmt.Errorf("failed to find free port: %w", err)
	}
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)
	return fmt.Sprintf("127.0.0.2:%d", addr.Port), nil
}

func (lm *LauncherManager) buildArgs(serviceID string) ([]string, error) {
	pb_data := path.Join(lm.dataDir, serviceID)
	return []string{
		"--dir", path.Join(pb_data, "pb_data"),
		"--hooksDir", path.Join(pb_data, "hooks"),
		"--publicDir", path.Join(pb_data, "public"),
		"--migrationsDir", path.Join(pb_data, "migrations"),
	}, nil
}

// initializeBootUser sets up the initial boot user for the service instance.
func (lm *LauncherManager) initializeBootUser(ctx context.Context,
	service models.Service, binaryPath string, baseArgs []string) error {

	if service.BootCompleted {
		return nil
	}

	args := append(baseArgs, "superuser", "create", service.BootUserEmail, service.BootUserPassword)
	cmd := exec.CommandContext(ctx, binaryPath, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		errorMessage := fmt.Sprintf("boot user :%s", err.Error())
		if updateErr := lm.repository.SetServiceError(ctx, service.ID, errorMessage); updateErr != nil {
			slog.Error("failed to update service error status",
				"serviceID", service.ID,
				"error", updateErr,
				"originalError", errorMessage,
			)
		}
		slog.Error("failed to initialize boot user",
			"service", service.ID,
			"email", service.BootUserEmail,
			"output", string(output),
			"error", err,
		)
		return err
	}

	if err := lm.repository.BootCompleted(ctx, service.ID); err != nil {
		slog.Error("failed to update boot completed flag",
			"service", service.ID,
			"error", err,
		)
		return err
	}

	return nil
}

func (lm *LauncherManager) startService(ctx context.Context, service models.Service) (string, error) {
	executablePath, err := lm.finder.FindBinary(ctx, service.RepositoryID, service.Version, service.ExecFilePattern)
	if err != nil {
		slog.Error("Failed to find binary", "serviceID", service.ID, "error", err)
		return "", err
	}
	listenIp, err := lm.findFreePort()
	if err != nil {
		slog.Error("Failed to find free port for service", "serviceID", service.ID, "error", err)
		return "", err
	}

	baseArgs, err := lm.buildArgs(service.ID)
	if err != nil {
		slog.Error("Failed to build arguments", "serviceID", service.ID, "error", err)
		return "", err
	}

	if existingProcess, exists := lm.processList[service.ID]; exists {
		if existingProcess.IsRunning() {
			slog.Info("Stopping existing process", "serviceID", service.ID)
			if err := existingProcess.Stop(); err != nil {
				slog.Error("Failed to stop existing process", "serviceID", service.ID, "error", err)
				return "", err
			}
		}
		delete(lm.processList, service.ID)
	}

	if err := lm.initializeBootUser(ctx, service, executablePath, baseArgs); err != nil {
		return "", err
	}

	serveArgs := append([]string{"serve"}, append(baseArgs, "--http", listenIp)...)
	newProcess := process.New(
		service.ID,
		executablePath,
		serveArgs,
		process.WithErrorChan(lm.errChan),
	)

	if err := newProcess.Start(); err != nil {
		slog.Error("Failed to start process", "serviceID", service.ID, "error", err)
		return "", err
	}

	lm.processList[service.ID] = newProcess

	return listenIp, err
}

func (lm *LauncherManager) parseIPPort(addr string) (string, string, error) {
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid address format: %s", addr)
	}
	return parts[0], parts[1], nil
}

// RecoveryLastState restores and starts all services that were active
// before pb_launcher was shut down.
func (lm *LauncherManager) RecoveryLastState(ctx context.Context) error {
	lm.Lock()
	defer lm.Unlock()
	go lm.handleServiceErrors()
	services, err := lm.repository.RunningServices(ctx)
	if err != nil {
		slog.Error("Failed to retrieve running services", "error", err)
		return err
	}

	for _, service := range services {
		listenIp, err := lm.startService(ctx, service)
		if err != nil {
			slog.Error("failed to start service",
				"serviceID", service.ID,
				"error", err,
			)

			continue
		}
		ip, port, err := lm.parseIPPort(listenIp)
		if err != nil {
			slog.Error("invalid listenIp format, expected ip:port",
				"listenIp", listenIp,
				"serviceID", service.ID,
			)
			continue
		}
		if err := lm.repository.SetServiceRunning(ctx, service.ID, ip, port); err != nil {
			slog.Error("failed to update service status to running",
				"serviceID", service.ID,
				"ip", ip,
				"port", port,
				"error", err,
			)
		}
	}

	return nil
}

func (lm *LauncherManager) handleServiceErrors() {
	for serviceErr := range lm.errChan {
		ctx := context.Background()
		var errorMessage string
		if serviceErr.Error != nil {
			errorMessage = serviceErr.Error.Error()
		}

		if err := lm.repository.SetServiceError(ctx, serviceErr.ID, errorMessage); err != nil {
			slog.Error("failed to update service falta error  status",
				"serviceID", serviceErr.ID,
				"error", err,
				"originalError", errorMessage,
			)
		}
	}
}

func (lm *LauncherManager) Run(ctx context.Context) error {
	fmt.Println("===========> LAUNCH RUN")
	return nil
}

func (lm *LauncherManager) Stop() error {
	lm.Lock()
	defer lm.Unlock()

	var wg sync.WaitGroup
	var mu sync.Mutex
	var combinedErr error

	collectError := func(err error) {
		mu.Lock()
		defer mu.Unlock()
		combinedErr = errors.Join(combinedErr, err)
	}

	for _, proc := range lm.processList {
		wg.Add(1)
		go func(p *process.Process) {
			defer wg.Done()
			if !p.IsRunning() {
				return
			}
			if err := p.Stop(); err != nil {
				collectError(err)
			}
		}(proc)
	}

	wg.Wait()
	close(lm.errChan)
	return combinedErr
}
