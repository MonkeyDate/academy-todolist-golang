package cli

import (
	"academy-todo/pkg/todo"
	"fmt"
	"strings"
)

func TodoListCli(args []string, todoList []todo.Item) (modified bool, list []todo.Item, err error) {
	// no args -> just list
	// add -> add item, then print
	// unknown args -> syntax message

	switch true {
	case todoList == nil && len(args) == 0:
		fmt.Println(noTodoListHelp)
		return false, todoList, nil

	case len(args) == 0:
		PrintList(todoList)
		fmt.Println()
		fmt.Println("For details of supported commands use -h")
		return false, todoList, nil
	}

	switch strings.ToLower(args[0]) {
	case "add", "a":
		todoList, err := addItemToListCommand(todoList, args[1:])
		if err != nil {
			return false, todoList, err
		}

		PrintList(todoList)
		return true, todoList, nil

	case "update", "u":
		todoList, err := updateItemByIndexCommand(todoList, args[1:])
		if err != nil {
			return false, todoList, err
		}

		PrintList(todoList)
		return true, todoList, nil

	case "delete", "d":
		todoList, err := deleteItemByIndexCommand(todoList, args[1:])
		if err != nil {
			return false, todoList, err
		}

		PrintList(todoList)
		return true, todoList, nil

	default:
		fmt.Println(generalCommandHelp)
		return false, todoList, nil
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
