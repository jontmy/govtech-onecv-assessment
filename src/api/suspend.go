package api

import (
	"govtech-onecv-assessment/src/database"
	"govtech-onecv-assessment/src/httputils"
	"net/http"
)

type Student struct {
	Student string `json:"student"`
}

// Suspend POST /api/suspend
func Suspend(res http.ResponseWriter, req *http.Request) {
	db := database.GetDB()

	// Check that the request method is POST.
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST is allowed.", http.StatusMethodNotAllowed)
		return
	}

	// Check that the student exists, otherwise return a 404.
	var student Student
	var studentExists bool
	httputils.ParseJSON(res, req, &student)
	row := db.QueryRow("SELECT EXISTS (SELECT 1 FROM students WHERE student_email = ?)", student.Student)
	err := row.Scan(&studentExists)
	if err != nil {
		httputils.HandleServerError(res, err)
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
		httputils.HandleServerError(res, err)
		return
	}
	res.WriteHeader(http.StatusNoContent)
}
