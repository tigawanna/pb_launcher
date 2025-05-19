package domain

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
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

func (lm *LauncherManager) buildArgs(serviceID string) ([]string, string, error) {
	freePort, err := lm.findFreePort()
	if err != nil {
		slog.Error("Failed to find free port for service", "serviceID", serviceID, "error", err)
		return nil, "", err
	}
	pb_data := path.Join(lm.dataDir, serviceID)
	var args = []string{"serve",
		"--dir", path.Join(pb_data, "pb_data"),
		"--hooksDir", path.Join(pb_data, "hooks"),
		"--publicDir", path.Join(pb_data, "public"),
		"--migrationsDir", path.Join(pb_data, "migrations"),
		"--http", freePort,
	}
	return args, freePort, nil
}

func (lm *LauncherManager) startService(ctx context.Context, service models.Service) (string, error) {
	executablePath, err := lm.finder.FindBinary(ctx, service.RepositoryID, service.Version, service.ExecFilePattern)
	if err != nil {
		slog.Error("Failed to find binary", "serviceID", service.ID, "error", err)
		return "", err
	}

	args, listenIp, err := lm.buildArgs(service.ID)
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

	newProcess := process.New(
		service.ID,
		executablePath,
		args,
		process.WithErrorChan(lm.errChan),
	)

	if err := newProcess.Start(); err != nil {
		slog.Error("Failed to start process", "serviceID", service.ID, "error", err)
		return "", err
	}

	lm.processList[service.ID] = newProcess

	return listenIp, err
}

func (lm *LauncherManager) Init(ctx context.Context) error {
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
		vals := strings.Split(listenIp, ":")
		if len(vals) != 2 {
			slog.Error("invalid listenIp format, expected ip:port",
				"listenIp", listenIp,
				"serviceID", service.ID,
			)
			continue
		}
		ip, port := vals[0], vals[1]
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
