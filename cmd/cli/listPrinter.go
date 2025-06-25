package main

import (
	"academy-todo/pkg/todo"
	"fmt"
)

func PrintList(todoList todo.List) {
	if len(todoList.Items) == 0 {
		fmt.Println("Great work, your TODO list is empty.")
		return
	}

	fmt.Println("<index>: <status> - <description>")

	for index, item := range todoList.Items {
		var statusColour string
		var descriptionColour = colorReset

		switch item.Status {
		case todo.NotStarted:
			statusColour = colorRed
			break

		case todo.Started:
			statusColour = colorYellow
			break

		case todo.Completed:
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
