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
		w.Debugf(string(p))
	case LevelInfo:
		w.Infof(string(p))
	case LevelWarn:
		w.Warnf(string(p))
	case LevelError:
		w.Errorf(string(p))
	case LevelCritical:
		w.Criticalf(string(p))
	case LevelFatal:
		w.Fatalf(string(p))
	default:
		w.Infof(string(p))
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
