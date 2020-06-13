package log

import (
	"strings"
	"unicode"
)

type logPrinter func(format string, args ...interface{})

// LogrAdapter is a wrapper around log.Logger to implement github.com/go-pkgz/lgr.L interface
type LogrAdapter struct {
	Logger
}

func NewLogrAdapter(logger Logger) *LogrAdapter {
	return &LogrAdapter{Logger: logger}
}

func (l *LogrAdapter) Logf(format string, args ...interface{}) {
	format, logPrinterFn := l.extractLevel(format)
	logPrinterFn(format, args...)
}

func (l *LogrAdapter) extractLevel(msg string) (string, logPrinter) {
	msg = strings.TrimSpace(msg)
	const (
		levelStartChar = '['
		levelEndChar   = ']'
	)
	lvlMap := map[string]logPrinter{
		"debug":    l.Logger.Debugf,
		"info":     l.Logger.Infof,
		"warn":     l.Logger.Warnf,
		"warning":  l.Logger.Warnf,
		"err":      l.Logger.Errorf,
		"error":    l.Logger.Errorf,
		"critical": l.Logger.Criticalf,
		"fatal":    l.Logger.Fatalf,
	}
	var (
		lvlNameBuilder strings.Builder
		closePos       int
	)
loop:
	for pos, char := range msg {
		if pos == 0 && char != levelStartChar {
			break
		}
		switch char {
		case levelStartChar:
		case levelEndChar:
			closePos = pos + 1
			break loop
		default:
			lvlNameBuilder.WriteRune(unicode.ToLower(char))
		}
	}

	if lvl, ok := lvlMap[lvlNameBuilder.String()]; ok {
		return strings.TrimSpace(msg[closePos:]), lvl
	}

	return msg, l.Logger.Infof
}
