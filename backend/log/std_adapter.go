package log

import (
	stdLog "log"
)

// Writer holds Logger and wraps with io.Writer interface
type Writer struct {
	Logger
	level Level
}

// Write to lgr.L
func (w *Writer) Write(p []byte) (n int, err error) {
	switch w.level {
	case LevelDebug:
		w.Debug(string(p))
	case LevelInfo:
		w.Info(string(p))
	case LevelWarn:
		w.Warn(string(p))
	case LevelError:
		w.Error(string(p))
	case LevelFatal:
		w.Fatal(string(p))
	default:
		w.Info(string(p))
	}
	return len(p), nil
}

// ToWriter makes io.Writer for given lgr.L with optional level
func ToWriter(logger Logger, level Level) *Writer {
	return &Writer{logger, level}
}

func ToStdLogger(logger Logger, level Level) *stdLog.Logger {
	return stdLog.New(ToWriter(logger, level), "", 0)
}
