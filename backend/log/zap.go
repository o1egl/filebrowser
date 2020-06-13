package log

import (
	"sort"

	"go.uber.org/zap"         //nolint:depguard
	"go.uber.org/zap/zapcore" //nolint:depguard
)

type zapLogger struct {
	sugaredLogger *zap.SugaredLogger
}

func getEncoder(format Format) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	switch format {
	case FormatJson:
		return zapcore.NewJSONEncoder(encoderConfig)
	case FormatPlain:
		return zapcore.NewConsoleEncoder(encoderConfig)
	default:
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
}

func getZapLevel(level Level) zapcore.Level {
	switch level {
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelError:
		return zapcore.ErrorLevel
	case LevelCritical:
		return zapcore.ErrorLevel
	case LevelFatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func newZapLogger(config Configuration) (Logger, error) { //nolint:unparam
	level := getZapLevel(config.LogLevel)
	writer := zapcore.Lock(config.Output)
	core := zapcore.NewCore(getEncoder(config.Format), writer, level)

	var options []zap.Option
	if config.LogLevel == LevelDebug {
		// AddCallerSkip skips 2 number of callers, this is important else the file that gets
		// logged will always be the wrapped file. In our case zap.go
		options = append(options, zap.AddCallerSkip(2), zap.AddCaller()) //nolint:gomnd
	}
	logger := zap.New(core, options...).Sugar()

	return &zapLogger{
		sugaredLogger: logger,
	}, nil
}

func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.sugaredLogger.Debugf(format, args...)
}

func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.sugaredLogger.Infof(format, args...)
}

func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.sugaredLogger.Warnf(format, args...)
}

func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.sugaredLogger.Errorf(format, args...)
}

func (l *zapLogger) Criticalf(format string, args ...interface{}) {
	l.sugaredLogger.Errorf(format, args...)
}

func (l *zapLogger) Fatalf(format string, args ...interface{}) {
	l.sugaredLogger.Fatalf(format, args...)
}

func (l *zapLogger) WithFields(fields Fields) Logger {
	var f = make([]interface{}, 0, len(fields)*2) //nolint:gomnd
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		f = append(f, k, fields[k])
	}
	newLogger := l.sugaredLogger.With(f...)
	return &zapLogger{newLogger}
}
