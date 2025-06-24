package cli

import (
	"academy-todo/pkg/todo"
	"testing"
)

func TestTodoListCli_ShouldNotModifyTodoList_IfNoCommandLineFlagsPassed(t *testing.T) {
	flags := make([]string, 0)
	todoList := todo.List{Items: make([]todo.Item, 1)}
	todoList.Items[0] = todo.Item{
		Description: "desc-1",
		Status:      todo.Started,
	}

	isModified, _, _ := TodoListCli(flags, todoList)

	if isModified {
		t.Errorf("TodoListCli should not modify the list")
	}
}
