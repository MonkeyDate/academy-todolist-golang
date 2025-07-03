package main

import (
	"academy-todo/internal/common"
	"academy-todo/pkg/todo"
	"encoding/json"
	"net/http"
	"strings"
)

// TODO: no way to distinguish between internal errors and api errors

func handleCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := common.GetLogger(ctx)
	description := r.URL.Query().Get("description")
	statusParam := r.URL.Query().Get("status")

	status := parseStatus(statusParam)

	result := CreateItem(ctx, description, status)
	if result.err != nil {
		logger.Error("Failed to add item to TODO list", "httpStatusCode", http.StatusInternalServerError, "sourceError", result.err, "errorCode", ErrorGenericError)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(APIError{Message: "Failed to add item to TODO list", Code: ErrorGenericError})
		return
	} else {
		_ = returnTodoListSuccess(w, r, result.list)
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := common.GetLogger(ctx)

	result := ReadItems(ctx)
	if result.err != nil {
		logger.Error("Failed to load TODO list", "httpStatusCode", http.StatusInternalServerError, "sourceError", result.err, "errorCode", ErrorGenericError)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(APIError{Message: "Failed to load TODO list", Code: ErrorGenericError})
		return
	} else {
		_ = returnTodoListSuccess(w, r, result.list)
	}
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := common.GetLogger(ctx)

	id := r.PathValue("ID")
	description := r.URL.Query().Get("description")
	statusParam := r.URL.Query().Get("status")
	status := parseStatus(statusParam)

	if id == "" {
		err := APIError{
			Message: "ID cannot be empty",
			Code:    ErrorInvalidParameter,
		}
		logger.Error(err.Message, "httpStatusCode", http.StatusBadRequest, "sourceError", err.Code)
		writeJSONError(w, http.StatusBadRequest, err)
		return
	}

	result := UpdateItem(ctx, id, description, status)
	if result.err != nil {
		err := APIError{
			Message: "Failed to update TODO item",
			Code:    ErrorGenericError,
		}
		logger.Error(err.Message, "httpStatusCode", http.StatusInternalServerError, "sourceError", result.err, "errorCode", err.Code)
		writeJSONError(w, http.StatusInternalServerError, err)
		return
	}

	_ = returnTodoListSuccess(w, r, result.list)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := common.GetLogger(ctx)

	id := r.PathValue("ID")
	if id == "" {
		err := APIError{
			Message: "ID cannot be empty",
			Code:    ErrorInvalidParameter,
		}
		logger.Error(err.Message, "httpStatusCode", http.StatusBadRequest, "sourceError", err.Code)
		writeJSONError(w, http.StatusBadRequest, err)
		return
	}

	result := DeleteItem(ctx, id)
	if result.err != nil {
		err := APIError{
			Message: "Failed to delete TODO item",
			Code:    ErrorGenericError,
		}
		logger.Error(err.Message, "httpStatusCode", http.StatusInternalServerError, "sourceError", result.err, "errorCode", err.Code)
		writeJSONError(w, http.StatusInternalServerError, err)
		return
	}

	_ = returnTodoListSuccess(w, r, result.list)
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

func returnTodoListSuccess(w http.ResponseWriter, r *http.Request, todoList todo.List) error {
	w.Header().Set("Content-Type", "application/json")

	var err error
	if err = json.NewEncoder(w).Encode(todoList); err != nil {
		logger := common.GetLogger(r.Context())
		logger.Error("Failed to encode JSON", "httpStatusCode", http.StatusInternalServerError, "sourceError", err, "errorCode", ErrorGenericError)
	}

	return err
}

func writeJSONError(w http.ResponseWriter, status int, err APIError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(err)
}
