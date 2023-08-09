package logger

import (
	"context"
	"golang.org/x/exp/slog"
)

type ContextHandler struct {
	slog.Handler
}

func NewContextHandler(handler slog.Handler) *ContextHandler {
	return &ContextHandler{Handler: handler}
}

func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if ctxAttrs, ok := ctx.Value(contextKey).([]slog.Attr); ok {
		r.AddAttrs(ctxAttrs...)
	}
	return h.Handler.Handle(ctx, r)
}

type contextKeyType string

const contextKey contextKeyType = "logger-key"

// ContextWithAttrs returns a new context with the given log attributes.
func ContextWithAttrs(ctx context.Context, attrs ...slog.Attr) context.Context {
	ctxAttrs, ok := ctx.Value(contextKey).([]slog.Attr)
	if !ok {
		return context.WithValue(ctx, contextKey, attrs)
	}
	return context.WithValue(ctx, contextKey, append(ctxAttrs, attrs...))
}
