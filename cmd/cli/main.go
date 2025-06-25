package main

import (
	"academy-todo/internal/common"
	"context"
	"fmt"
	"github.com/google/uuid"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	ctx := setupContext()

	todoList, err := common.LoadTodoList(ctx)
	if err != nil {
		fmt.Println("There was a problem loading the TODO list")
		fmt.Println(err)
		return
	}

	isModified, todoList, err := TodoListCli(os.Args[1:], todoList)
	if err != nil {
		fmt.Println(err)
		return
	}

	if isModified {
		err = common.SaveTodoList(ctx, todoList)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func setupContext() (ctx context.Context) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT)

	traceId := uuid.New().String()
	ctx = context.WithValue(ctx, common.CtxTraceID{}, traceId)

	logger, loggerCleanup := common.CreateJsonLogger(traceId)
	go func() {
		<-ctx.Done() // blocks until context is canceled
		stop()
		loggerCleanup()
	}()

	ctx = common.SetLogger(ctx, logger)
	logger.Info("Starting: " + strings.Join(os.Args[1:], " "))

	return
}
