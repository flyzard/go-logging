package logging

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/phuslu/log"
)

// logEntry represents the structure of our JSON log output
type logEntry struct {
	Level   string `json:"level"`
	Time    string `json:"time"`
	Message string `json:"message"`
}

// testWriter creates a logger with a buffer capture
func testLogger(level LogLevel) (*Logger, *bytes.Buffer) {
	buf := new(bytes.Buffer)
	l := NewLogger(level)
	l.logger = &log.Logger{
		Writer: &log.IOWriter{Writer: buf},
	}
	return l, buf
}

func parseLogEntry(buf *bytes.Buffer) (logEntry, error) {
	var entry logEntry
	err := json.Unmarshal(buf.Bytes(), &entry)
	return entry, err
}

func TestNewLogger(t *testing.T) {
	t.Run("default log level", func(t *testing.T) {
		logger, _ := testLogger(LogLevelInfo)
		if logger.logLevel != LogLevelInfo {
			t.Errorf("Expected default log level Info, got %v", logger.logLevel)
		}
	})
}

func TestLogLevels(t *testing.T) {
	testCases := []struct {
		name      string
		setLevel  LogLevel
		shouldLog map[string]bool // method name -> should log
	}{
		{
			name:     "Info level",
			setLevel: LogLevelInfo,
			shouldLog: map[string]bool{
				"Info":    true,
				"Warning": true,
				"Error":   true,
			},
		},
		{
			name:     "Warning level",
			setLevel: LogLevelWarning,
			shouldLog: map[string]bool{
				"Info":    false,
				"Warning": true,
				"Error":   true,
			},
		},
		{
			name:     "Error level",
			setLevel: LogLevelError,
			shouldLog: map[string]bool{
				"Info":    false,
				"Warning": false,
				"Error":   true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logger, buf := testLogger(tc.setLevel)

			tests := []struct {
				name     string
				logFunc  func()
				expected bool
			}{
				{"Info", func() { logger.Info("test") }, tc.shouldLog["Info"]},
				{"Warning", func() { logger.Warning("test") }, tc.shouldLog["Warning"]},
				{"Error", func() { logger.Error("test") }, tc.shouldLog["Error"]},
			}

			for _, test := range tests {
				buf.Reset()
				test.logFunc()

				if test.expected && buf.Len() == 0 {
					t.Errorf("%s should have logged but didn't", test.name)
				}
				if !test.expected && buf.Len() > 0 {
					t.Errorf("%s shouldn't have logged but did", test.name)
				}
			}
		})
	}
}

func TestLogOutput(t *testing.T) {
	logger, buf := testLogger(LogLevelInfo)
	expectedMsg := "test message 42"

	logger.Info("test message %d", 42)
	entry, err := parseLogEntry(buf)
	if err != nil {
		t.Fatalf("Failed to parse log entry: %v", err)
	}

	if entry.Message != expectedMsg {
		t.Errorf("Expected message '%s', got '%s'", expectedMsg, entry.Message)
	}

	if entry.Level != "info" {
		t.Errorf("Expected level 'info', got '%s'", entry.Level)
	}

	if _, err := time.Parse(time.RFC3339, entry.Time); err != nil {
		t.Errorf("Invalid timestamp format: %v", err)
	}
}

func TestSetLogLevel(t *testing.T) {
	logger, buf := testLogger(LogLevelError)

	// Verify initial level
	logger.Info("test")
	if buf.Len() > 0 {
		t.Error("Info should not log at Error level")
	}

	// Change level and verify new behavior
	logger.SetLogLevel(LogLevelInfo)
	buf.Reset()
	logger.Info("test")
	if buf.Len() == 0 {
		t.Error("Info should log after level change")
	}
}

func TestErrorLogging(t *testing.T) {
	logger, buf := testLogger(LogLevelError)
	err := errors.New("critical error")

	logger.Error("operation failed: %v", err)
	entry, _ := parseLogEntry(buf)

	expectedMsg := "operation failed: critical error"
	if entry.Message != expectedMsg {
		t.Errorf("Expected '%s', got '%s'", expectedMsg, entry.Message)
	}

	if entry.Level != "error" {
		t.Errorf("Expected level 'error', got '%s'", entry.Level)
	}
}

func TestLogFormatting(t *testing.T) {
	logger, buf := testLogger(LogLevelInfo)

	testCases := []struct {
		name     string
		logFunc  func()
		expected string
	}{
		{
			"Complex formatting",
			func() { logger.Info("User %s (%d) logged in", "john", 123) },
			"User john (123) logged in",
		},
		{
			"Special characters",
			func() { logger.Warning("Path: %q", "/etc/passwd") },
			`Path: "/etc/passwd"`,
		},
		{
			"Error formatting",
			func() { logger.Error("Failed: %v", errors.New("timeout")) },
			"Failed: timeout",
		},
		// Add more test cases for different formatting scenarios
		{
			"Multiple placeholders",
			func() { logger.Info("%s: %.2f%% complete", "Download", 75.5) },
			"Download: 75.50% complete",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf.Reset()
			tc.logFunc()

			entry, err := parseLogEntry(buf)
			if err != nil {
				t.Fatalf("Failed to parse log entry: %v", err)
			}

			if entry.Message != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, entry.Message)
			}
		})
	}
}
