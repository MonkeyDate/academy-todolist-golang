package main

import (
	"academy-todo/internal/common"
	"academy-todo/pkg/todo"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

func handleCreate(w http.ResponseWriter, r *http.Request) {
	description := r.URL.Query().Get("description")
	statusParam := r.URL.Query().Get("status")

	if description == "" {
		description = "new-item-" + time.Now().Format(time.RFC3339)
	}

	status := parseStatus(statusParam)

	todoList, ok := loadTodolist(w, r)
	if !ok {
		return
	}

	todoList.Items = append(todoList.Items, todo.Item{Description: description, Status: status})

	ok = saveTodolist(w, r, todoList)
	if !ok {
		return
	}

	returnTodoListSuccess(w, r, todoList)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	todoList, ok := loadTodolist(w, r)
	if !ok {
		return
	}

	returnTodoListSuccess(w, r, todoList)
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	description := r.URL.Query().Get("description")
	statusParam := r.URL.Query().Get("status")

	if description == "" {
		description = "new-item-" + time.Now().Format(time.RFC3339)
	}

	status := parseStatus(statusParam)

	todoList, ok := loadTodolist(w, r)
	if !ok {
		return
	}

	todoList.Items = append(todoList.Items, todo.Item{Description: description, Status: status})

	ok = saveTodolist(w, r, todoList)
	if !ok {
		return
	}

	returnTodoListSuccess(w, r, todoList)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	description := r.URL.Query().Get("description")
	statusParam := r.URL.Query().Get("status")

	if description == "" {
		description = "new-item-" + time.Now().Format(time.RFC3339)
	}

	status := parseStatus(statusParam)

	todoList, ok := loadTodolist(w, r)
	if !ok {
		return
	}

	todoList.Items = append(todoList.Items, todo.Item{Description: description, Status: status})

	ok = saveTodolist(w, r, todoList)
	if !ok {
		return
	}

	returnTodoListSuccess(w, r, todoList)
}

func loadTodolist(w http.ResponseWriter, r *http.Request) (todoList todo.List, ok bool) {
	ctx := r.Context()

	todoList, err := common.LoadTodoList(ctx)
	if err != nil {
		logStructuredError(ctx, http.StatusInternalServerError, "Failed loading TODO list", ErrorGenericError, err)
		returnStructuredError(w, http.StatusInternalServerError, "Unable to load TODO list", ErrorGenericError)
		ok = false
		return
	}

	ok = true
	return
}

func saveTodolist(w http.ResponseWriter, r *http.Request, list todo.List) (ok bool) {
	ctx := r.Context()

	err := common.SaveTodoList(ctx, list)
	if err != nil {
		logStructuredError(ctx, http.StatusInternalServerError, "Failed to save TODO list", ErrorGenericError, err)
		returnStructuredError(w, http.StatusInternalServerError, "Unable to save TODO list", ErrorGenericError)
		ok = false
		return
	}

	ok = true
	return
}

func parseStatus(statusParam string) todo.ItemStatus {
	var status todo.ItemStatus

	switch strings.ToLower(statusParam) {
	case "not-started", "not_started", "not started":
		status = todo.NotStarted
		break
	case "started":
		status = todo.Started
		break
	case "complete", "completed":
		status = todo.Completed
		break
	default:
		status = todo.NotStarted
	}
	return status
}

func returnTodoListSuccess(w http.ResponseWriter, r *http.Request, todoList todo.List) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todoList); err != nil {
		logStructuredError(ctx, http.StatusInternalServerError, "Failed to encode JSON", ErrorGenericError, err)
		returnStructuredError(w, http.StatusInternalServerError, "Unable to return updated TODO list", ErrorGenericError)
		return
	}
}

func returnStructuredError(w http.ResponseWriter, httpStatusCode int, errorMessage string, errorCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	_ = json.NewEncoder(w).Encode(BuildStructuredErrorJson(errorMessage, errorCode))
}

func logStructuredError(ctx context.Context, httpStatusCode int, errorMessage string, errorCode int, sourceError error) {
	args := []any{
		"httpStatusCode", httpStatusCode,
		"sourceError", sourceError,
		"error", BuildStructuredErrorJson(errorMessage, errorCode),
	}

	logger := common.GetLogger(ctx)
	logger.Error(errorMessage, args...)
}
