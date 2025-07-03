package main

import (
	"academy-todo/internal/common"
	"academy-todo/pkg/todo"
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

// TODO: no way to distinguish between internal errors and api errors

func handleCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	description := r.URL.Query().Get("description")
	status := parseStatus(r.URL.Query().Get("status"))

	result := CreateItem(ctx, description, status)
	if result.err != nil {
		writeAPIFailure(ctx, w, http.StatusInternalServerError, "Create: Failed to add item to TODO list", result.err, ErrorGenericError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(result.list)
	if err != nil {
		writeLogError(r.Context(), "Failed to encode TODO list", err, ErrorGenericError, http.StatusInternalServerError)
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := ReadItems(ctx)
	if result.err != nil {
		writeAPIFailure(ctx, w, http.StatusInternalServerError, "Get: Failed to load TODO list", result.err, ErrorGenericError)
		return
	}

	_ = writeTodoListSuccess(w, r, result.list)
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("ID")
	description := r.URL.Query().Get("description")
	status := parseStatus(r.URL.Query().Get("status"))

	if id == "" {
		writeAPIError(ctx, w, http.StatusBadRequest, "ID cannot be empty", ErrorInvalidParameter)
		return
	}

	result := UpdateItem(ctx, id, description, status)
	if result.err != nil {
		writeAPIFailure(ctx, w, http.StatusInternalServerError, "Update: Failed to update TODO item", result.err, ErrorGenericError)
		return
	}

	_ = writeTodoListSuccess(w, r, result.list)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("ID")
	if id == "" {
		writeAPIError(ctx, w, http.StatusBadRequest, "ID cannot be empty", ErrorInvalidParameter)
		return
	}

	result := DeleteItem(ctx, id)
	if result.err != nil {
		writeAPIFailure(ctx, w, http.StatusInternalServerError, "Delete: Failed to delete TODO item", result.err, ErrorGenericError)
		return
	}

	_ = writeTodoListSuccess(w, r, result.list)
}

func parseStatus(statusParam string) todo.ItemStatus {
	switch strings.ToLower(statusParam) {
	case "not-started", "not_started", "not started":

		return todo.NotStarted
	case "started":

		return todo.Started
	case "complete", "completed":
		return todo.Completed
	default:
		return todo.NotStarted
	}
}

func writeTodoListSuccess(w http.ResponseWriter, r *http.Request, list todo.List) error {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(list)
	if err != nil {
		writeLogError(r.Context(), "Failed to encode TODO list", err, ErrorGenericError, http.StatusInternalServerError)
	}
	return err
}

func writeAPIError(ctx context.Context, w http.ResponseWriter, status int, msg string, code int) {
	writeLogError(ctx, msg, nil, code, status)
	writeJSONError(w, status, APIError{Message: msg, Code: code})
}

func writeAPIFailure(ctx context.Context, w http.ResponseWriter, status int, msg string, sourceErr error, code int) {
	writeLogError(ctx, msg, sourceErr, code, status)
	writeJSONError(w, status, APIError{Message: msg, Code: code})
}

// TODO: rename
func writeLogError(ctx context.Context, msg string, sourceErr error, code, httpStatus int) {
	logger := common.GetLogger(ctx)
	args := []any{
		"httpStatusCode", httpStatus,
		"errorCode", code,
	}
	if sourceErr != nil {
		args = append(args, "sourceError", sourceErr)
	}
	logger.Error(msg, args...)
}

func writeJSONError(w http.ResponseWriter, status int, err APIError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(err)
}
