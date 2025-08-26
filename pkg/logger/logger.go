package logger

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/tkcrm/mx/logger"
	"go.uber.org/zap"
)

type Logger interface {
	Debug(msg string, params ...any)
	Info(msg string, params ...any)
	Warn(msg string, params ...any)
	Error(msg string, err error, params ...any)
	Fatal(msg string, err error, params ...any)
}

type WrappedLogger struct {
	logger logger.Logger
}

type LogPrams struct {
	Slug      string
	ProcessID uuid.UUID
}

var _ Logger = (*WrappedLogger)(nil)

func New(appVersion string, conf logger.Config) Logger {
	log := logger.New(
		logger.WithLogFormat(logger.LoggerFormatJSON),
		logger.WithAppVersion(appVersion),
		logger.WithConfig(conf),
	)

	return &WrappedLogger{logger: log}
}

func (l *WrappedLogger) Debug(msg string, params ...any) {
	l.logger.Debugw(msg, params...)
}

func (l *WrappedLogger) Info(msg string, params ...any) {
	l.logger.Infow(msg, params...)
}

func (l *WrappedLogger) Warn(msg string, params ...any) {
	l.logger.Warnw(msg, params...)
}

func (l *WrappedLogger) Error(msg string, err error, params ...any) {
	l.logger.Errorw(fmt.Sprintf("%s: %s", msg, err), params...)
}

func (l *WrappedLogger) Fatal(msg string, err error, params ...any) {
	l.logger.Fatalw(msg, append(params, zap.Error(err))...)
}
