package cli

import (
	"academy-todo/models"
	"errors"
	"flag"
)

func deleteItemByIndexCommand(todoList []models.TodoItem, args []string) ([]models.TodoItem, error) {
	deleteCmd := flag.NewFlagSet("delete", flag.ContinueOnError)
	index := deleteCmd.Int("i", -1, "index of TODO item to update")

	err := deleteCmd.Parse(args)
	if err != nil {
		return todoList, err
	}

	if *index < 0 || *index >= len(todoList) {
		return todoList, errors.New("index must reference an item in the TODO list")
	}

	return append(todoList[:*index], todoList[*index+1:]...), nil
}
