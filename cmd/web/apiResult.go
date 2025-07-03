package main

import "academy-todo/pkg/todo"

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type APIList struct {
	// TODO: WIP
	List todo.List `json:"list""`
	ID   string    `json:"id"`
}

const (
	ErrorGenericError     = -1
	ErrorInvalidParameter = 100
)
