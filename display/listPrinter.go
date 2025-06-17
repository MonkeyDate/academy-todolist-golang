package display

import (
	"academy-todo/models"
	"fmt"
)

func PrintList(todoList []models.TodoItem) {
	fmt.Println("<index>: <status> - <description>")

	for index, item := range todoList {
		fmt.Printf("%d: %s - %s\n", index, item.Status, item.Description)
	}
}
