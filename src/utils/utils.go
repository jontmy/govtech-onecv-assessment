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

func ParseJSON[T any](res http.ResponseWriter, req *http.Request, val *T) {
	// Check that the request body is JSON.
	if req.Header.Get("Content-Type") != "application/json" {
		res.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	// Parse the request body.
	err := json.NewDecoder(req.Body).Decode(val)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
	}
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
