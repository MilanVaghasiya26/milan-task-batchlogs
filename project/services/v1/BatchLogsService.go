package v1Srv

import (
	"fmt"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/team-scaletech/common/config"
	"github.com/team-scaletech/common/helpers"
	"github.com/team-scaletech/common/logging"
	"github.com/team-scaletech/data_model/model"
	"github.com/team-scaletech/project/utils/message"

	v1Repo "github.com/team-scaletech/project/repository/v1"
	v1Req "github.com/team-scaletech/project/resources/request/v1"
)

// LogEntry represents a single log entry.
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Body      string    `json:"body"`
	Service   string    `json:"service"`
	Severity  string    `json:"severity"`
}

var (
	memoryStore []LogEntry // In-memory store for logs.
	storeMux    sync.Mutex // Mutex for thread-safe access to memoryStore.
)

const maxMemorySize = 10 * 1024 * 1024 // 10MB

// Regular expression to parse log entries
// ([\d-T:.Z]+)  -> Captures timestamp
// (\w+)         -> Captures severity (INFO, WARN, etc.)
// \[(\w+)]      -> Captures service name (apache, nginx, etc.)
// (.+)          -> Captures the full log message
var logRegex = regexp.MustCompile(`^([\d-T:.Z]+) (\w+) \[(\w+)] (.+)$`)

// IBatchLogsService defines the interface for log operations.
type IBatchLogsService interface {
	BatchLogsCreate(c *gin.Context, req v1Req.LogsEntry) error
	BatchLogsList(c *gin.Context, start, end, searchText string) error
}

// BatchLogsService handles business logic for batch logs.
type BatchLogsService struct {
	BatchLogsRepo v1Repo.IBatchLogsRepository
	Config        config.Config
}

// NewBatchLogsService initializes a new BatchLogsService instance.
func NewBatchLogsService(cf config.Config) IBatchLogsService {
	BatchLogsRepo := v1Repo.NewBatchLogsWriter()
	return &BatchLogsService{
		BatchLogsRepo: BatchLogsRepo,
		Config:        cf,
	}
}

// BatchLogsCreate processes incoming logs and persists them to the database.
func (us *BatchLogsService) BatchLogsCreate(c *gin.Context, req v1Req.LogsEntry) error {
	logger := logging.GetRequestLog(c)
	storeMux.Lock()
	defer storeMux.Unlock()

	// Parse and store logs in memory.
	for _, entry := range req.Logs {
		parsedLog := parseLog(entry)
		memoryStore = append(memoryStore, parsedLog)
	}

	// Flush to the database if memory limit is exceeded.
	if getMemorySize() > maxMemorySize {
		for _, log := range memoryStore {
			logEntry := &model.LogEntry{
				Timestamp: log.Timestamp,
				Body:      log.Body,
				Service:   log.Service,
				Severity:  log.Severity,
			}
			err := us.BatchLogsRepo.BatchLogsCreate(logEntry)
			if err != nil {
				logger.Error().Err(err).Msg("Error creating batch logs")
				return helpers.ServiceError{Message: message.SomethingWrong, Code: http.StatusInternalServerError}
			}
		}
		memoryStore = nil // Clear memory after flushing.
	}

	return nil
}

// parseLog extracts fields from a log entry string.
func parseLog(entry string) LogEntry {
	matches := logRegex.FindStringSubmatch(entry)
	if len(matches) < 5 {
		fmt.Println("Failed to parse log:", entry)
		return LogEntry{}
	}

	timestamp, severity, service, body := matches[1], matches[2], matches[3], matches[4]
	ts, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		fmt.Println("Error parsing timestamp:", err)
		return LogEntry{}
	}

	return LogEntry{Timestamp: ts, Severity: severity, Service: service, Body: body}
}

// getMemorySize estimates the memory usage of logs in memory.
func getMemorySize() int {
	return len(memoryStore) * 256 // Approximate size per log entry.
}

// BatchLogsList retrieves logs filtered by time range and search text.
func (us *BatchLogsService) BatchLogsList(c *gin.Context, start, end, searchText string) error {
	storeMux.Lock()
	defer storeMux.Unlock()

	// Return empty response if no logs in memory.
	if len(memoryStore) == 0 {
		c.JSON(http.StatusOK, []LogEntry{})
		return nil
	}

	// Parse start and end times or use defaults.
	startTime, endTime := memoryStore[0].Timestamp.UTC(), memoryStore[len(memoryStore)-1].Timestamp.UTC()
	if parsedStart, err := parseTime(start); err == nil {
		startTime = parsedStart
	}
	if parsedEnd, err := parseTime(end); err == nil {
		endTime = parsedEnd
	}

	// Filter logs based on time range and search text.
	var result []LogEntry
	for _, log := range memoryStore {
		logTime := log.Timestamp.UTC()
		if !logTime.Before(startTime) && !logTime.After(endTime) &&
			(searchText == "" || contains(log, searchText)) {
			result = append(result, log)
		}
	}

	// Send filtered logs as the response.
	c.JSON(http.StatusOK, result)
	return nil
}

// parseTime safely parses a time string in RFC3339 format.
func parseTime(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, fmt.Errorf("empty time value")
	}
	parsedTime, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time format: %w", err)
	}
	return parsedTime.UTC(), nil
}

// contains checks if a log entry contains the search text.
func contains(log LogEntry, search string) bool {
	return log.Body == search || log.Service == search || log.Severity == search
}
