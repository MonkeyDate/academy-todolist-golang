package main

import (
	"academy-todo/internal/common"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

func TraceIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			traceID := r.Header.Get("X-Trace-Id")

			if traceID == "" {
				traceID = uuid.New().String()
			}

			ctx := common.SetTraceID(r.Context(), traceID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
}

func LoggerMiddleware(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			traceID := common.GetTraceID(r.Context())
			logger = logger.With(slog.String("trace_id", traceID))
			ctx := common.SetLogger(r.Context(), logger)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
}
