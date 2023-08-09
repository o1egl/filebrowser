package logger

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
	"regexp"
	"strings"
	"testing"
)

func TestWithContextAttrs(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewContextHandler(slog.NewTextHandler(buf, nil))
	logger := slog.New(handler)

	ctx := context.Background()

	t.Run("no attrs", func(t *testing.T) {
		buf.Reset()
		logger.InfoContext(ctx, "hello")

		assert.Equal(t, stripPrefix(buf.String()), "msg=hello")
	})

	t.Run("add key attr", func(t *testing.T) {
		buf.Reset()
		ctx := ContextWithAttrs(ctx, slog.String("key", "value"))
		logger.InfoContext(ctx, "hello")

		assert.Equal(t, stripPrefix(buf.String()), "msg=hello key=value")

		t.Run("add another key", func(t *testing.T) {
			buf.Reset()
			ctx := ContextWithAttrs(ctx, slog.String("key2", "value2"))
			logger.InfoContext(ctx, "hello")

			assert.Equal(t, stripPrefix(buf.String()), "msg=hello key=value key2=value2")
		})
	})
}

func stripPrefix(s string) string {
	timeRegexp := regexp.MustCompile(`time=\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?(Z|[+-]\d{2}:\d{2}) level=INFO`)
	return strings.TrimSpace(timeRegexp.ReplaceAllString(s, ""))
}
