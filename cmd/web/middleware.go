package main

import (
	"academy-todo/internal/common"
	"context"
	"github.com/google/uuid"
	"net/http"
)

func TraceIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: get traceid from header first
		traceID := uuid.New().String()

		ctx := context.WithValue(r.Context(), common.CtxTraceID{}, traceID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID, _ := r.Context().Value(common.CtxTraceID{}).(string)
		logger, loggerCleanup := common.CreateJsonLogger(traceID)

		ctx := common.SetLogger(r.Context(), logger)
		go func() {
			<-ctx.Done()
			loggerCleanup()
		}()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
