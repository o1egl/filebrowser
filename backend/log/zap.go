package log

import (
	"go.uber.org/zap"         //nolint:depguard
	"go.uber.org/zap/zapcore" //nolint:depguard
)

type zapLogger struct {
	logger *zap.Logger
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

	return &zapLogger{
		logger: zap.New(core, options...),
	}, nil
}

func (z *zapLogger) Debug(msg string, fields ...zap.Field) {
	z.logger.Debug(msg, fields...)
}

func (z *zapLogger) Info(msg string, fields ...zap.Field) {
	z.logger.Info(msg, fields...)
}

func (z *zapLogger) Warn(msg string, fields ...zap.Field) {
	z.logger.Warn(msg, fields...)
}

func (z *zapLogger) Error(msg string, fields ...zap.Field) {
	z.logger.Error(msg, fields...)
}

func (z *zapLogger) Fatal(msg string, fields ...zap.Field) {
	z.logger.Fatal(msg, fields...)
}

func (z *zapLogger) With(fields ...zap.Field) Logger {
	return &zapLogger{
		logger: z.logger.With(fields...),
	}
}

func (z *zapLogger) Sugared() SugaredLogger {
	return z.logger.Sugar()
}
