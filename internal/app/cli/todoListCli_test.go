package cli

import (
	"academy-todo/pkg/todo"
	"testing"
)

func TestTodoListCli_ShouldNotModifyTodoList_IfNoCommandLineFlagsPassed(t *testing.T) {
	flags := make([]string, 0)
	todoList := make([]todo.Item, 1)

	_ = append(todoList, todo.Item{
		Description: "desc-1",
		Status:      todo.Started,
	})

	isModified, _, _ := TodoListCli(flags, todoList)

	if isModified {
		t.Errorf("TodoListCli should not modify the list")
	}
}
