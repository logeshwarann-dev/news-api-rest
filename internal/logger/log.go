package logger

import (
	"context"
	"log/slog"
	"os"
)

type ctxKey struct{}

func CtxWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	if logger == nil {
		return ctx
	}

	if ctxLog, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok && ctxLog == logger {
		return ctx
	}
	return context.WithValue(ctx, ctxKey{}, logger)
}

func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok {
		return logger
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
}
