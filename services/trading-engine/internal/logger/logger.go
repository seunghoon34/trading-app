// file: trading-engine/internal/logger/logger.go
package logger

import (
	"github.com/sirupsen/logrus"
)

type CustomFormatter struct {
	ServiceName string
}

type Logger struct {
	*logrus.Logger
}

var defaultLogger *Logger

func init() {
	// Initialize default logger
	defaultLogger = New()
}

func New() *Logger {
	log := logrus.New()

	// Use JSON formatter for ELK
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z",
	})
	log.SetLevel(logrus.InfoLevel)

	return &Logger{Logger: log}
}

// Global functions that use the default logger
func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}

func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
}

// WithFields for structured logging
func WithFields(fields map[string]interface{}) *logrus.Entry {
	// Always include service name
	fields["service"] = "trading-engine"
	return defaultLogger.WithFields(fields)
}

// GetLogger returns the default logger instance
func GetLogger() *Logger {
	return defaultLogger
}
