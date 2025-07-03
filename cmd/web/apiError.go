package main

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

const (
	ErrorGenericError     = -1
	ErrorInvalidParameter = 100
)
