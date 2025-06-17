package cli

import (
	"academy-todo/models"
	"container/list"
	"testing"
)

func TestTodoListCli_ShouldNotModifyTodoList_IfNoCommandLineFlagsPassed(t *testing.T) {
	flags := make([]string, 0)
	todoList := list.New()

	todoList.PushFront(models.TodoItem{
		Description: "desc-1",
		Status:      models.Started,
	})

	isModified := TodoListCli(flags, todoList)

	if isModified {
		t.Errorf("TodoListCli should not modify the list")
	}
}
