package cli

import (
	"container/list"
)

func TodoListCli(commandLineFlags []string, todoList *list.List) (modified bool) {
	// no args -> just list
	// add -> add item, then print
	// unknown args -> syntax message

	if len(commandLineFlags) == 0 {
		// print list
		return false
	}

	return false
}
