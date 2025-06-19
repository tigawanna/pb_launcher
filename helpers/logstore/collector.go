package logstore

import (
	"errors"
	"io"
	"time"

	"github.com/pocketbase/dbx"
)

type StreamType string

const (
	StreamStdout StreamType = "stdout"
	StreamStderr StreamType = "stderr"
)

type ServiceLog struct {
	ID        int64     `db:"id"`
	ServiceID string    `db:"service_id"`
	Stream    string    `db:"stream"`
	Message   string    `db:"message"`
	Timestamp time.Time `db:"timestamp"`
}

type ServiceLogDB struct {
	db *dbx.DB
}

func NewServiceLogDB() (*ServiceLogDB, error) {
	dsn := "file::memory:?_pragma=busy_timeout(10000)&_pragma=synchronous(OFF)&_pragma=temp_store(MEMORY)"
	db, err := dbx.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	const schema = `
	CREATE TABLE IF NOT EXISTS service_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		service_id TEXT NOT NULL,
		stream TEXT CHECK(stream IN ('stdout', 'stderr')) NOT NULL,
		message TEXT NOT NULL,
		timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.NewQuery(schema).Execute(); err != nil {
		return nil, err
	}

	return &ServiceLogDB{db: db}, nil
}

func (s *ServiceLogDB) InsertLog(serviceID string, stream StreamType, message string) error {
	if stream != StreamStdout && stream != StreamStderr {
		return errors.New("invalid stream type")
	}

	query := `
		INSERT INTO service_logs (service_id, stream, message)
		VALUES (:service_id, :stream, :message)
	`

	_, err := s.db.NewQuery(query).Bind(dbx.Params{
		"service_id": serviceID,
		"stream":     stream,
		"message":    message,
	}).Execute()

	return err
}

func (s *ServiceLogDB) GetLogsByService(serviceID string) ([]ServiceLog, error) {
	var logs []ServiceLog
	err := s.db.Select().From("service_logs").
		Where(dbx.HashExp{"service_id": serviceID}).
		OrderBy("timestamp ASC").
		All(&logs)
	return logs, err
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
