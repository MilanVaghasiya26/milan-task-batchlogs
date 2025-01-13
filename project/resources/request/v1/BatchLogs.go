package v1Req

// LogsEntry represents a batch of log messages to be processed.
type LogsEntry struct {
	Logs []string `json:"logs"` // Logs is a slice of log entries in string format.
}
