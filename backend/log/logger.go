//go:generate ../../tools/bin/go-enum --sql --marshal --nocase --names --file $GOFILE
//go:generate ../../tools/bin/mockgen -destination mock/logger_mock.go . Logger
package log

import (
	"context"
	"fmt"
	"io"
	"os"

	"go.uber.org/zap"
)

// DefaultLogger instance that can be used immediately
var DefaultLogger Logger

/*
ENUM(
debug
info
warn
error
fatal
)
*/
type Level int

/*
ENUM(
plain
json
)
*/
type Format int

func init() {
	defaultConfig := Configuration{
		LogLevel: LevelInfo,
		Format:   FormatPlain,
		Output:   os.Stderr,
	}

	logger, err := NewLogger(defaultConfig)
	if err != nil {
		panic(fmt.Sprintf("failed to initiate logger: %v", err))
	}

	DefaultLogger = logger
}

// Logger is the contract of logger
type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	With(fields ...zap.Field) Logger
	Sugared() SugaredLogger
}

// SugaredLogger is the contract of logger
type SugaredLogger interface {
	Debugf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
	Warnf(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
	Fatalf(msg string, args ...interface{})
}

type WriteSyncer interface {
	io.Writer
	Sync() error
}

// Configuration for the logger
type Configuration struct {
	LogLevel Level
	Format   Format
	Output   WriteSyncer
}

// NewLogger returns an instance of logger
func NewLogger(config Configuration) (Logger, error) {
	logger, err := newZapLogger(config)
	if err != nil {
		return nil, err
	}

	return logger, nil
}

type loggerKeyType int

const loggerKey loggerKeyType = iota

func ToContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func FromContext(ctx context.Context) Logger {
	logger, ok := ctx.Value(loggerKey).(Logger)
	if !ok {
		logger = DefaultLogger
	}
	return logger.(Logger)
}
