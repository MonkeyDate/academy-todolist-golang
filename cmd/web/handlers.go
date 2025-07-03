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

	resultChan := BeginCreateItem(ctx, description, status)
	result := <-resultChan

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

	resultChan := BeginReadItems(ctx)
	result := <-resultChan

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

	description := r.URL.Query().Get("description")
	statusParam := r.URL.Query().Get("status")
	ID := r.URL.Query().Get("ID")

	if len(ID) == 0 {
		logger.Error("ID cannot be empty string", "httpStatusCode", http.StatusBadRequest, "sourceError", ErrorInvalidParameter)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(APIError{Message: "ID cannot be empty string", Code: ErrorInvalidParameter})
		return
	}

	status := parseStatus(statusParam)

	resultChan := BeginUpdateItem(ctx, ID, description, status)
	result := <-resultChan

	if result.err != nil {
		logger.Error("Failed to update TODO list", "httpStatusCode", http.StatusInternalServerError, "sourceError", result.err, "errorCode", ErrorGenericError)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(APIError{Message: "Failed to update TODO list", Code: ErrorGenericError})
		return
	} else {
		_ = returnTodoListSuccess(w, r, result.list)
	}
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := common.GetLogger(ctx)

	ID := r.URL.Query().Get("ID")

	if len(ID) == 0 {
		logger.Error("ID cannot be empty string", "httpStatusCode", http.StatusBadRequest, "sourceError", ErrorInvalidParameter)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(APIError{Message: "ID cannot be empty string", Code: ErrorInvalidParameter})
		return
	}

	resultChan := BeginDeleteItem(ctx, ID)
	result := <-resultChan

	if result.err != nil {
		logger.Error("Failed to update TODO list", "httpStatusCode", http.StatusInternalServerError, "sourceError", result.err, "errorCode", ErrorGenericError)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(APIError{Message: "Failed to update TODO list", Code: ErrorGenericError})
		return
	} else {
		_ = returnTodoListSuccess(w, r, result.list)
	}
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
