package types

import (
	pkglogger "github.com/MetaDiv-AI/logger"

	"go.uber.org/zap"
)

func NewLogger(app string, handlerName string) *logger {
	log := pkglogger.New().Build()
	return &logger{
		log:         log.With(zap.String("app", app), zap.String("handler", handlerName)),
		app:         app,
		handlerName: handlerName,
	}
}

type Logger interface {
	Logger() *logger
	// Info logs an info message with structured fields
	LogInfo(message string, fields ...zap.Field)
	// Error logs an error message with structured fields
	LogError(message string, fields ...zap.Field)
	// Debug logs a debug message with structured fields
	LogDebug(message string, fields ...zap.Field)
	// Warn logs a warning message with structured fields
	LogWarn(message string, fields ...zap.Field)
}

type logger struct {
	log         pkglogger.Logger
	app         string
	handlerName string
}

func (l *logger) Logger() *logger {
	return l
}

func (l *logger) LogInfo(message string, fields ...zap.Field) {
	l.log.Info(message, fields...)
}

func (l *logger) LogError(message string, fields ...zap.Field) {
	l.log.Error(message, fields...)
}

func (l *logger) LogDebug(message string, fields ...zap.Field) {
	l.log.Debug(message, fields...)
}

func (l *logger) LogWarn(message string, fields ...zap.Field) {
	l.log.Warn(message, fields...)
}
