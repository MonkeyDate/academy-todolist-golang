package cli

import (
	"academy-todo/pkg/todo"
	"errors"
	"flag"
	"log/slog"
)

func addItemToListCommand(todoList []todo.Item, args []string) ([]todo.Item, error) {
	addCmd := flag.NewFlagSet("add", flag.ContinueOnError)
	description := addCmd.String("d", "new-item", "description of TODO item")
	started := addCmd.Bool("started", false, "has the TODO item already started")

	err := addCmd.Parse(args)
	if err != nil {
		return todoList, err
	}

	var status todo.ItemStatus
	if *started == true {
		status = todo.Started
	} else {
		status = todo.NotStarted
	}

	return append(todoList, todo.Item{Description: *description, Status: status}), nil
}

func deleteItemByIndexCommand(todoList []todo.Item, args []string) ([]todo.Item, error) {
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

func updateItemByIndexCommand(todoList []todo.Item, args []string) ([]todo.Item, error) {
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

	var status todo.ItemStatus
	if *complete == true {
		status = todo.Completed
	} else if *started == true {
		status = todo.Started
	} else {
		status = todo.NotStarted
	}

	itemToUpdate := &todoList[*index]

	itemToUpdate.Status = status

	if *description != "" {
		itemToUpdate.Description = *description
	}

	return todoList, nil
}
