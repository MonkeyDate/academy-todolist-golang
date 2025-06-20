package main

import (
	"academy-todo/cli"
	appContext "academy-todo/context"
	"academy-todo/storage"
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

	ctx, stop := setupContext()
	defer stop()

	todoList, err := storage.LoadTodoList(ctx)
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
		err = storage.SaveTodoList(ctx, todoList)
		if err != nil {
			fmt.Println(err)
		}
	}

	select {
	case <-ctx.Done():
	}
}

func setupContext() (context.Context, context.CancelFunc) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT)

	traceId := uuid.New().String()
	ctx = context.WithValue(ctx, appContext.CtxTraceID{}, traceId)

	logger := createJsonLogger(traceId)
	ctx = context.WithValue(ctx, appContext.CtxLogger{}, *logger)
	logger.Info("Starting: " + strings.Join(os.Args[1:], " "))

	return ctx, stop
}

func createJsonLogger(traceId string) *slog.Logger {
	var logger *slog.Logger
	logFile, _ := os.OpenFile(logFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, fs.ModePerm)
	if logFile != nil {
		logger = slog.New(slog.NewJSONHandler(logFile, nil))
	} else {
		logger = slog.Default()
	}

	logger = logger.With(slog.String("trace_id", traceId))

	return logger
}

const logFilename = "todolist.log"
