package models

type TodoItemStatus string

const (
	NotStarted TodoItemStatus = "not started"
	Started    TodoItemStatus = "started"
	Completed  TodoItemStatus = "completed"
)

type TodoItem struct {
	Description string
	Status      TodoItemStatus
}
