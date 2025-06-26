package main

import (
	"encoding/json"
	"strconv"
)

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	// Detail  string `json:"detail,omitempty"` // optional machine-readable detail
}

func BuildStructuredErrorJson(errorMessage string, errorCode int) string {
	structuredError := APIError{Message: errorMessage, Code: errorCode}
	if structuredJson, err := json.Marshal(structuredError); err == nil {
		return string(structuredJson)
	}

	return "{ \"code\": " + strconv.Itoa(errorBuildStructuredErrorJson) + " }"
}

const (
	ErrorGenericError             = -1
	errorBuildStructuredErrorJson = -2
	ErrorInvalidParameter         = 100
)
