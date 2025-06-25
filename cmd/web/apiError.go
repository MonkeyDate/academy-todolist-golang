package main

import (
	"academy-todo/internal/common"
	"context"
	"encoding/json"
	"net/http"
)

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	// Detail  string `json:"detail,omitempty"` // optional machine-readable detail
}

func ReturnStructuredError(w http.ResponseWriter, httpStatusCode int, errorMessage string, errorCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)

	structuredError := APIError{Code: errorCode, Message: errorMessage}
	_ = json.NewEncoder(w).Encode(structuredError)
}

func LogStructuredError(ctx context.Context, httpStatusCode int, errorMessage string, errorCode int, sourceError error) {
	logger := common.GetLogger(ctx)

	args := []any{
		"httpStatusCode", httpStatusCode,
		"sourceError", sourceError,
	}

	structuredError := APIError{Code: errorCode, Message: errorMessage}
	if structuredJson, err := json.Marshal(structuredError); err == nil {
		args = append(args, "error", string(structuredJson))
	}

	logger.Error(errorMessage, args...)
}

const (
	ErrorGenericError = -1
)
