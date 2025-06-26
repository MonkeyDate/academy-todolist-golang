package main

import (
	"academy-todo/internal/common"
	"html/template"
	"net/http"
	"strings"
	"time"
)

func setupBasicHtmlHandler() http.HandlerFunc {
	funcs := template.FuncMap{
		"upper": strings.ToUpper,
		"now":   func() string { return time.Now().Format(time.RFC1123) },
	}

	tmpl := template.Must(
		template.New("basic.html").
			Funcs(funcs).
			ParseFiles("./inetpub/templates/basic.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		logger := common.GetLogger(r.Context())

		data := struct {
			Name      string
			Timestamp string
		}{
			Name:      "Gary",
			Timestamp: time.Now().Format(time.RFC1123),
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, data); err != nil {
			logger.Warn("index required and must be a number", "httpStatusCode", http.StatusBadRequest, "sourceError", err, "errorCode", ErrorInvalidParameter)
			http.Error(w, "Template rendering error", http.StatusInternalServerError)
		}
	}
}
