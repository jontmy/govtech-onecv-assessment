package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strings"
)

var db *sql.DB

// POST /api/register
func register(res http.ResponseWriter, req *http.Request) {
	type Registration struct {
		Teacher  string   `json:"teacher"`
		Students []string `json:"students"`
	}

	// Check that the request method is POST.
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST is allowed.", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body.
	var registration Registration
	parseJSON(res, req, &registration)

	// Insert the teacher into the database if they don't already exist.
	_, err := db.Exec("INSERT IGNORE INTO teachers (teacher_email) VALUES (?)", registration.Teacher)
	if err != nil {
		handleServerError(res, err)
		return
	}

	// Insert the students into the database if they don't already exist.
	for _, student := range registration.Students {
		_, err = db.Exec(`
			INSERT IGNORE INTO students (student_email) VALUES (?);
		`, student)
		if err != nil {
			handleServerError(res, err)
			return
		}
		_, err = db.Exec(`
			INSERT IGNORE INTO class (teacher_email, student_email) VALUES (?, ?);
		`, registration.Teacher, student)
		if err != nil {
			handleServerError(res, err)
			return
		}
	}
	// Send the response.
	res.WriteHeader(http.StatusNoContent)
}

// GET /api/commonstudents
func commonStudents(res http.ResponseWriter, req *http.Request) {
	type Students struct {
		Students []string `json:"students"`
	}

	// Check that the request method is GET.
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET is allowed.", http.StatusMethodNotAllowed)
		return
	}

	// Get the list of teachers from the query string.
	teachers := req.URL.Query()["teacher"]
	query, args, err := sqlx.In(`
		SELECT DISTINCT student_email
		FROM class NATURAL JOIN teachers NATURAL JOIN students
		WHERE teacher_email IN (?)
		GROUP BY student_email
		HAVING COUNT(*) = ?
	`, teachers, len(teachers))
	if err != nil {
		handleServerError(res, err)
		return
	}
	rows, err := db.Query(query, args...)
	if err != nil {
		handleServerError(res, err)
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			handleServerError(res, err)
			return
		}
	}(rows)

	// Convert the rows into a list of students.
	var students Students
	for rows.Next() {
		var student string
		err = rows.Scan(&student)
		if err != nil {
			handleServerError(res, err)
			return
		}
		students.Students = append(students.Students, student)
	}

	// Send the response.
	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(students)
	if err != nil {
		handleServerError(res, err)
		return
	}
	res.WriteHeader(http.StatusOK)
}

// POST /api/suspend
func suspend(res http.ResponseWriter, req *http.Request) {
	type Student struct {
		Student string `json:"student"`
	}

	// Check that the request method is POST.
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST is allowed.", http.StatusMethodNotAllowed)
		return
	}

	// Check that the student exists, otherwise return a 404.
	var student Student
	var studentExists bool
	parseJSON(res, req, &student)
	row := db.QueryRow("SELECT EXISTS (SELECT 1 FROM students WHERE student_email = ?)", student.Student)
	err := row.Scan(&studentExists)
	if err != nil {
		handleServerError(res, err)
		return
	}
	if !studentExists {
		http.Error(res, "Student not found.", http.StatusNotFound)
		return
	}

	// Suspend the student.
	_, err = db.Exec(`
		UPDATE students
		SET is_suspended = TRUE
		WHERE student_email = ?
	`, student.Student)
	if err != nil {
		handleServerError(res, err)
		return
	}
	res.WriteHeader(http.StatusNoContent)
}

// POST /api/retrievefornotifications
func retrieveForNotifications(res http.ResponseWriter, req *http.Request) {
	type Notification struct {
		Teacher      string `json:"teacher"`
		Notification string `json:"notification"`
	}
	type Recipients struct {
		Recipients []string `json:"recipients"`
	}

	// Check that the request method is POST.
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST is allowed.", http.StatusMethodNotAllowed)
		return
	}
	// Deserialize the request JSON.
	var notification Notification
	parseJSON(res, req, &notification)
	// Get the message and the list of all @-mentioned students.
	fmt.Println(notification.Notification)
	splits := strings.Split(notification.Notification, " @")
	students := splits[1:]
	fmt.Println(students)
	fmt.Println(notification.Teacher)

	// Return the list of students who can receive a given notification.
	var rows *sql.Rows
	var err error

	// IN () is not a valid expression, so we need to handle the case where there are no @-mentioned students.
	if len(students) == 0 {
		rows, err = db.Query(`
		SELECT DISTINCT c.student_email
		FROM class c, students s
		WHERE c.student_email = s.student_email
			AND s.is_suspended = FALSE
			AND c.teacher_email = ?;
		`, notification.Teacher)
	} else {
		query, args, err := sqlx.In(`
		SELECT DISTINCT c.student_email
		FROM class c, students s
		WHERE c.student_email = s.student_email
			AND s.is_suspended = FALSE
			AND (c.teacher_email = ? OR c.student_email IN (?));
		`, notification.Teacher, students)
		if err != nil {
			handleServerError(res, err)
			return
		}
		rows, err = db.Query(query, args...)
	}
	if err != nil {
		handleServerError(res, err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			handleServerError(res, err)
			return
		}
	}(rows)

	// Convert the rows into a list of students.
	var recipients Recipients
	for rows.Next() {
		var student string
		err = rows.Scan(&student)
		if err != nil {
			handleServerError(res, err)
			return
		}
		recipients.Recipients = append(recipients.Recipients, student)
	}

	// Send the response.
	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(recipients)
	if err != nil {
		handleServerError(res, err)
	}
	res.WriteHeader(http.StatusOK)
}

// DELETE /api/reset
func reset(res http.ResponseWriter, req *http.Request) {
	// Check that the request method is DELETE.
	if req.Method != http.MethodDelete {
		http.Error(res, "Only DELETE is allowed.", http.StatusMethodNotAllowed)
		return
	}
	// Clear all values from the database.
	if _, err := db.Exec(`TRUNCATE TABLE class;`); err != nil {
		handleServerError(res, err)
		return
	}
	if _, err := db.Exec(`TRUNCATE TABLE teachers;`); err != nil {
		handleServerError(res, err)
		return
	}
	if _, err := db.Exec(`TRUNCATE TABLE students;`); err != nil {
		handleServerError(res, err)
		return
	}
	res.WriteHeader(http.StatusNoContent)
}
