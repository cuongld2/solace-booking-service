package main

import (
	"github.com/gorilla/mux"
)

func (app *application) routes() *mux.Router {
	// Register handler functions.
	r := mux.NewRouter()
	r.HandleFunc("/api/meetings/", app.all).Methods("GET")
	r.HandleFunc("/api/meetings/{id}", app.findByID).Methods("GET")
	r.HandleFunc("/api/meetings/", app.insert).Methods("POST")
	r.HandleFunc("/api/meetings/{id}", app.delete).Methods("DELETE")

	return r
}
