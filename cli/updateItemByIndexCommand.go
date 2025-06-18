package cli

import (
	"academy-todo/models"
	"errors"
	"flag"
)

func updateItemByIndexCommand(todoList []models.TodoItem, args []string) ([]models.TodoItem, error) {
	addCmd := flag.NewFlagSet("update", flag.ExitOnError)
	index := addCmd.Int("i", -1, "index of TODO item to update")
	description := addCmd.String("d", "", "new description of TODO item, or blank")
	started := addCmd.Bool("started", false, "the TODO item has started")
	complete := addCmd.Bool("complete", false, "the TODO item has completed")

	err := addCmd.Parse(args)
	if err != nil {
		return todoList, err
	}

	if *index < 0 || *index >= len(todoList) {
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
