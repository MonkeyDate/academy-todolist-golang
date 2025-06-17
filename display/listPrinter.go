package display

import (
	"academy-todo/models"
	"fmt"
)

func PrintList(todoList []models.TodoItem) {
	if len(todoList) == 0 {
		fmt.Println("Great work, your TODO list is empty.")
		return
	}

	fmt.Println("<index>: <status> - <description>")

	for index, item := range todoList {
		fmt.Printf("%d: %s - %s\n", index, item.Status, item.Description)
	}
}
