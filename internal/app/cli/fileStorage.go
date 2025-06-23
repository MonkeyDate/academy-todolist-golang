package cli

import (
	"academy-todo/pkg/todo"
	"context"
	"encoding/csv"
	"errors"
	"io"
	"io/fs"
	"log/slog"
	"os"
)

const filename string = "todolist.csv"

func SaveTodoList(ctx context.Context, list []todo.Item) error {
	logger := ctx.Value(CtxLogger{}).(slog.Logger)
	logger.Info("Saving TodoList...")

	f, err := os.Create(filename)
	if err != nil {
		logger.Error("TODO list cannot be saved", "filename", filename, "err", err)
		return err
	}

	defer func() { _ = f.Close() }()

	w := csv.NewWriter(f)
	for _, item := range list {
		lineValues := []string{string(item.Status), item.Description}
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

func LoadTodoList(ctx context.Context) ([]todo.Item, error) {
	logger := ctx.Value(CtxLogger{}).(slog.Logger)

	f, err := os.Open(filename)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			logger.Info("TODO list does not exist", "filename", filename)
			return nil, nil
		}

		logger.Error("TODO list cannot be loaded", "filename", filename, "err", err)
		return nil, err
	}

	defer func() { _ = f.Close() }()

	r := csv.NewReader(f)
	todoItems := make([]todo.Item, 0)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Error("TODO list cannot be loaded", "filename", filename, "err", err)
			return nil, err
		}

		todoItems = append(todoItems, todo.Item{Status: todo.ItemStatus(record[0]), Description: record[1]})
	}

	return todoItems, nil
}
