package main

import (
	"academy-todo/internal/common"
	"fmt"
	"net/http"
)

func Start() (err error) {
	logger, cleanup := common.CreateJsonLogger2()
	defer cleanup()

	StartTodolistStoreActor(logger)

	mux := http.NewServeMux()

	//
	// api crud endpoints
	// using GET for everything to keep testing easy
	mux.HandleFunc("GET /create", handleCreate) // TODO: lookup this syntax, check for param filtering etc
	mux.HandleFunc("GET /get", handleGet)
	mux.HandleFunc("GET /update/{ID}", handleUpdate)
	mux.HandleFunc("GET /delete/{ID}", handleDelete)

	//
	// static file
	fs := http.FileServer(http.Dir("./inetpub/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	//
	// template example
	mux.HandleFunc("/t/basic.html", setupBasicHtmlHandler())
	mux.HandleFunc("/t/list.html", setupTodolistTemplateHandler())

	fmt.Println("Server running on http://localhost:8080")
	serverStack := TraceIDMiddleware(LoggerMiddleware(mux, logger))
	return http.ListenAndServe(":8080", serverStack)
}
