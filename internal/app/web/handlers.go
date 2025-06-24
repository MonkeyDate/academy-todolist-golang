package web

import (
	"academy-todo/internal/common"
	"encoding/json"
	"fmt"
	"net/http"
)

func handleCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"status":"ok"}`)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	ctx, ctxCleanup := setupContext()
	defer ctxCleanup()
	defer func() { <-ctx.Done() }()

	todoList, err := common.LoadTodoList(ctx)
	if err != nil {
		fmt.Println("There was a problem loading the TODO list")
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todoList); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {}

func handleDelete(w http.ResponseWriter, r *http.Request) {}
