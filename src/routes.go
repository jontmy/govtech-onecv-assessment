package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

var db *sql.DB

type Registration struct {
	Teacher  string   `json:"teacher"`
	Students []string `json:"students"`
}

// POST /api/register
func register(res http.ResponseWriter, req *http.Request) {
	// Check that the request method is POST.
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// Check that the request body is JSON.
	if req.Header.Get("Content-Type") != "application/json" {
		res.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	// Parse the request body.
	var registration Registration
	err := json.NewDecoder(req.Body).Decode(&registration)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	// Insert the teacher into the database if they don't already exist.
	_, err = db.Exec("INSERT IGNORE INTO teachers (email) VALUES (?)", registration.Teacher)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Insert the students into the database if they don't already exist.
	for _, student := range registration.Students {
		_, err = db.Exec("INSERT IGNORE INTO students (email) VALUES (?)", student)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	// Send the response.
	res.WriteHeader(http.StatusNoContent)
}

// GET /api/commonstudents
func commonStudents(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	res.WriteHeader(http.StatusOK)
}

// POST /api/suspend
func suspend(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	res.WriteHeader(http.StatusNoContent)
}

// POST /api/retrievefornotifications
func retrieveForNotifications(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	res.WriteHeader(http.StatusOK)
}
