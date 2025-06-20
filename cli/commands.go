package cli

import (
	"academy-todo/models"
	"errors"
	"flag"
	"log/slog"
)

func addItemToListCommand(todoList []models.TodoItem, args []string) ([]models.TodoItem, error) {
	addCmd := flag.NewFlagSet("add", flag.ContinueOnError)
	description := addCmd.String("d", "new-item", "description of TODO item")
	started := addCmd.Bool("started", false, "has the TODO item already started")

	err := addCmd.Parse(args)
	if err != nil {
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

func deleteItemByIndexCommand(todoList []models.TodoItem, args []string) ([]models.TodoItem, error) {
	deleteCmd := flag.NewFlagSet("delete", flag.ContinueOnError)
	index := deleteCmd.Int("i", -1, "index of TODO item to update")

	err := deleteCmd.Parse(args)
	if err != nil {
		slog.Error("item cannot be removed from TODO list", "err", err)
		return todoList, err
	}

	if *index < 0 || *index >= len(todoList) {
		slog.Warn("item cannot be removed form TODO list, bad index", "index", *index, "todListSize", len(todoList))
		return todoList, errors.New("index must reference an item in the TODO list")
	}

	return append(todoList[:*index], todoList[*index+1:]...), nil
}

func updateItemByIndexCommand(todoList []models.TodoItem, args []string) ([]models.TodoItem, error) {
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	index := updateCmd.Int("i", -1, "index of TODO item to update")
	description := updateCmd.String("d", "", "new description of TODO item, or blank")
	started := updateCmd.Bool("started", false, "the TODO item has started")
	complete := updateCmd.Bool("complete", false, "the TODO item has completed")

	err := updateCmd.Parse(args)
	if err != nil {
		slog.Error("item cannot be updated in TODO list", "err", err)
		return todoList, err
	}

	if *index < 0 || *index >= len(todoList) {
		slog.Warn("item cannot be updated in TODO list, bad index", "index", *index, "todListSize", len(todoList))
		return todoList, errors.New("index must reference an item in the TODO list")
	}

	var status models.TodoItemStatus
	if *complete == true {
		status = models.Completed
	} else if *started == true {
		status = models.Started
	} else {
		status = models.NotStarted
	}

	itemToUpdate := &todoList[*index]

	itemToUpdate.Status = status

	if *description != "" {
		itemToUpdate.Description = *description
	}

	return todoList, nil
}
