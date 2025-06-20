package logstore

import (
	"errors"
	"io"
	"log/slog"
	"sync"
	"time"
)

type StreamType string

const (
	StreamStdout StreamType = "stdout"
	StreamStderr StreamType = "stderr"
)

const maxLogsPerService = 100

type ServiceLog struct {
	ID        int64     `json:"id"`
	ServiceID string    `json:"service_id"`
	Stream    string    `json:"stream"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type ServiceLogDB struct {
	mu    sync.RWMutex
	logs  map[string][]*ServiceLog
	idSeq int64
}

func NewServiceLogDB() (*ServiceLogDB, error) {
	return &ServiceLogDB{
		logs: make(map[string][]*ServiceLog),
	}, nil
}

func (s *ServiceLogDB) InsertLog(serviceID string, stream StreamType, message string) error {
	if stream != StreamStdout && stream != StreamStderr {
		return errors.New("invalid stream type")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.idSeq++
	log := &ServiceLog{
		ID:        s.idSeq,
		ServiceID: serviceID,
		Stream:    string(stream),
		Message:   message,
		Timestamp: time.Now(),
	}
	s.logs[serviceID] = append(s.logs[serviceID], log)

	// Truncar si se excede el máximo
	if len(s.logs[serviceID]) > maxLogsPerService {
		excess := len(s.logs[serviceID]) - maxLogsPerService
		s.logs[serviceID] = s.logs[serviceID][excess:]
	}

	return nil
}

func (s *ServiceLogDB) GetLogsByService(serviceID string) ([]ServiceLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := s.logs[serviceID]
	n := len(list)
	result := make([]ServiceLog, 0, n)
	for i := n - 1; i >= 0; i-- {
		result = append(result, *list[i])
		if len(result) >= maxLogsPerService {
			break
		}
	}
	return result, nil
}

func (s *ServiceLogDB) Cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for serviceID, entries := range s.logs {
		if len(entries) > maxLogsPerService {
			s.logs[serviceID] = entries[len(entries)-maxLogsPerService:]
		}
	}
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
