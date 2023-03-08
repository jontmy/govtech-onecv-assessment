package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"govtech-onecv-assessment/src/database"
	"govtech-onecv-assessment/src/httputils"
	"net/http"
	"strings"
)

type Notification struct {
	Teacher      string `json:"teacher"`
	Notification string `json:"notification"`
}
type Recipients struct {
	Recipients []string `json:"recipients"`
}

// RetrieveForNotifications Implements POST /api/retrievefornotifications.
func RetrieveForNotifications(res http.ResponseWriter, req *http.Request) {
	db := database.GetDB()

	// Check that the request method is POST.
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST is allowed.", http.StatusMethodNotAllowed)
		return
	}
	// Deserialize the request JSON.
	var notification Notification
	httputils.ParseJSON(res, req, &notification)
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
			httputils.HandleServerError(res, err)
			return
		}
		rows, err = db.Query(query, args...)
	}
	if err != nil {
		httputils.HandleServerError(res, err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			httputils.HandleServerError(res, err)
			return
		}
	}(rows)

	// Convert the rows into a list of students.
	var recipients Recipients
	for rows.Next() {
		var student string
		err = rows.Scan(&student)
		if err != nil {
			httputils.HandleServerError(res, err)
			return
		}
		recipients.Recipients = append(recipients.Recipients, student)
	}

	// Send the response.
	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(recipients)
	if err != nil {
		httputils.HandleServerError(res, err)
	}
	res.WriteHeader(http.StatusOK)
}
