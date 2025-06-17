package cli

import (
	"academy-todo/models"
)

func TodoListCli(commandLineFlags []string, todoList []models.TodoItem) (modified bool) {
	// no args -> just list
	// add -> add item, then print
	// unknown args -> syntax message

	if len(commandLineFlags) == 0 {
		// print list
		return false
	}

	return false
}
