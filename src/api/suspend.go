package api

import (
	"govtech-onecv-assessment/src/database"
	"govtech-onecv-assessment/src/utils"
	"net/http"
)

type Student struct {
	Student string `json:"student"`
}

// Suspend implements POST /api/suspend
func Suspend(res http.ResponseWriter, req *http.Request) {
	db := database.GetDB()

	// Check that the request method is POST.
	if req.Method != http.MethodPost {
		utils.HandleCustomError(res, "Only POST is allowed.", http.StatusMethodNotAllowed)
		return
	}

	// Check that the student exists, otherwise return a 422 if unspecified, or 404 otherwise.
	var student Student
	var studentExists bool
	if !utils.ParseJSON(res, req, &student) {
		return
	}
	row := db.QueryRow("SELECT EXISTS (SELECT 1 FROM students WHERE student_email = ?)", student.Student)
	err := row.Scan(&studentExists)
	if err != nil {
		utils.HandleServerError(res, err)
		return
	}
	if student.Student == "" {
		utils.HandleCustomError(res, "Student email must be specified.", http.StatusUnprocessableEntity)
		return
	}
	if !studentExists {
		utils.HandleCustomError(res, "Student does not exist.", http.StatusNotFound)
		return
	}

	// Suspend the student.
	_, err = db.Exec(`
		UPDATE students
		SET is_suspended = TRUE
		WHERE student_email = ?
	`, student.Student)
	if err != nil {
		utils.HandleServerError(res, err)
		return
	}
	res.WriteHeader(http.StatusNoContent)
}
