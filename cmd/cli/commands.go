package main

import (
	"academy-todo/pkg/todo"
	"errors"
	"flag"
	"log/slog"
)

func addItemToListCommand(todoList todo.List, args []string) (todo.List, error) {
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

	todoList.Items = append(todoList.Items, todo.Item{Description: *description, Status: status})
	return todoList, nil
}

func deleteItemByIndexCommand(todoList todo.List, args []string) (todo.List, error) {
	deleteCmd := flag.NewFlagSet("delete", flag.ContinueOnError)
	index := deleteCmd.Int("i", -1, "index of TODO item to update")

	err := deleteCmd.Parse(args)
	if err != nil {
		slog.Error("item cannot be removed from TODO list", "err", err)
		return todoList, err
	}

	if *index < 0 || *index >= len(todoList.Items) {
		slog.Warn("item cannot be removed form TODO list, bad index", "index", *index, "todListSize", len(todoList.Items))
		return todoList, errors.New("index must reference an item in the TODO list")
	}

	todoList.Items = append(todoList.Items[:*index], todoList.Items[*index+1:]...)
	return todoList, nil
}

func updateItemByIndexCommand(todoList todo.List, args []string) (todo.List, error) {
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

	if *index < 0 || *index >= len(todoList.Items) {
		slog.Warn("item cannot be updated in TODO list, bad index", "index", *index, "todListSize", len(todoList.Items))
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

	itemToUpdate := &todoList.Items[*index]

	itemToUpdate.Status = status

	if *description != "" {
		itemToUpdate.Description = *description
	}

	return todoList, nil
}
