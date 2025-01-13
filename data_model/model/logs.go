package model

import "time"

// LogEntry represents a log stored in the database.
type LogEntry struct {
	BaseEntity
	Timestamp time.Time `json:"timestamp"`
	Body      string    `json:"body"`
	Service   string    `json:"service"`
	Severity  string    `json:"severity"`
}

// TableName specifies the database table name for the User model.
func (le *LogEntry) TableName() string {
	return "logs"
}
