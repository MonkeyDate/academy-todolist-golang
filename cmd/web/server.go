package main

import (
	"fmt"
	"net/http"
)

func Start() (err error) {
	// TODO: start the server on a seperate thread

	StartTodolistStoreActor()

	mux := http.NewServeMux()

	//
	// api crud endpoints
	// using GET for everything to keep testing easy
	mux.HandleFunc("GET /create", handleCreate) // TODO: lookup this syntax, check for param filtering etc
	mux.HandleFunc("/get", onlyOnGET(handleGet))
	mux.HandleFunc("/update", onlyOnGET(handleUpdate))
	mux.HandleFunc("/delete", onlyOnGET(handleDelete))

	//
	// static file
	fs := http.FileServer(http.Dir("./inetpub/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	//
	// template example
	mux.HandleFunc("/t/basic.html", setupBasicHtmlHandler())
	mux.HandleFunc("/t/list.html", setupTodolistTemplateHandler())

	fmt.Println("Server running on http://localhost:8080")
	serverStack := TraceIDMiddleware(LoggerMiddleware(mux))
	return http.ListenAndServe(":8080", serverStack)
}

func onlyOnGET(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	}
}
