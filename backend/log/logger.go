//go:generate go-enum --sql --marshal --nocase --names --file $GOFILE
//go:generate mockgen -destination mock/logger_mock.go . Logger
package log

import (
	"context"
	"fmt"
	"io"
	"os"
)

// DefaultLogger instance that can be used immediately
var DefaultLogger Logger

// Fields Type to pass when we want to call WithFields for structured logging
type Fields map[string]interface{}

/*
ENUM(
debug
info
warn
error
critical
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
	// Debugf log debug information which only be output on development env
	Debugf(format string, args ...interface{})
	// Infof log general information
	Infof(format string, args ...interface{})
	// Warnf log warning log exception that is already been handled
	Warnf(format string, args ...interface{})
	// Errorf log un-handled exception
	Errorf(format string, args ...interface{})
	// Criticalf like error, but a single instance of critical exception should trigger immediate response
	Criticalf(format string, args ...interface{})
	// Fatalf print error message which can not be recovered. Application will immediately halt
	Fatalf(format string, args ...interface{})
	// WithFields add structured KV to the log message
	WithFields(keyValues Fields) Logger
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

func NewContext(ctx context.Context, fields Fields) context.Context {
	return context.WithValue(ctx, loggerKey, WithContext(ctx).WithFields(fields))
}

func WithContext(ctx context.Context) Logger {
	if ctx == nil {
		return DefaultLogger
	}
	if logger, ok := ctx.Value(loggerKey).(Logger); ok {
		return logger
	}
	return DefaultLogger
}

// Debugf log debug information
func Debugf(format string, args ...interface{}) {
	DefaultLogger.Debugf(format, args...)
}

// Infof log general information
func Infof(format string, args ...interface{}) {
	DefaultLogger.Infof(format, args...)
}

// Warnf log warning log exception that is already been handled
func Warnf(format string, args ...interface{}) {
	DefaultLogger.Warnf(format, args...)
}

// Errorf log un-handled exception
func Errorf(format string, args ...interface{}) {
	DefaultLogger.Errorf(format, args...)
}

// Criticalf like error, but a single instance of critical exception should trigger immediate response
func Criticalf(format string, args ...interface{}) {
	DefaultLogger.Criticalf(format, args...)
}

// Fatalf print error message which can not be recovered. Application will immediately halt
func Fatalf(format string, args ...interface{}) {
	DefaultLogger.Fatalf(format, args...)
}

// WithFields add structured KV to the log message
func WithFields(keyValues Fields) Logger {
	return DefaultLogger.WithFields(keyValues)
}
