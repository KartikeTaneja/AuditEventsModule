// audit/db_audit.go
package audit

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// DBAuditLogger implements AuditLogger interface using database storage
type DBAuditLogger struct {
	db *sql.DB
}

// NewDBAuditLogger creates a new DBAuditLogger
func NewDBAuditLogger(dbPath string) (*DBAuditLogger, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Create table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS audit_events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		action_string TEXT NOT NULL,
		extra_msg TEXT,
		epoch_timestamp_sec INTEGER NOT NULL,
		org_id INTEGER NOT NULL,
		metadata TEXT
	);
	CREATE INDEX IF NOT EXISTS idx_org_time ON audit_events(org_id, epoch_timestamp_sec);
	`
	
	_, err = db.Exec(createTableSQL)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create audit events table: %v", err)
	}

	return &DBAuditLogger{
		db: db,
	}, nil
}

// Close closes the database connection
func (l *DBAuditLogger) Close() error {
	return l.db.Close()
}

// CreateAuditEvent logs a user action to the database
func (l *DBAuditLogger) CreateAuditEvent(username, actionString, extraMsg string, epochTimestampSec, orgID int64, metadata interface{}) error {
	// If timestamp is 0, use current time
	if epochTimestampSec == 0 {
		epochTimestampSec = time.Now().Unix()
	}

	var metadataJSON []byte
	var err error

	if metadata != nil {
		metadataJSON, err = json.Marshal(metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %v", err)
		}
	}

	insertSQL := `
	INSERT INTO audit_events (username, action_string, extra_msg, epoch_timestamp_sec, org_id, metadata)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err = l.db.Exec(insertSQL, username, actionString, extraMsg, epochTimestampSec, orgID, string(metadataJSON))
	if err != nil {
		return fmt.Errorf("failed to insert audit event: %v", err)
	}

	return nil
}

// ReadAuditEvents reads audit events from the database for a specific organization and time range
func (l *DBAuditLogger) ReadAuditEvents(orgID int64, startEpochSec, endEpochSec int64) ([]AuditEvent, error) {
	query := `
	SELECT username, action_string, extra_msg, epoch_timestamp_sec, org_id, metadata
	FROM audit_events
	WHERE org_id = ? AND epoch_timestamp_sec >= ?
	`
	args := []interface{}{orgID, startEpochSec}

	if endEpochSec > 0 {
		query += " AND epoch_timestamp_sec <= ?"
		args = append(args, endEpochSec)
	}

	query += " ORDER BY epoch_timestamp_sec ASC"

	rows, err := l.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query audit events: %v", err)
	}
	defer rows.Close()

	var events []AuditEvent
	for rows.Next() {
		var event AuditEvent
		var metadataStr sql.NullString

		err := rows.Scan(
			&event.Username,
			&event.ActionString,
			&event.ExtraMsg,
			&event.EpochTimestampSec,
			&event.OrgID,
			&metadataStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan audit event row: %v", err)
		}

		if metadataStr.Valid && metadataStr.String != "" {
			var metadata interface{}
			if err := json.Unmarshal([]byte(metadataStr.String), &metadata); err == nil {
				event.Metadata = metadata
			}
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating audit event rows: %v", err)
	}

	return events, nil
}