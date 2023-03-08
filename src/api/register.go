package api

import (
	"govtech-onecv-assessment/src/database"
	"govtech-onecv-assessment/src/utils"
	"net/http"
)

type Registration struct {
	Teacher  string   `json:"teacher"`
	Students []string `json:"students"`
}

// Register Implements POST /api/register.
func Register(res http.ResponseWriter, req *http.Request) {
	db := database.GetDB()

	// Check that the request method is POST.
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST is allowed.", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body.
	var registration Registration
	utils.ParseJSON(res, req, &registration)

	// Insert the teacher into the database if they don't already exist.
	_, err := db.Exec("INSERT IGNORE INTO teachers (teacher_email) VALUES (?)", registration.Teacher)
	if err != nil {
		utils.HandleServerError(res, err)
		return
	}

	// Insert the students into the database if they don't already exist.
	for _, student := range registration.Students {
		_, err = db.Exec(`
			INSERT IGNORE INTO students (student_email) VALUES (?);
		`, student)
		if err != nil {
			utils.HandleServerError(res, err)
			return
		}
		_, err = db.Exec(`
			INSERT IGNORE INTO class (teacher_email, student_email) VALUES (?, ?);
		`, registration.Teacher, student)
		if err != nil {
			utils.HandleServerError(res, err)
			return
		}
	}
	// Send the response.
	res.WriteHeader(http.StatusNoContent)
}
