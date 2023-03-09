package api

import (
	"database/sql"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"govtech-onecv-assessment/src/database"
	"govtech-onecv-assessment/src/utils"
	"net/http"
)

type Students struct {
	Students []string `json:"students"`
}

// CommonStudents implements GET /api/commonstudents.
func CommonStudents(res http.ResponseWriter, req *http.Request) {
	db := database.GetDB()

	// Check that the request method is GET.
	if req.Method != http.MethodGet {
		utils.HandleCustomError(res, "Only GET is allowed.", http.StatusMethodNotAllowed)
		return
	}

	// Get the list of teachers from the query string.
	teachers := req.URL.Query()["teacher"]
	var query string
	var args []interface{}
	var err error
	if len(teachers) == 0 {
		// If there are no teachers specified, return all students.
		query = `
			SELECT student_email
			FROM students
		`
	} else {
		query, args, err = sqlx.In(`
		SELECT DISTINCT student_email
		FROM class NATURAL JOIN teachers NATURAL JOIN students
		WHERE teacher_email IN (?)
		GROUP BY student_email
		HAVING COUNT(*) = ?
	`, teachers, len(teachers))
	}

	if err != nil {
		utils.HandleServerError(res, err)
		return
	}
	rows, err := db.Query(query, args...)
	if err != nil {
		utils.HandleServerError(res, err)
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			utils.HandleServerError(res, err)
			return
		}
	}(rows)

	// Convert the rows into a list of students.
	// If there are no students, return an empty list.
	students := Students{Students: []string{}}
	for rows.Next() {
		var student string
		err = rows.Scan(&student)
		if err != nil {
			utils.HandleServerError(res, err)
			return
		}
		students.Students = append(students.Students, student)
	}

	// Send the response.
	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(students)
	if err != nil {
		utils.HandleServerError(res, err)
		return
	}
	res.WriteHeader(http.StatusOK)
}
