package main

import (
	"academy-todo/cli"
	"academy-todo/storage"
	"container/list"
	"fmt"
	"os"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	todoList, err := storage.LoadTodoList()
	if err != nil {
		todoList = list.New()
	}

	isModified := cli.TodoListCli(os.Args[1:], todoList)

	if isModified {
		err = storage.SaveTodoList(todoList)
		if err != nil {
			fmt.Println(err)
		}
	}
}
