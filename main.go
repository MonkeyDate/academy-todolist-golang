package main

import (
	"academy-todo/cli"
	"academy-todo/models"
	"academy-todo/storage"
	"fmt"
	"os"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	todoList, err := storage.LoadTodoList()
	if err != nil {
		// TODO: consider cli option to initialise list
		todoList = make([]models.TodoItem, 0)
	}

	isModified := cli.TodoListCli(os.Args[1:], todoList)

	todoList = append(todoList, models.TodoItem{
		Description: "Test item",
		Status:      models.NotStarted,
	})
	todoList = append(todoList, models.TodoItem{
		Description: "Another item",
		Status:      models.Started,
	})
	todoList = append(todoList, models.TodoItem{
		Description: "Yet another \"item",
		Status:      models.Started,
	})
	isModified = true

	if isModified {
		err = storage.SaveTodoList(todoList)
		if err != nil {
			fmt.Println(err)
		}
	}
}
