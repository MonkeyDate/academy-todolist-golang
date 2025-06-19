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
		var statusColour string
		var descriptionColour = colorReset

		switch item.Status {
		case models.NotStarted:
			statusColour = colorRed
			break

		case models.Started:
			statusColour = colorYellow
			break

		case models.Completed:
			statusColour = colorGreen
			descriptionColour = colorDarkGrey
			break
		}

		fmt.Printf("%3d%s: %s%11s %s- %s%s%s\n",
			index, colorLightGrey,
			statusColour, item.Status,
			colorLightGrey,
			descriptionColour, item.Description, colorReset,
		)
	}
}

const colorReset = "\033[0m"
const colorRed = "\033[31m"
const colorYellow = "\033[33m"
const colorGreen = "\033[32m"
const colorLightGrey = "\033[38;5;7m"
const colorDarkGrey = "\033[38;5;7m"
