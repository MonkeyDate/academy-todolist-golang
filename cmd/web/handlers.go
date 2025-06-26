package main

import (
	"academy-todo/internal/common"
	"academy-todo/pkg/todo"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func handleCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := common.GetLogger(ctx)
	description := r.URL.Query().Get("description")
	statusParam := r.URL.Query().Get("status")

	if description == "" {
		description = "new-item-" + time.Now().Format(time.RFC3339)
	}

	status := parseStatus(statusParam)

	todoList, err := common.LoadTodoList(ctx)
	if err != nil {
		logger.Error("Failed loading TODO list", "httpStatusCode", http.StatusInternalServerError, "sourceError", err, "errorCode", ErrorGenericError)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(APIError{Message: "Unable to load TODO list", Code: ErrorGenericError})
		return
	}

	todoList.Items = append(todoList.Items, todo.Item{Description: description, Status: status})

	err = common.SaveTodoList(ctx, todoList)
	if err != nil {
		logger.Error("Failed to save TODO list", "httpStatusCode", http.StatusInternalServerError, "sourceError", err, "errorCode", ErrorGenericError)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(APIError{Message: "Unable to save TODO list", Code: ErrorGenericError})
		return
	}

	_ = returnTodoListSuccess(w, r, todoList)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	todoList, err := common.LoadTodoList(ctx)
	if err != nil {
		logger := common.GetLogger(ctx)
		logger.Error("Failed loading TODO list", "httpStatusCode", http.StatusInternalServerError, "sourceError", err, "errorCode", ErrorGenericError)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(APIError{Message: "Unable to load TODO list", Code: ErrorGenericError})
		return
	}

	_ = returnTodoListSuccess(w, r, todoList)
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := common.GetLogger(ctx)
	description := r.URL.Query().Get("description")
	statusParam := r.URL.Query().Get("status")
	indexParam := r.URL.Query().Get("index")

	index, err := strconv.Atoi(indexParam)
	if err != nil {
		logger.Error("index required and must be a number", "httpStatusCode", http.StatusBadRequest, "sourceError", err, "errorCode", ErrorInvalidParameter)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(APIError{Message: "index required and must be a number", Code: ErrorInvalidParameter})
		return
	}

	if description == "" {
		description = "new-item-" + time.Now().Format(time.RFC3339)
	}

	status := parseStatus(statusParam)

	todoList, err := common.LoadTodoList(ctx)
	if err != nil {
		logger := common.GetLogger(ctx)
		logger.Error("Failed loading TODO list", "httpStatusCode", http.StatusInternalServerError, "sourceError", err, "errorCode", ErrorGenericError)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(APIError{Message: "Unable to load TODO list", Code: ErrorGenericError})
		return
	}

	if index < 0 || index >= len(todoList.Items) {
		logger.Warn("item cannot be updated, bad index", "httpStatusCode", http.StatusBadRequest, "sourceError", err, "errorCode", ErrorInvalidParameter)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(APIError{Message: "item cannot be updated, bad index", Code: ErrorInvalidParameter})
		return
	}

	itemToUpdate := &todoList.Items[index]
	itemToUpdate.Status = status

	if description != "" {
		itemToUpdate.Description = description
	}

	err = common.SaveTodoList(ctx, todoList)
	if err != nil {
		logger.Error("Failed to save TODO list", "httpStatusCode", http.StatusInternalServerError, "sourceError", err, "errorCode", ErrorGenericError)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(APIError{Message: "Unable to save TODO list", Code: ErrorGenericError})
		return
	}

	_ = returnTodoListSuccess(w, r, todoList)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := common.GetLogger(ctx)
	indexParam := r.URL.Query().Get("index")

	index, err := strconv.Atoi(indexParam)
	if err != nil {
		logger.Warn("index required and must be a number", "httpStatusCode", http.StatusBadRequest, "sourceError", err, "errorCode", ErrorInvalidParameter)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(APIError{Message: "index required and must be a number", Code: ErrorInvalidParameter})
		return
	}

	todoList, err := common.LoadTodoList(ctx)
	if err != nil {
		logger.Error("Failed loading TODO list", "httpStatusCode", http.StatusInternalServerError, "sourceError", err, "errorCode", ErrorGenericError)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(APIError{Message: "Unable to load TODO list", Code: ErrorGenericError})
		return
	}

	if index < 0 || index >= len(todoList.Items) {
		logger.Warn("item cannot be removed form TODO list, bad index", "httpStatusCode", http.StatusBadRequest, "sourceError", err, "errorCode", ErrorInvalidParameter)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(APIError{Message: "item cannot be removed form TODO list, bad index", Code: ErrorInvalidParameter})
		return
	}

	todoList.Items = append(todoList.Items[:index], todoList.Items[index+1:]...)

	err = common.SaveTodoList(ctx, todoList)
	if err != nil {
		logger.Error("Failed to save TODO list", "httpStatusCode", http.StatusInternalServerError, "sourceError", err, "errorCode", ErrorGenericError)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(APIError{Message: "Unable to save TODO list", Code: ErrorGenericError})
		return
	}

	_ = returnTodoListSuccess(w, r, todoList)
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
