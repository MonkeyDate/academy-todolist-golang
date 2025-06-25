package common

import (
	"context"
	"log/slog"
)

type (
	ctxLogger  struct{}
	ctxTraceID struct{}
)

func SetLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxLogger{}, *logger)
}

func GetLogger(ctx context.Context) slog.Logger {
	logger := ctx.Value(ctxLogger{}).(slog.Logger)
	return logger
}

func SetTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, ctxTraceID{}, traceID)
}

func GetTraceID(ctx context.Context) string {
	traceID := ctx.Value(ctxTraceID{}).(string)
	return traceID
}
