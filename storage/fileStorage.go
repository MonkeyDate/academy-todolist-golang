package storage

import (
	"container/list"
	"fmt"
)

type FileStorage struct{}

func SaveTodoList(list *list.List) error {
	return fmt.Errorf("not implemented")
}

func LoadTodoList() (*list.List, error) {
	return nil, fmt.Errorf("not implemented")
}
