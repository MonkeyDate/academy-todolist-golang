package cli

import (
	"academy-todo/pkg/todo"
	"fmt"
	"strings"
)

func TodoListCli(args []string, list todo.List) (modified bool, newList todo.List, err error) {
	// no args -> just list
	// add -> add item, then print
	// unknown args -> syntax message

	switch true {
	case len(list.Items) == 0 && len(args) == 0:
		fmt.Println(noTodoListHelp)
		return false, list, nil

	case len(args) == 0:
		PrintList(list)
		fmt.Println()
		fmt.Println("For details of supported commands use -h")
		return false, list, nil
	}

	switch strings.ToLower(args[0]) {
	case "add", "a":
		list, err := addItemToListCommand(list, args[1:])
		if err != nil {
			return false, list, err
		}

		PrintList(list)
		return true, list, nil

	case "update", "u":
		list, err := updateItemByIndexCommand(list, args[1:])
		if err != nil {
			return false, list, err
		}

		PrintList(list)
		return true, list, nil

	case "delete", "d":
		list, err := deleteItemByIndexCommand(list, args[1:])
		if err != nil {
			return false, list, err
		}

		PrintList(list)
		return true, list, nil

	default:
		fmt.Println(generalCommandHelp)
		return false, list, nil
	}
}

const noTodoListHelp = `There is no TODO list

Use "academy-todo add -h" for more information about adding items to the list.`

const generalCommandHelp = `Simple TODO list manager

Usage:

academy-todo <command> [arguments]

The commands are:

	add     add a new item to the list
	update  update an item in the list
	delete  delete an item form the list

Use "academy-todo <command> -h" for more information about a command.`
