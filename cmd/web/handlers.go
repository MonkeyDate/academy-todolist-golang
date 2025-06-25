package main

import (
	"academy-todo/internal/common"
	"encoding/json"
	"net/http"
)

func handleCreate(w http.ResponseWriter, r *http.Request) {
	ReturnStructuredError(w, http.StatusNotImplemented, "CREATE not supported", ErrorGenericError)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	todoList, err := common.LoadTodoList(r.Context())
	if err != nil {
		LogStructuredError(r.Context(), http.StatusInternalServerError, "There was a problem loading the TODO list", ErrorGenericError, err)
		ReturnStructuredError(w, http.StatusInternalServerError, "There was a problem loading the TODO list", ErrorGenericError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todoList); err != nil {
		LogStructuredError(r.Context(), http.StatusInternalServerError, "Failed to encode JSON", ErrorGenericError, err)
		ReturnStructuredError(w, http.StatusInternalServerError, "Failed to encode JSON", ErrorGenericError)
		return
	}
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	ReturnStructuredError(w, http.StatusNotImplemented, "UPDATE not supported", ErrorGenericError)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	ReturnStructuredError(w, http.StatusNotImplemented, "DELETE not supported", ErrorGenericError)
}
