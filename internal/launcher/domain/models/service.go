package models

type ServiceStatus string
type RestartPolicy string

const (
	Idle           ServiceStatus = "idle"            // Created but never started
	Starting       ServiceStatus = "starting"        // In the process of starting
	Running        ServiceStatus = "running"         // Active and running
	Stopping       ServiceStatus = "stopping"        // In the process of stopping
	Stopped        ServiceStatus = "stopped"         // Stopped manually
	StartFailed    ServiceStatus = "start-failed"    // Failed to start
	UnexpectedExit ServiceStatus = "unexpected-exit" // Crashed or stopped unexpectedly
)

const (
	Always    RestartPolicy = "always"     // Always restart, except if manually stopped
	OnFailure RestartPolicy = "on-failure" // Restart only on errors (start-failed, unexpected-exit)
	Never     RestartPolicy = "never"      // Never restart automatically
)

type Service struct {
	ID            string
	Status        ServiceStatus
	RestartPolicy RestartPolicy
}
