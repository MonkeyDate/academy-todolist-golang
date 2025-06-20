package main

import (
	"academy-todo/cli"
	appContext "academy-todo/context"
	"context"
	"fmt"
	"github.com/google/uuid"
	"io/fs"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	ctx, ctxCleanup := setupContext()
	defer ctxCleanup()
	defer func() { <-ctx.Done() }()

	todoList, err := LoadTodoList(ctx)
	if err != nil {
		fmt.Println("There was a problem loading the TODO list")
		fmt.Println(err)
		return
	}

	isModified, todoList, err := cli.TodoListCli(os.Args[1:], todoList)
	if err != nil {
		fmt.Println(err)
		return
	}

	if isModified {
		err = SaveTodoList(ctx, todoList)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func setupContext() (ctx context.Context, cleanup func()) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT)

	traceId := uuid.New().String()
	ctx = context.WithValue(ctx, appContext.CtxTraceID{}, traceId)

	logger, loggerCleanup := createJsonLogger(traceId)
	ctx = context.WithValue(ctx, appContext.CtxLogger{}, *logger)
	logger.Info("Starting: " + strings.Join(os.Args[1:], " "))

	cleanup = func() {
		stop()
		loggerCleanup()
	}

	return
}

func createJsonLogger(traceID string) (logger *slog.Logger, cleanup func()) {
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

const logFilename = "todolist.log"
