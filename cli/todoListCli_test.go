package cli

import (
	"academy-todo/models"
	"testing"
)

func TestTodoListCli_ShouldNotModifyTodoList_IfNoCommandLineFlagsPassed(t *testing.T) {
	flags := make([]string, 0)
	todoList := make([]models.TodoItem, 1)

	_ = append(todoList, models.TodoItem{
		Description: "desc-1",
		Status:      models.Started,
	})

	isModified, _, _ := TodoListCli(flags, todoList)

	if isModified {
		t.Errorf("TodoListCli should not modify the list")
	}
}
