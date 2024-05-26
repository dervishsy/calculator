package logger

import (
	"fmt"
	"log/slog"
)

// Logger is a wrapper around the zap logger.
type Logger struct {
	*slog.Logger
}

// NewLogger creates a new instance of the Logger.
func NewLogger(service string) (Logger, error) {
	logger := slog.Default()
	return Logger{logger}, nil
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Debug(fmt.Sprintf(format, args...))
}
