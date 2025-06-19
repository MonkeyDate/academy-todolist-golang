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
	"strings"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	ctx := context.Background()
	traceId := uuid.New().String()
	ctx = context.WithValue(ctx, appContext.CtxTraceID{}, traceId)

	logger := createJsonLogger(traceId)
	ctx = context.WithValue(ctx, appContext.CtxLogger{}, *logger)
	logger.Info("Starting: " + strings.Join(os.Args[1:], " "))

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
