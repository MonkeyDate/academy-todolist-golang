package common

import (
	"context"
	"log/slog"
)

type (
	ctxLogger  struct{}
	CtxTraceID struct{}
)

func SetLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxLogger{}, *logger)
}

func GetLogger(ctx context.Context) slog.Logger {
	logger := ctx.Value(ctxLogger{}).(slog.Logger)
	return logger
}
