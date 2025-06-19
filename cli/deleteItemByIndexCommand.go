package cli

import (
	"academy-todo/models"
	"errors"
	"flag"
	"log/slog"
)

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
