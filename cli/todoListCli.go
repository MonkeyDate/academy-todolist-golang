package cli

import (
	"academy-todo/display"
	"academy-todo/models"
)

func TodoListCli(args []string, todoList []models.TodoItem) (modified bool, list []models.TodoItem, err error) {
	// no args -> just list
	// add -> add item, then print
	// unknown args -> syntax message

	if len(args) == 0 {
		display.PrintList(todoList)
		return false, todoList, nil
	}

	switch args[0] {
	case "add":
		todoList, err := addItemToListCommand(todoList, args[1:])
		if err != nil {
			return false, todoList, err
		}

		display.PrintList(todoList)
		return true, todoList, nil

	case "update":
		todoList, err := updateItemByIndexCommand(todoList, args[1:])
		if err != nil {
			return false, todoList, err
		}

		display.PrintList(todoList)
		return true, todoList, nil

	case "delete":
		todoList, err := deleteItemByIndexCommand(todoList, args[1:])
		if err != nil {
			return false, todoList, err
		}

		display.PrintList(todoList)
		return true, todoList, nil
	}

	return false, todoList, nil
}
