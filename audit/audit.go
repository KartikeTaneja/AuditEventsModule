// audit/audit.go
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// AuditEvent represents a single user action in the system
type AuditEvent struct {
	Username          string      `json:"username"`
	ActionString      string      `json:"actionString"`
	ExtraMsg          string      `json:"extraMsg,omitempty"`
	EpochTimestampSec int64       `json:"epochTimestampSec"`
	OrgID             int64       `json:"orgId"`
	Metadata          interface{} `json:"metadata,omitempty"`
}

// AuditLogger interface defines the methods for audit logging
type AuditLogger interface {
	CreateAuditEvent(username, actionString, extraMsg string, epochTimestampSec, orgID int64, metadata interface{}) error
	ReadAuditEvents(orgID int64, startEpochSec, endEpochSec int64) ([]AuditEvent, error)
}

// FileAuditLogger implements AuditLogger interface using file storage
type FileAuditLogger struct {
	filePath string
	mu       sync.Mutex
}

// NewFileAuditLogger creates a new FileAuditLogger
func NewFileAuditLogger(filePath string) (*FileAuditLogger, error) {
	// Check if file exists, if not create it
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to create audit log file: %v", err)
		}
		file.Close()
	}

	return &FileAuditLogger{
		filePath: filePath,
	}, nil
}

// CreateAuditEvent logs a user action to the audit log file
func (l *FileAuditLogger) CreateAuditEvent(username, actionString, extraMsg string, epochTimestampSec, orgID int64, metadata interface{}) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// If timestamp is 0, use current time
	if epochTimestampSec == 0 {
		epochTimestampSec = time.Now().Unix()
	}

	event := AuditEvent{
		Username:          username,
		ActionString:      actionString,
		ExtraMsg:          extraMsg,
		EpochTimestampSec: epochTimestampSec,
		OrgID:             orgID,
		Metadata:          metadata,
	}

	// Convert to JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal audit event: %v", err)
	}

	// Append to file
	file, err := os.OpenFile(l.filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open audit log file: %v", err)
	}
	defer file.Close()

	// Add a newline after each JSON object for better readability
	if _, err := file.Write(append(eventJSON, '\n')); err != nil {
		return fmt.Errorf("failed to write to audit log file: %v", err)
	}

	return nil
}

// ReadAuditEvents reads audit events from the log file for a specific organization and time range
func (l *FileAuditLogger) ReadAuditEvents(orgID int64, startEpochSec, endEpochSec int64) ([]AuditEvent, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	file, err := os.Open(l.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open audit log file: %v", err)
	}
	defer file.Close()

	var events []AuditEvent
	scanner := NewJSONScanner(file)

	for scanner.Scan() {
		var event AuditEvent
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			continue // Skip invalid entries
		}

		// Filter by organization ID and time range
		if event.OrgID == orgID &&
			event.EpochTimestampSec >= startEpochSec &&
			(endEpochSec == 0 || event.EpochTimestampSec <= endEpochSec) {
			events = append(events, event)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading audit log: %v", err)
	}

	return events, nil
}

// Helper functions

// For a real implementation, you might want a more sophisticated JSON scanner
// This is a simplified version
type JSONScanner struct {
	file    *os.File
	buffer  []byte
	err     error
	current []byte
}

func NewJSONScanner(file *os.File) *JSONScanner {
	return &JSONScanner{
		file:   file,
		buffer: make([]byte, 4096),
	}
}

func (s *JSONScanner) Scan() bool {
	line := []byte{}
	buf := make([]byte, 1)
	for {
		n, err := s.file.Read(buf)
		if n == 0 || err != nil {
			if len(line) > 0 {
				s.current = line
				return true
			}
			return false
		}
		if buf[0] == '\n' {
			s.current = line
			return true
		}
		line = append(line, buf[0])
	}
}

func (s *JSONScanner) Bytes() []byte {
	return s.current
}

func (s *JSONScanner) Err() error {
	return s.err
}