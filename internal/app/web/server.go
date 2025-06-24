package web

import (
	"academy-todo/internal/common"
	"context"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func Start() (err error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/create", onlyOnGET(handleCreate))
	mux.HandleFunc("/get", onlyOnGET(handleGet))
	mux.HandleFunc("/update", onlyOnGET(handleUpdate))
	mux.HandleFunc("/delete", onlyOnGET(handleDelete))

	fmt.Println("Server running on http://localhost:8080")
	return http.ListenAndServe(":8080", mux)
}

func onlyOnGET(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	}
}

func setupContext() (ctx context.Context, cleanup func()) {
	// TODO: just copied from the cli version
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT)

	traceId := uuid.New().String()
	ctx = context.WithValue(ctx, common.CtxTraceID{}, traceId)

	logger, loggerCleanup := common.CreateJsonLogger(traceId)
	ctx = context.WithValue(ctx, common.CtxLogger{}, *logger)
	logger.Info("Starting: " + strings.Join(os.Args[1:], " "))

	cleanup = func() {
		stop()
		loggerCleanup()
	}

	return
}
