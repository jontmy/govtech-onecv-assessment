package main

import (
	"database/sql"
	"encoding/json"
	"github.com/jmoiron/sqlx"
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
	_, err = db.Exec("INSERT IGNORE INTO teachers (teacher_email) VALUES (?)", registration.Teacher)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Insert the students into the database if they don't already exist.
	for _, student := range registration.Students {
		_, err = db.Exec(`
			INSERT IGNORE INTO students (student_email) VALUES (?);
		`, student)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = db.Exec(`
			INSERT IGNORE INTO class (teacher_email, student_email) VALUES (?, ?);
		`, registration.Teacher, student)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	// Send the response.
	res.WriteHeader(http.StatusNoContent)
}

type Students struct {
	Students []string `json:"students"`
}

// GET /api/commonstudents
func commonStudents(res http.ResponseWriter, req *http.Request) {
	// Check that the request method is GET.
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Get the list of teachers from the query string.
	teachers := req.URL.Query()["teacher"]
	query, args, err := sqlx.In(`
		SELECT student_email
		FROM class NATURAL JOIN teachers NATURAL JOIN students
		WHERE teacher_email IN (?)
		GROUP BY student_email
		HAVING COUNT(*) = ?
	`, teachers, len(teachers))
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	}
	rows, err := db.Query(query, args...)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Convert the rows into a list of students.
	var students Students
	for rows.Next() {
		var student string
		err = rows.Scan(&student)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		students.Students = append(students.Students, student)
	}

	// Send the response.
	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(students)
	if err != nil {
		res.Header().Del("Content-Type")
		res.WriteHeader(http.StatusInternalServerError)
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
