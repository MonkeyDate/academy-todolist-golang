package common

import (
	"academy-todo/pkg/todo"
	"context"
	"encoding/csv"
	"errors"
	"io"
	"io/fs"
	"os"
)

const filename string = "todolist.csv"

func SaveTodoList(ctx context.Context, list todo.List) error {
	logger := GetLogger(ctx)
	logger.Info("Saving TodoList...")

	f, err := os.Create(filename)
	if err != nil {
		logger.Error("TODO list cannot be saved", "filename", filename, "err", err)
		return err
	}

	defer func() { _ = f.Close() }()

	w := csv.NewWriter(f)
	for _, item := range list.Items {
		lineValues := []string{item.ID, string(item.Status), item.Description}
		err = w.Write(lineValues)
		if err != nil {
			logger.Error("TODO list cannot be saved", "filename", filename, "err", err)
			return err
		}
	}

	w.Flush()
	logger.Info("TodoList saved.")
	return nil
}

func LoadTodoList(ctx context.Context) (todo.List, error) {
	logger := GetLogger(ctx)

	f, err := os.Open(filename)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			logger.Info("TODO list does not exist", "filename", filename)
			return todo.List{}, nil
		}

		logger.Error("TODO list cannot be loaded", "filename", filename, "err", err)
		return todo.List{}, err
	}

	defer func() { _ = f.Close() }()

	r := csv.NewReader(f)
	todoItems := make([]todo.Item, 0)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(record) != 3 {
			logger.Error("TODO list cannot be loaded", "filename", filename, "err", err, "fields", len(record))
			return todo.List{}, err
		}

		todoItems = append(todoItems, todo.Item{ID: record[0], Status: todo.ItemStatus(record[1]), Description: record[2]})
	}

	todoList := todo.List{Items: todoItems}
	return todoList, nil
}
