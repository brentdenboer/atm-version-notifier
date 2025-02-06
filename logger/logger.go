package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// Level represents the logging level
type Level int

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

var levelStrings = map[Level]string{
	DEBUG:   "DEBUG",
	INFO:    "INFO",
	WARNING: "WARN",
	ERROR:   "ERROR",
	FATAL:   "FATAL",
}

// Logger represents our custom logger
type Logger struct {
	logger  *log.Logger
	level   Level
	testing bool // Used to control exit behavior in tests
}

// New creates a new Logger instance
func New(writers ...io.Writer) *Logger {
	// If no writers provided, default to stdout
	if len(writers) == 0 {
		writers = []io.Writer{os.Stdout}
	}

	// Create multi-writer if multiple writers provided
	var writer io.Writer
	if len(writers) == 1 {
		writer = writers[0]
	} else {
		writer = io.MultiWriter(writers...)
	}

	return &Logger{
		logger:  log.New(writer, "", 0),
		level:   INFO, // Default level
		testing: false,
	}
}

// SetLevel sets the minimum logging level
func (l *Logger) SetLevel(level Level) {
	l.level = level
}

// SetTesting sets the testing mode for the logger
func (l *Logger) SetTesting(testing bool) {
	l.testing = testing
}

// formatMessage formats the log message with timestamp, level, and caller info
func (l *Logger) formatMessage(level Level, format string, args ...interface{}) string {
	// Get caller information
	_, file, line, ok := runtime.Caller(3)
	callerInfo := "???"
	if ok {
		callerInfo = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}

	// Format the message
	msg := fmt.Sprintf(format, args...)
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	return fmt.Sprintf("%s [%s] %s: %s", timestamp, levelStrings[level], callerInfo, msg)
}

// log performs the actual logging
func (l *Logger) log(level Level, format string, args ...interface{}) {
	if level >= l.level {
		l.logger.Output(3, l.formatMessage(level, format, args...))
	}
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARNING, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
	if l.testing {
		panic("fatal error in testing mode")
	}
	os.Exit(1)
}

// Default logger instance
var defaultLogger = New(os.Stdout)

// Package-level functions that use the default logger
func Debug(format string, args ...interface{}) { defaultLogger.log(DEBUG, format, args...) }
func Info(format string, args ...interface{})  { defaultLogger.log(INFO, format, args...) }
func Warn(format string, args ...interface{})  { defaultLogger.log(WARNING, format, args...) }
func Error(format string, args ...interface{}) { defaultLogger.log(ERROR, format, args...) }
func Fatal(format string, args ...interface{}) { defaultLogger.log(FATAL, format, args...) }

// SetLevel sets the level for the default logger
func SetLevel(level Level) { defaultLogger.SetLevel(level) }

// SetOutput sets the output for the default logger
func SetOutput(w io.Writer) { defaultLogger = New(w) }

// SetTesting sets the testing mode for the default logger
func SetTesting(testing bool) { defaultLogger.SetTesting(testing) }
