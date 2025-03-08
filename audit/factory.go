// audit/factory.go
package audit

import (
	"fmt"
	"sync"
)

// LoggerType defines the type of audit logger to use
type LoggerType string

const (
	FileLoggerType LoggerType = "file"
	DBLoggerType   LoggerType = "db"
)

var (
	loggerInstance AuditLogger
	loggerMutex    sync.Mutex
)

// InitAuditLogger initializes the audit logger with the specified type and configuration
func InitAuditLogger(loggerType LoggerType, config map[string]string) error {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()

	var err error

	switch loggerType {
	case FileLoggerType:
		filePath, ok := config["filePath"]
		if !ok {
			return fmt.Errorf("file path not provided for file logger")
		}
		loggerInstance, err = NewFileAuditLogger(filePath)
	case DBLoggerType:
		dbPath, ok := config["dbPath"]
		if !ok {
			return fmt.Errorf("database path not provided for DB logger")
		}
		loggerInstance, err = NewDBAuditLogger(dbPath)
	default:
		return fmt.Errorf("unsupported logger type: %s", loggerType)
	}

	return err
}

// GetAuditLogger returns the initialized audit logger
func GetAuditLogger() (AuditLogger, error) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()

	if loggerInstance == nil {
		return nil, fmt.Errorf("audit logger not initialized, call InitAuditLogger first")
	}

	return loggerInstance, nil
}

// CreateAuditEvent is a convenience function to create an audit event without getting the logger
func CreateAuditEvent(username, actionString, extraMsg string, epochTimestampSec, orgID int64, metadata interface{}) error {
	logger, err := GetAuditLogger()
	if err != nil {
		return err
	}

	return logger.CreateAuditEvent(username, actionString, extraMsg, epochTimestampSec, orgID, metadata)
}

// ReadAuditEvents is a convenience function to read audit events without getting the logger
func ReadAuditEvents(orgID int64, startEpochSec, endEpochSec int64) ([]AuditEvent, error) {
	logger, err := GetAuditLogger()
	if err != nil {
		return nil, err
	}

	return logger.ReadAuditEvents(orgID, startEpochSec, endEpochSec)
}