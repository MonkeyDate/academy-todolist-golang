package storage

import (
	"academy-todo/models"
	"encoding/csv"
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
)

const filename string = "todolist.csv"

func SaveTodoList(list []models.TodoItem) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	// TODO: from example, not looked into it
	defer f.Close() // no error handling

	w := csv.NewWriter(f)
	for _, item := range list {
		lineValues := []string{string(item.Status), item.Description}
		err = w.Write(lineValues)
		if err != nil {
			return err
		}
	}

	w.Flush()
	return nil
}

func LoadTodoList() ([]models.TodoItem, error) {
	f, err := os.Open(filename)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	// TODO: look into defer and how we return errors from it
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	r := csv.NewReader(f)
	todoItems := make([]models.TodoItem, 0)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		todoItems = append(todoItems, models.TodoItem{Status: models.TodoItemStatus(record[0]), Description: record[1]})
	}

	return todoItems, nil
}
