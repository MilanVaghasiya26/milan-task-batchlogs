package logging

import (
	"bytes"
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/team-scaletech/common/config"
)

const (
	RequestIDKey    = "X-Request-ID"
	LogRequestIDKey = "req-id"
)

// Data is
type Data struct {
	IPAddress string // the ip address of the caller
}

// ILogger is
type ILogger interface {
	Debug() *zerolog.Event
	Info() *zerolog.Event
	Warn() *zerolog.Event
	Error() *zerolog.Event
	Fatal() *zerolog.Event
	Panic() *zerolog.Event
}

// Logger Impl is
type Logger struct {
	appName    string
	appVersion string

	filePath string
	level    zerolog.Level
	maxAge   time.Duration

	theLogger zerolog.Logger
}

var defaultLogger Logger
var defaultLoggerOnce sync.Once
var logLevel = zerolog.DebugLevel

func NewLogger(cf config.Config) {
	// Set log level
	level, err := zerolog.ParseLevel(cf.Level)
	if err != nil {
		// Default to debug if the log level is invalid
		level = zerolog.DebugLevel
	}
	logLevel = level

	createLogger()
}

func createLogger() {
	zerolog.TimeFieldFormat = time.RFC3339Nano

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(logLevel)

	filename := "logs.log"
	logger := &lumberjack.Logger{
		Filename:   filepath.Join(defaultLogger.filePath, filename), // Log file location
		MaxSize:    10,                                              // Max size in MB before a new file is created
		MaxBackups: 3,                                               // Max number of old log files to keep
		MaxAge:     28,                                              // Max age in days to keep an old log file
		Compress:   true,                                            // Compress old log files (gzip)
	}

	fileLogger := zerolog.New(logger).Hook(&myCustomHook{writer: logger}) // Set file output

	// TODO: Once fix the log insertion issue to log file remove the following code.
	hook := &myCustomHook{writer: logger}

	// Define a custom hook for rotating log files
	defaultLogger.theLogger.Hook(hook)
	// Use multi-level logger to log to both console and file
	defaultLogger.theLogger = zerolog.New(zerolog.MultiLevelWriter(fileLogger, zerolog.ConsoleWriter{Out: os.Stdout})).With().Timestamp().Logger()
}

// Define the myCustomHook type
type myCustomHook struct {
	writer io.Writer
}

// Run Implement the zerolog.Hook interface
func (h *myCustomHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	//Not required now
	//fmt.Fprintf(h.writer, "%s\n", msg)
}

// GetLog is
func GetLog() ILogger {
	defaultLoggerOnce.Do(createLogger)
	return &defaultLogger
}

func GetRequestLog(c *gin.Context) ILogger {
	defaultLoggerOnce.Do(createLogger)

	l := defaultLogger

	if c != nil {
		reqID, ok := c.Get(RequestIDKey)
		if ok {
			l.theLogger = l.theLogger.With().Str(LogRequestIDKey, reqID.(string)).Logger()
		}
	}
	return &l
}

func Middleware(c *gin.Context) {
	start := time.Now()
	c.Next()
	status := c.Writer.Status()
	end := time.Now()
	latency := end.Sub(start)
	ctxRqId, ok := c.Value(RequestIDKey).(string)
	if ok {
		defaultLogger.theLogger.Info().
			Str(LogRequestIDKey, ctxRqId).
			Str("uri", c.Request.RequestURI).
			Str("method", c.Request.Method).
			Str("ip", c.ClientIP()).
			Dur("latency", latency).
			Int("status", status).
			Msg("")
	} else {
		defaultLogger.theLogger.Info().
			Str("uri", c.Request.RequestURI).
			Str("method", c.Request.Method).
			Str("ip", c.ClientIP()).
			Dur("latency", latency).
			Int("status", status).
			Msg("")
	}
}

func (logger *Logger) getLogEntry() *zerolog.Logger {
	pc, _, _, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	_, line := runtime.FuncForPC(pc).FileLine(pc)
	var buffer bytes.Buffer

	buffer.WriteString("fn:")

	x := strings.LastIndex(funcName, "/")
	buffer.WriteString(funcName[x+1:] + " L:" + fmt.Sprint(line))

	var zLog = logger.theLogger.With().Str("info", buffer.String()).Logger()
	return &zLog
}

// Debug returns a new debug event
func (logger *Logger) Debug() *zerolog.Event {
	return logger.getLogEntry().Debug()
}

// Info returns a new info event
func (logger *Logger) Info() *zerolog.Event {
	return logger.getLogEntry().Info()
}

// Warn returns a new warn event
func (logger *Logger) Warn() *zerolog.Event {
	return logger.getLogEntry().Warn()
}

// Error returns a new error event
func (logger *Logger) Error() *zerolog.Event {
	return logger.getLogEntry().Error()
}

// Fatal returns a new fatal event
func (logger *Logger) Fatal() *zerolog.Event {
	return logger.getLogEntry().Fatal()
}

// Panic returns a new panic event
func (logger *Logger) Panic() *zerolog.Event {
	return logger.getLogEntry().Panic()
}

func formatLatency(milliseconds int64) string {
	duration := time.Duration(milliseconds) * time.Millisecond
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	millisecondsLeft := duration.Milliseconds() % 1000

	formatted := ""
	if hours > 0 {
		formatted += fmt.Sprintf("%dh ", hours)
	}
	if minutes > 0 {
		formatted += fmt.Sprintf("%dm ", minutes)
	}
	if seconds > 0 {
		formatted += fmt.Sprintf("%ds ", seconds)
	}
	if millisecondsLeft > 0 {
		formatted += fmt.Sprintf("%dms", millisecondsLeft)
	}

	return formatted
}
