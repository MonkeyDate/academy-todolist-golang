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

// handleCreate processes a request to create a new TODO item with a description and status from the query parameters.
// It adds the item and responds with the updated TODO list in JSON format.
// If an error occurs during creation or encoding the response, it logs the error and sends an appropriate API failure response.
func handleCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	description := r.URL.Query().Get("description")
	status := parseStatus(r.URL.Query().Get("status"))

	//
	// create a wrapped context that has a timeout/timelimit and pass that to CreateItem
	// allowing us to timeout the operation without the operation having to explicitly
	// have a timeout

	result := CreateItem(ctx, description, status)
	if result.Err != nil {
		writeAPIFailure(ctx, w, http.StatusInternalServerError, "Create: Failed to add item to TODO list", result.Err, ErrorGenericError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(result.List)
	if err != nil {
		writeLogError(r.Context(), "Failed to encode TODO list", err, ErrorGenericError, http.StatusInternalServerError)
	}
}

// handleGet handles GET requests to retrieve the TODO list and responds with the list in JSON format.
// If an error occurs during fetching or encoding the response, it logs the error and sends an API failure response.
func handleGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := ReadItems(ctx)
	if result.Err != nil {
		writeAPIFailure(ctx, w, http.StatusInternalServerError, "Get: Failed to load TODO list", result.Err, ErrorGenericError)
		return
	}

	_ = writeTodoListSuccess(w, r, result.List)
}

// handleUpdate processes a request to update a TODO item identified by its ID with new description and status.
// It updates the item, and responds with the updated TODO list in JSON format.
// If validation or the update process fails, it logs the error and sends an appropriate API failure response.
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
	if result.Err != nil {
		writeAPIFailure(ctx, w, http.StatusInternalServerError, "Update: Failed to update TODO item", result.Err, ErrorGenericError)
		return
	}

	_ = writeTodoListSuccess(w, r, result.List)
}

// handleDelete processes a DELETE request for a TODO item identified by its ID and responds with the updated TODO list.
// If the ID is missing or deletion fails, it logs the error and sends an appropriate API failure response.
func handleDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("ID")
	if id == "" {
		writeAPIError(ctx, w, http.StatusBadRequest, "ID cannot be empty", ErrorInvalidParameter)
		return
	}

	result := DeleteItem(ctx, id)
	if result.Err != nil {
		writeAPIFailure(ctx, w, http.StatusInternalServerError, "Delete: Failed to delete TODO item", result.Err, ErrorGenericError)
		return
	}

	_ = writeTodoListSuccess(w, r, result.List)
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
func writeLogError(ctx context.Context, msg string, sourceErr error, code int, httpStatus int) {
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
