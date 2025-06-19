package main

import (
	"academy-todo/cli"
	"academy-todo/storage"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"strings"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	logFile, _ := os.OpenFile(logFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, fs.ModePerm)
	if logFile != nil {
		defer func() { _ = logFile.Close() }()
		slog.SetDefault(slog.New(slog.NewJSONHandler(logFile, nil)))
	}

	slog.Info("Starting: ", strings.Join(os.Args[1:], " "))

	todoList, err := storage.LoadTodoList()
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
		err = storage.SaveTodoList(todoList)
		if err != nil {
			fmt.Println(err)
		}
	}
}

const logFilename = "todolist.log"
