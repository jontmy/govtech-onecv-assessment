package api

import (
	"govtech-onecv-assessment/src/database"
	"govtech-onecv-assessment/src/utils"
	"net/http"
)

// Reset Implements DELETE /api/reset.
func Reset(res http.ResponseWriter, req *http.Request) {
	db := database.GetDB()

	// Check that the request method is DELETE.
	if req.Method != http.MethodDelete {
		http.Error(res, "Only DELETE is allowed.", http.StatusMethodNotAllowed)
		return
	}
	// Clear all values from the database.
	if _, err := db.Exec(`TRUNCATE TABLE class;`); err != nil {
		utils.HandleServerError(res, err)
		return
	}
	if _, err := db.Exec(`TRUNCATE TABLE teachers;`); err != nil {
		utils.HandleServerError(res, err)
		return
	}
	if _, err := db.Exec(`TRUNCATE TABLE students;`); err != nil {
		utils.HandleServerError(res, err)
		return
	}
	res.WriteHeader(http.StatusNoContent)
}
