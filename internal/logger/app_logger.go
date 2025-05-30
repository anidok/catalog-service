package logger

import (
	"catalog-service/internal/appcontext"
	"context"
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

const (
	LogMethod = "Method"
	LogError  = "Error"
)

func Setup(logLevel, logFormat string) {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		log.Fatalf("fatal error: %v", err)
	}

	logger = &logrus.Logger{
		Out:       os.Stdout,
		Hooks:     make(logrus.LevelHooks),
		Level:     level,
		Formatter: &logrus.JSONFormatter{},
	}

	if logFormat != "json" {
		logger.Formatter = &logrus.TextFormatter{}
	}
	NonContext = &ContextLogger{
		entry: logrus.NewEntry(logger),
	}
}

type ContextLogger struct {
	entry *logrus.Entry
}

var NonContext *ContextLogger

func NewContextLogger(ctx context.Context, method string) *ContextLogger {
	logEntry := logger.WithFields(appcontext.LogFields(ctx)).WithField(LogMethod, method)
	return &ContextLogger{
		entry: logEntry,
	}
}

func (l *ContextLogger) Fatal(errMessage string, err error) {
	l.entry.
		WithField(LogError, err).
		Fatal(errMessage)
}

func (l *ContextLogger) Error(errMessage string, err error) {
	l.entry.WithField(LogError, err).Error(errMessage)
}

func (l *ContextLogger) Errorf(err error, errMessageFormat string, errMessages ...interface{}) {
	l.entry.WithField(LogError, err).Errorf(errMessageFormat, errMessages...)
}

func (l *ContextLogger) ErrorWithFields(msg string, fields map[string]interface{}, err error) {
	for key, val := range fields {
		l.entry = l.entry.WithField(key, val)
	}

	l.entry.WithField(LogError, err).Error(msg)
}

func (l *ContextLogger) Info(msg string) {
	l.entry.Info(msg)

}

func (l *ContextLogger) Infof(msg string, args ...interface{}) {
	l.entry.Infof(msg, args...)
}

func (l *ContextLogger) Debugf(msg string, args ...interface{}) {
	l.entry.Debugf(msg, args...)
}

func (l *ContextLogger) Debug(msg string) {
	l.entry.Debug(msg)
}

func (l *ContextLogger) InfoWithFields(msg string, fields map[string]interface{}) {
	for key, val := range fields {
		l.entry = l.entry.WithField(key, val)
	}

	l.entry.Info(msg)
}

func (l *ContextLogger) DebugWithFields(msg string, fields map[string]interface{}) {
	for key, val := range fields {
		l.entry = l.entry.WithField(key, val)
	}

	l.entry.Debug(msg)
}

func (l *ContextLogger) Warn(msg string) {
	l.entry.Warn(msg)
}

func (l *ContextLogger) Warnf(msg string, args ...interface{}) {
	l.entry.Warnf(msg, args...)
}

func (l *ContextLogger) WarnWithFields(msg string, fields map[string]interface{}, err error) {
	for key, val := range fields {
		l.entry = l.entry.WithField(key, val)
	}

	l.entry.WithField(LogError, err)
	l.entry.Warn(msg)
}
