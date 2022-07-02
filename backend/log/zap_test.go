package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap" //nolint:depguard
)

func Test_getZapLevel(t *testing.T) {
	assert.Equal(t, zap.DebugLevel, getZapLevel(LevelDebug))
	assert.Equal(t, zap.InfoLevel, getZapLevel(LevelInfo))
	assert.Equal(t, zap.WarnLevel, getZapLevel(LevelWarn))
	assert.Equal(t, zap.ErrorLevel, getZapLevel(LevelError))
	assert.Equal(t, zap.FatalLevel, getZapLevel(LevelFatal))
	assert.Equal(t, zap.InfoLevel, getZapLevel(Level(-999)))
}
