// audit/audit_test.go
package audit

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFileAuditLogger(t *testing.T) {
	// Create temporary file for testing
	tempDir, err := os.MkdirTemp("", "audit_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logPath := filepath.Join(tempDir, "audit.log")
	
	// Create logger
	logger, err := NewFileAuditLogger(logPath)
	if err != nil {
		t.Fatalf("Failed to create file audit logger: %v", err)
	}
	
	// Test creating audit events
	testCreateAuditEvents(t, logger)
	
	// Test reading audit events
	testReadAuditEvents(t, logger)
}

func TestDBAuditLogger(t *testing.T) {
	// Create temporary database for testing
	tempDir, err := os.MkdirTemp("", "audit_test_db")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "audit.db")
	
	// Create logger
	logger, err := NewDBAuditLogger(dbPath)
	if err != nil {
		t.Fatalf("Failed to create DB audit logger: %v", err)
	}
	defer func() {
		if dbLogger, ok := logger.(*DBAuditLogger); ok {
			dbLogger.Close()
		}
	}()
	
	// Test creating audit events
	testCreateAuditEvents(t, logger)
	
	// Test reading audit events
	testReadAuditEvents(t, logger)
}

func TestFactoryPattern(t *testing.T) {
	// Create temporary directories
	tempDir, err := os.MkdirTemp("", "audit_factory_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logPath := filepath.Join(tempDir, "audit.log")
	dbPath := filepath.Join(tempDir, "audit.db")
	
	// Test file logger initialization
	err = InitAuditLogger(FileLoggerType, map[string]string{
		"filePath": logPath,
	})
	if err != nil {
		t.Fatalf("Failed to initialize file logger: %v", err)
	}
	
	logger, err := GetAuditLogger()
	if err != nil {
		t.Fatalf("Failed to get initialized logger: %v", err)
	}
	
	// Verify it's a file logger
	_, ok := logger.(*FileAuditLogger)
	if !ok {
		t.Fatal("Expected FileAuditLogger, got something else")
	}
	
	// Test DB logger initialization
	err = InitAuditLogger(DBLoggerType, map[string]string{
		"dbPath": dbPath,
	})
	if err != nil {
		t.Fatalf("Failed to initialize DB logger: %v", err)
	}
	
	logger, err = GetAuditLogger()
	if err != nil {
		t.Fatalf("Failed to get initialized logger: %v", err)
	}
	
	// Verify it's a DB logger
	dbLogger, ok := logger.(*DBAuditLogger)
	if !ok {
		t.Fatal("Expected DBAuditLogger, got something else")
	}
	
	// Clean up
	dbLogger.Close()
}

func TestConvenienceFunctions(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "audit_convenience_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logPath := filepath.Join(tempDir, "audit.log")
	
	// Initialize logger
	err = InitAuditLogger(FileLoggerType, map[string]string{
		"filePath": logPath,
	})
	if err != nil {
		t.Fatalf("Failed to initialize file logger: %v", err)
	}
	
	// Test convenience function for creating events
	err = CreateAuditEvent("testuser", ActionUserLogin, "Test login", time.Now().Unix(), 123, nil)
	if err != nil {
		t.Fatalf("Failed to create audit event: %v", err)
	}
	
	// Test convenience function for reading events
	events, err := ReadAuditEvents(123, 0, time.Now().Unix() + 1000)
	if err != nil {
		t.Fatalf("Failed to read audit events: %v", err)
	}
	
	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}
	
	if events[0].Username != "testuser" || events[0].ActionString != ActionUserLogin {
		t.Fatalf("Event data mismatch: %+v", events[0])
	}
}

// Helper functions for testing both logger implementations
func testCreateAuditEvents(t *testing.T, logger AuditLogger) {
	now := time.Now().Unix()
	testCases := []struct {
		username      string
		actionString  string
		extraMsg      string
		timestamp     int64
		orgID         int64
		metadata      interface{}
	}{
		{"alice", ActionUserLogin, "Login from 192.168.1.1", now, 123, nil},
		{"bob", ActionDashboardCreate, "Created dashboard 'System Overview'", now + 60, 123, 
			map[string]interface{}{"dashboardId": "dash-123", "title": "System Overview"}},
		{"charlie", ActionIndexDelete, "Deleted index 'logs-2023'", now + 120, 456, nil},
	}
	
	for _, tc := range testCases {
		err := logger.CreateAuditEvent(tc.username, tc.actionString, tc.extraMsg, tc.timestamp, tc.orgID, tc.metadata)
		if err != nil {
			t.Fatalf("Failed to create audit event: %v", err)
		}
	}
}

func testReadAuditEvents(t *testing.T, logger AuditLogger) {
	// Test reading events for a specific org
	events, err := logger.ReadAuditEvents(123, 0, time.Now().Unix() + 1000)
	if err != nil {
		t.Fatalf("Failed to read audit events: %v", err)
	}
	
	if len(events) != 2 {
		t.Fatalf("Expected 2 events for org 123, got %d", len(events))
	}
	
	// Test reading events for another org
	events, err = logger.ReadAuditEvents(456, 0, time.Now().Unix() + 1000)
	if err != nil {
		t.Fatalf("Failed to read audit events: %v", err)
	}
	
	if len(events) != 1 {
		t.Fatalf("Expected 1 event for org 456, got %d", len(events))
	}
	
	// Test time range filtering
	now := time.Now().Unix()
	events, err = logger.ReadAuditEvents(123, now + 30, now + 90)
	if err != nil {
		t.Fatalf("Failed to read audit events: %v", err)
	}
	
	if len(events) != 1 {
		t.Fatalf("Expected 1 event in time range, got %d", len(events))
	}
	
	// Verify metadata
	if events[0].ActionString != ActionDashboardCreate {
		t.Fatalf("Expected dashboard create action, got %s", events[0].ActionString)
	}
	
	metadataMap, ok := events[0].Metadata.(map[string]interface{})
	if !ok {
		// For DB logger, we may need to unmarshal JSON
		if jsonStr, ok := events[0].Metadata.(string); ok {
			var parsed map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &parsed); err == nil {
				metadataMap = parsed
			}
		}
	}
	
	if metadataMap == nil {
		t.Fatalf("Failed to parse metadata: %v", events[0].Metadata)
	}
	
	if dashID, ok := metadataMap["dashboardId"]; !ok || dashID != "dash-123" {
		t.Fatalf("Metadata mismatch: %v", metadataMap)
	}
}