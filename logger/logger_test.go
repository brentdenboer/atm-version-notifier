package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	// Create a buffer to capture log output
	var buf bytes.Buffer

	// Create a new logger instance
	logger := New(&buf)

	tests := []struct {
		name        string
		level       Level
		logFunc     func(string, ...interface{})
		message     string
		shouldLog   bool
		levelString string
	}{
		{
			name:        "debug message at debug level",
			level:       DEBUG,
			logFunc:     logger.Debug,
			message:     "debug message",
			shouldLog:   true,
			levelString: "DEBUG",
		},
		{
			name:        "info message at debug level",
			level:       DEBUG,
			logFunc:     logger.Info,
			message:     "info message",
			shouldLog:   true,
			levelString: "INFO",
		},
		{
			name:        "debug message at info level",
			level:       INFO,
			logFunc:     logger.Debug,
			message:     "debug message",
			shouldLog:   false,
			levelString: "DEBUG",
		},
		{
			name:        "warning message at info level",
			level:       INFO,
			logFunc:     logger.Warn,
			message:     "warning message",
			shouldLog:   true,
			levelString: "WARN",
		},
		{
			name:        "error message at warning level",
			level:       WARNING,
			logFunc:     logger.Error,
			message:     "error message",
			shouldLog:   true,
			levelString: "ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset buffer and set log level
			buf.Reset()
			logger.SetLevel(tt.level)

			// Log the message
			tt.logFunc(tt.message)

			// Check if message was logged
			output := buf.String()
			if tt.shouldLog {
				if !strings.Contains(output, tt.message) {
					t.Errorf("Expected message '%s' to be logged, but it wasn't. Output: %s", tt.message, output)
				}
				if !strings.Contains(output, tt.levelString) {
					t.Errorf("Expected level '%s' in log output, but it wasn't found. Output: %s", tt.levelString, output)
				}
				if !strings.Contains(output, "logger_test.go") {
					t.Errorf("Expected source file info in log output, but it wasn't found. Output: %s", output)
				}
			} else {
				if output != "" {
					t.Errorf("Expected no output, but got: %s", output)
				}
			}
		})
	}
}

func TestMultipleWriters(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	logger := New(&buf1, &buf2)
	logger.SetLevel(INFO)

	testMessage := "test multiple writers"
	logger.Info(testMessage)

	for i, buf := range []*bytes.Buffer{&buf1, &buf2} {
		output := buf.String()
		if !strings.Contains(output, testMessage) {
			t.Errorf("Writer %d: Expected message '%s' to be logged, but it wasn't. Output: %s", i+1, testMessage, output)
		}
	}
}

func TestDefaultLogger(t *testing.T) {
	// Create a buffer to capture output
	var buf bytes.Buffer
	SetOutput(&buf)

	// Test default logger functions
	testMessage := "test default logger"
	Info(testMessage)

	// Check the output
	output := buf.String()
	if !strings.Contains(output, testMessage) {
		t.Errorf("Expected message '%s' to be logged, but it wasn't. Output: %s", testMessage, output)
	}
	if !strings.Contains(output, "INFO") {
		t.Errorf("Expected 'INFO' level in log output, but it wasn't found. Output: %s", output)
	}
	if !strings.Contains(output, "logger_test.go") {
		t.Errorf("Expected source file info in log output, but it wasn't found. Output: %s", output)
	}
}

func TestSetOutput(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)

	testMessage := "test set output"
	Info(testMessage)

	output := buf.String()
	if !strings.Contains(output, testMessage) {
		t.Errorf("Expected message '%s' to be logged, but it wasn't. Output: %s", testMessage, output)
	}
}

func TestFatal(t *testing.T) {
	// Create a buffer to capture log output
	var buf bytes.Buffer
	logger := New(&buf)
	logger.SetTesting(true) // Enable testing mode

	// Use a deferred function to catch the panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected Fatal to panic in testing mode, but it didn't")
		} else if r != "fatal error in testing mode" {
			t.Errorf("Expected panic with 'fatal error in testing mode', got %v", r)
		}
	}()

	testMessage := "fatal error"
	logger.Fatal(testMessage)

	// We shouldn't reach this point, but if we do, check the output
	output := buf.String()
	if !strings.Contains(output, testMessage) {
		t.Errorf("Expected message '%s' to be logged, but it wasn't. Output: %s", testMessage, output)
	}
	if !strings.Contains(output, "FATAL") {
		t.Errorf("Expected 'FATAL' level in log output, but it wasn't found. Output: %s", output)
	}
}
