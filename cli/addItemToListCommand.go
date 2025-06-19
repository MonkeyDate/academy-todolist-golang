package cli

import (
	"academy-todo/models"
	"flag"
	"log/slog"
)

func addItemToListCommand(todoList []models.TodoItem, args []string) ([]models.TodoItem, error) {
	addCmd := flag.NewFlagSet("add", flag.ContinueOnError)
	description := addCmd.String("d", "new-item", "description of TODO item")
	started := addCmd.Bool("started", false, "has the TODO item already started")

	err := addCmd.Parse(args)
	if err != nil {
		slog.Error("item cannot be added to TODO list", "err", err)
		return todoList, err
	}

	var status models.TodoItemStatus
	if *started == true {
		status = models.Started
	} else {
		status = models.NotStarted
	}

	return append(todoList, models.TodoItem{Description: *description, Status: status}), nil
}
