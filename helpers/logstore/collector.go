package logstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"slices"
	"sync"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"go.uber.org/fx"
)

type StreamType string

const (
	StreamStdout StreamType = "stdout"
	StreamStderr StreamType = "stderr"
)

const maxLogsPerService = 500

type ServiceLog struct {
	ID        int64     `json:"id"`
	ServiceID string    `json:"service_id"`
	Stream    string    `json:"stream"`
	Message   []byte    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

var _ json.Marshaler = (*ServiceLog)(nil)
var _ json.Marshaler = ServiceLog{}

func (s ServiceLog) MarshalJSON() ([]byte, error) {
	type Alias ServiceLog
	return json.Marshal(&struct {
		Message string `json:"message"`
		Alias
	}{
		Message: string(s.Message),
		Alias:   (Alias)(s),
	})
}

type ServiceLogDB struct {
	mu sync.RWMutex
	db *dbx.DB
}

func NewServiceLogDB(lc fx.Lifecycle, app *pocketbase.PocketBase) (*ServiceLogDB, error) {
	dataDir := app.DataDir()
	dbPath := filepath.Join(dataDir, "service_logs.db")
	dsn := fmt.Sprintf(
		"%s?_pragma=busy_timeout(10000)&_pragma=synchronous(NORMAL)&_pragma=journal_mode(WAL)&_pragma=temp_store(MEMORY)",
		dbPath,
	)
	db, err := dbx.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open service logs database: %w", err)
	}

	const schema = `
	CREATE TABLE IF NOT EXISTS service_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		service_id TEXT NOT NULL,
		stream TEXT NOT NULL,
		message BLOB NOT NULL,
		timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_service_logs_service_id ON service_logs(service_id);
	`

	if _, err := db.NewQuery(schema).Execute(); err != nil {
		return nil, fmt.Errorf("failed to initialize service_logs schema: %w", err)
	}
	lc.Append(fx.StopHook(func() {
		if err := db.Close(); err != nil {
			slog.Error("failed to close service logs database", slog.Any("error", err))
		} else {
			slog.Info("service logs database closed successfully")
		}
	}))
	return &ServiceLogDB{db: db}, nil
}

func (s *ServiceLogDB) InsertLog(serviceID string, stream StreamType, message string) error {
	if stream != StreamStdout && stream != StreamStderr {
		return errors.New("invalid stream type")
	}

	query := `
		INSERT INTO service_logs (service_id, stream, message, timestamp)
		VALUES ({:service_id}, {:stream}, {:message}, CURRENT_TIMESTAMP)
	`

	_, err := s.db.NewQuery(query).
		Bind(dbx.Params{
			"service_id": serviceID,
			"stream":     string(stream),
			"message":    []byte(message),
		}).Execute()

	return err
}

func (s *ServiceLogDB) GetLogsByService(serviceID string, limit int) ([]ServiceLog, error) {
	if limit < 0 {
		limit = 0
	}
	if limit == 0 {
		limit = maxLogsPerService
	}

	query := `
		SELECT id, service_id, stream, message, timestamp
		FROM service_logs
		WHERE service_id = {:service_id}
		ORDER BY id DESC
		LIMIT {:limit}
	`

	var logs []ServiceLog
	err := s.db.NewQuery(query).
		Bind(dbx.Params{
			"service_id": serviceID,
			"limit":      limit,
		}).
		All(&logs)

	slices.Reverse(logs)
	return logs, err
}

func (s *ServiceLogDB) Cleanup() error {
	query := `
		DELETE FROM service_logs
		WHERE id IN (
			SELECT id FROM service_logs
			WHERE service_id = logs.service_id
			ORDER BY timestamp DESC
			LIMIT -1 OFFSET {:max}
		)
		FROM (
			SELECT DISTINCT service_id FROM service_logs
		) AS logs
	`

	_, err := s.db.NewQuery(query).
		Bind(dbx.Params{"max": maxLogsPerService}).
		Execute()

	return err
}

type ServiceLogger struct {
	serviceID string
	stream    StreamType
	logger    *ServiceLogDB
}

var _ io.Writer = (*ServiceLogger)(nil)

func (s *ServiceLogger) Write(p []byte) (int, error) {
	message := string(p)
	if err := s.logger.InsertLog(s.serviceID, s.stream, message); err != nil {
		slog.Error("failed to insert log",
			slog.String("service_id", s.serviceID),
			slog.String("stream", string(s.stream)),
			slog.String("message", message),
			slog.Any("error", err),
		)
		return 0, err
	}
	return len(p), nil
}

func (s *ServiceLogDB) NewWriter(serviceID string, stream StreamType) io.Writer {
	return &ServiceLogger{
		serviceID: serviceID,
		stream:    stream,
		logger:    s,
	}
}
