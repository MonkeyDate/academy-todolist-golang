package main

import (
	"fmt"
	"net/http"
)

func Start() (err error) {
	// TODO: start the server on a seperate thread

	mux := http.NewServeMux()
	mux.HandleFunc("/create", onlyOnGET(handleCreate))
	mux.HandleFunc("/get", onlyOnGET(handleGet))
	mux.HandleFunc("/update", onlyOnGET(handleUpdate))
	mux.HandleFunc("/delete", onlyOnGET(handleDelete))

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
