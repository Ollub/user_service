package log

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
)

const (
	WarnLevel = iota
	InfoLevel
	DebugLevel
)

type Fields map[string]interface{}

type ILogger interface {
	Debug(string, ...Fields)
	Info(string, ...Fields)
	Warn(string, ...Fields)
	Error(string, ...Fields)
	SetContext(Fields)
}

var std = logrus.New()
var defaultLogger = New()

func Debug(msg string, fields ...Fields) {
	defaultLogger.Debug(msg, fields...)
}
func Info(msg string, fields ...Fields) {
	defaultLogger.Info(msg, fields...)
}
func Warn(msg string, fields ...Fields) {
	defaultLogger.Warn(msg, fields...)
}
func Error(msg string, fields ...Fields) {
	defaultLogger.Error(msg, fields...)
}

type Logger struct {
	logger *logrus.Entry
}

func New(fields ...Fields) *Logger {
	var logger *logrus.Entry
	if len(fields) != 0 {
		logger = std.WithFields(logrus.Fields(fields[0]))
	} else {
		logger = logrus.NewEntry(std)
	}
	return &Logger{logger: logger}
}

func (log *Logger) SetContext(fields Fields) {
	log.logger = log.logger.WithFields(logrus.Fields(fields))
}

func (log *Logger) Debug(msg string, fields ...Fields) {
	log.withFields(fields...).Debug(msg)
}

func (log *Logger) Info(msg string, fields ...Fields) {
	log.withFields(fields...).Info(msg)
}

func (log *Logger) Warn(msg string, fields ...Fields) {
	log.withFields(fields...).Warn(msg)
}

func (log *Logger) Error(msg string, fields ...Fields) {
	log.withFields(fields...).Error(msg)
}

func (log *Logger) withFields(fields ...Fields) *logrus.Entry {
	if len(fields) == 0 {
		return log.logger
	}
	return log.logger.WithFields(logrus.Fields(fields[0]))
}

func SetupLogger(level int) {
	switch level {
	case WarnLevel:
		std.SetLevel(logrus.WarnLevel)
	case InfoLevel:
		std.SetLevel(logrus.InfoLevel)
	case DebugLevel:
		std.SetLevel(logrus.DebugLevel)
	}
	std.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	Debug("Logger configured", Fields{"level": level})
}

func Rlog(r *http.Request) ILogger {
	return fromContext(r.Context())
}

func Clog(ctx context.Context) ILogger {
	return fromContext(ctx)
}

func fromContext(ctx context.Context) ILogger {
	logger, ok := ctx.Value(LoggerKey).(ILogger)
	if !ok || logger == nil {
		return defaultLogger
	}
	return logger
}
