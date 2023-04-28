package main

import (
	"encoding/json"
	"net/http"

	"cuongld2.com/api/pkg/models"
	"github.com/gorilla/mux"
)

func (app *application) all(w http.ResponseWriter, r *http.Request) {
	// Get all meetings stored
	meetings, err := app.meetings.All()
	if err != nil {
		app.serverError(w, err)
	}

	// Convert meeting list into json encoding
	b, err := json.Marshal(meetings)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Println("meetings have been listed")

	// Send response back
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) findByID(w http.ResponseWriter, r *http.Request) {
	// Get id from incoming url
	vars := mux.Vars(r)
	id := vars["id"]

	// Find meeting by id
	m, err := app.meetings.FindByID(id)
	if err != nil {
		if err.Error() == "ErrNoDocuments" {
			app.infoLog.Println("meeting not found")
			return
		}
		// Any other error will send an internal server error
		app.serverError(w, err)
	}

	// Convert meeting to json encoding
	b, err := json.Marshal(m)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Println("Have been found a meeting")

	// Send response back
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) insert(w http.ResponseWriter, r *http.Request) {
	// Define meeting model
	var m models.Meeting
	// Get request information
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		app.serverError(w, err)
	}

	// Insert new meeting
	insertResult, err := app.meetings.Insert(m)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Printf("New meeting have been created, id=%s", insertResult.InsertedID)
}

func (app *application) delete(w http.ResponseWriter, r *http.Request) {
	// Get id from incoming url
	vars := mux.Vars(r)
	id := vars["id"]

	// Delete meeting by id
	deleteResult, err := app.meetings.Delete(id)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Printf("Have been eliminated %d meeting(s)", deleteResult.DeletedCount)
}
