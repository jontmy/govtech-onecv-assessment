package utils

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

func GetEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("error loading .env: %s\n", err)
	}
	return os.Getenv(key)
}

// ParseJSON ensures that the request body and header specifies JSON, and parses it into the given value.
// Returns true if the request body is parsed successfully, false otherwise.
func ParseJSON[T any](res http.ResponseWriter, req *http.Request, val *T) bool {
	// Check that the request body is JSON.
	if req.Header.Get("Content-Type") != "application/json" {
		HandleCustomError(res, "Content-Type must be application/json.", http.StatusUnsupportedMediaType)
		return false
	}
	// Check JSON body is not empty.
	if req.ContentLength == 0 {
		HandleCustomError(res, "Empty request body.", http.StatusBadRequest)
		return false
	}
	// Parse the request body.
	err := json.NewDecoder(req.Body).Decode(val)
	if err != nil {
		HandleCustomError(res, "Invalid JSON body.\n"+err.Error(), http.StatusBadRequest)
		return false
	}
	return true
}

type ErrorMessage struct {
	Message string `json:"message"`
}

// HandleServerError handles an error by logging it and sending a 500 response.
func HandleServerError(res http.ResponseWriter, err error) {
	if err == nil {
		return
	}
	fmt.Println(err)
	// Error message should be encoded in JSON.
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusInternalServerError)
	msg := ErrorMessage{Message: err.Error()}
	err = json.NewEncoder(res).Encode(msg)
	if err != nil {
		// Not much else we can do here that won't also need error handling.
		return
	}
}

// HandleCustomError handles an error by sending a custom error message in the JSON body with a specified HTTP status code.
func HandleCustomError(res http.ResponseWriter, msg string, code int) {
	// Error message should be encoded in JSON.
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(code)
	errMsg := ErrorMessage{Message: msg}
	err := json.NewEncoder(res).Encode(errMsg)
	if err != nil {
		// Not much else we can do here that won't also need error handling.
		return
	}
}
