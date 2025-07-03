package common

import (
	"io/fs"
	"log/slog"
	"os"
)

func CreateJsonLogger(traceID string) (logger *slog.Logger, cleanup func()) {
	logFile, _ := os.OpenFile(logFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, fs.ModePerm)
	cleanup = func() {
		if logFile != nil {
			_ = logFile.Close()
		}
	}

	if logFile != nil {
		logger = slog.New(slog.NewJSONHandler(logFile, nil))
	} else {
		logger = slog.Default()
	}

	logger = logger.With(slog.String("trace_id", traceID))

	return
}

func CreateJsonLogger2() (logger *slog.Logger, cleanup func()) {
	logFile, _ := os.OpenFile(logFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, fs.ModePerm)
	cleanup = func() {
		if logFile != nil {
			_ = logFile.Close()
		}
	}

	if logFile != nil {
		logger = slog.New(slog.NewJSONHandler(logFile, nil))
	} else {
		logger = slog.Default()
	}

	return
}

const logFilename = "todolist.log"
