package main

import (
	appContext "academy-todo/context"
	"academy-todo/models"
	"context"
	"encoding/csv"
	"errors"
	"io"
	"io/fs"
	"log/slog"
	"os"
)

const filename string = "todolist.csv"

func SaveTodoList(ctx context.Context, list []models.TodoItem) error {
	logger := ctx.Value(appContext.CtxLogger{}).(slog.Logger)
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

func LoadTodoList(ctx context.Context) ([]models.TodoItem, error) {
	logger := ctx.Value(appContext.CtxLogger{}).(slog.Logger)

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
	todoItems := make([]models.TodoItem, 0)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Error("TODO list cannot be loaded", "filename", filename, "err", err)
			return nil, err
		}

		todoItems = append(todoItems, models.TodoItem{Status: models.TodoItemStatus(record[0]), Description: record[1]})
	}

	return todoItems, nil
}
