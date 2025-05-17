package process_test

import (
	"pb_luncher/helpers/process"
	"sync"
	"sync/atomic"
	"testing"
)

func TestProcess_StartAndStop(t *testing.T) {
	errChan := make(chan process.ProcessErrorMessage, 1)
	service := process.New("test-service", "sleep", []string{"2"}, process.WithErrorChan(errChan))

	if err := service.Start(); err != nil {
		t.Fatalf("failed to start service: %v", err)
	}

	if !service.IsRunning() {
		t.Fatalf("service should be running")
	}

	if err := service.Stop(); err != nil {
		t.Fatalf("failed to stop service: %v", err)
	}

	if service.IsRunning() {
		t.Fatalf("service should not be running")
	}
}

func TestProcess_Restart(t *testing.T) {
	errChan := make(chan process.ProcessErrorMessage, 1)
	service := process.New("test-service", "sleep", []string{"2"}, process.WithErrorChan(errChan))

	if err := service.Start(); err != nil {
		t.Fatalf("failed to start service: %v", err)
	}

	if err := service.Restart(); err != nil {
		t.Fatalf("failed to restart service: %v", err)
	}

	if !service.IsRunning() {
		t.Fatalf("service should be running after restart")
	}

	if err := service.Stop(); err != nil {
		t.Fatalf("failed to stop service after restart: %v", err)
	}
}

func TestProcess_AggressiveRestart(t *testing.T) {
	errChan := make(chan process.ProcessErrorMessage, 1)
	service := process.New("test-service", "sleep", []string{"5"}, process.WithErrorChan(errChan))

	// Intentar múltiples reinicios rápidos
	for i := range 5 {
		if err := service.Start(); err != nil {
			t.Fatalf("failed to start service (iteration %d): %v", i, err)
		}

		if !service.IsRunning() {
			t.Fatalf("service should be running (iteration %d)", i)
		}

		if err := service.Restart(); err != nil {
			t.Fatalf("failed to restart service (iteration %d): %v", i, err)
		}

		if !service.IsRunning() {
			t.Fatalf("service should be running after restart (iteration %d)", i)
		}
	}

	// Simular paradas bruscas y reinicios
	for i := range 5 {
		if err := service.Stop(); err != nil {
			t.Fatalf("failed to stop service (iteration %d): %v", i, err)
		}

		if service.IsRunning() {
			t.Fatalf("service should not be running after stop (iteration %d)", i)
		}

		if err := service.Start(); err != nil {
			t.Fatalf("failed to restart service after stop (iteration %d): %v", i, err)
		}

		if !service.IsRunning() {
			t.Fatalf("service should be running after restart (iteration %d)", i)
		}
	}

	// Probar que el canal de errores funcione correctamente
	if err := service.Stop(); err != nil {
		t.Fatalf("failed to stop service at final phase: %v", err)
	}

	// Verificar que no queden errores bloqueando el canal
	select {
	case err := <-errChan:
		t.Fatalf("unexpected error received after stop: %v", err)
	default:
	}

	// Intentar reiniciar después de detener
	if err := service.Restart(); err != nil {
		t.Fatalf("failed to restart service after final stop: %v", err)
	}

	if !service.IsRunning() {
		t.Fatalf("service should be running after final restart")
	}

	// Detener finalmente
	if err := service.Stop(); err != nil {
		t.Fatalf("failed to stop service after aggressive testing: %v", err)
	}

	if service.IsRunning() {
		t.Fatalf("service should not be running after final stop")
	}
}

func TestProcess_ConcurrentStress(t *testing.T) {
	errChan := make(chan process.ProcessErrorMessage, 1)
	service := process.New("test-service", "sleep", []string{"5"}, process.WithErrorChan(errChan))

	var wg sync.WaitGroup
	numGoroutines := 200
	var startCount, stopCount, restartCount int32

	// Prueba de concurrencia con múltiples operaciones simultáneas
	for i := range numGoroutines {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			switch id % 4 {
			case 0:
				// Start
				if err := service.Start(); err != nil {
					t.Errorf("unexpected start error: %v", err)
				} else if service.IsRunning() {
					atomic.AddInt32(&startCount, 1)
				}
			case 1:
				// Stop
				if err := service.Stop(); err != nil {
					t.Errorf("unexpected stop error: %v", err)
				} else {
					atomic.AddInt32(&stopCount, 1)
				}
			case 2:
				// Restart
				if err := service.Restart(); err != nil {
					t.Errorf("unexpected restart error: %v", err)
				} else if service.IsRunning() {
					atomic.AddInt32(&restartCount, 1)
				}
			case 3:
				// Check if running
				_ = service.IsRunning()
			}
		}(i)
	}

	// Espera a que todas las gorutinas terminen
	wg.Wait()

	// Verifica que se hayan logrado hacer algunas operaciones exitosamente
	if startCount == 0 && restartCount == 0 {
		t.Fatalf("service was never successfully started or restarted during the test")
	}

	// Verifica que el servicio esté detenido al final
	if service.IsRunning() {
		t.Fatalf("service should not be running after stress test")
	}

	// Verifica que no haya errores inesperados en el canal
	close(errChan)
	for err := range errChan {
		t.Errorf("unexpected error received: %v", err)
	}

	t.Logf("start: %d, stop: %d, restart: %d", startCount, stopCount, restartCount)
}
