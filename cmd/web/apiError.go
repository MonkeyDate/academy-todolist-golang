package main

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	// Detail  string `json:"detail,omitempty"` // optional machine-readable detail
}

const (
	ErrorGenericError     = -1
	ErrorInvalidParameter = 100
)
