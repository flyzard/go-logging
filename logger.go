// Package logging provides a simple logging interface for the application.
package logging

import (
	"os"

	"github.com/phuslu/log"
)

// LogLevel defines the severity of the log message.
type LogLevel int

// Log levels.
const (
	LogLevelInfo LogLevel = iota
	LogLevelWarning
	LogLevelError
)

// LoggerInterface is the interface for the application's logging.
type LoggerInterface interface {
	Info(format string, v ...any)
	Warning(format string, v ...any)
	Error(format string, v ...any)
	SetLogLevel(level LogLevel)
}

// Logger is the application's logging interface.
type Logger struct {
	logger   *log.Logger
	logLevel LogLevel
}

// NewLogger creates a new Logger instance.
func NewLogger(logLevel LogLevel) *Logger {
	l := log.Logger{
		Writer: &log.ConsoleWriter{
			Writer:         os.Stdout,
			ColorOutput:    true,
			QuoteString:    true,
			EndWithMessage: true,
		},
		TimeFormat: "2006-01-02 15:04:05",
	}
	return &Logger{
		logger:   &l,
		logLevel: logLevel,
	}
}

// Info logs informational messages.
func (l *Logger) Info(format string, v ...any) {
	if l.logLevel <= LogLevelInfo {
		l.logger.Info().Msgf(format, v...)
	}
}

// Warning logs warning messages.
func (l *Logger) Warning(format string, v ...any) {
	if l.logLevel <= LogLevelWarning {
		l.logger.Warn().Msgf(format, v...)
	}
}

// Error logs error messages.
func (l *Logger) Error(format string, v ...any) {
	if l.logLevel <= LogLevelError {
		l.logger.Error().Msgf(format, v...)
	}
}

// SetLogLevel changes the current log level of the logger.
func (l *Logger) SetLogLevel(level LogLevel) {
	l.logLevel = level
}
