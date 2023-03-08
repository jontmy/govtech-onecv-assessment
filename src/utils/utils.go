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

func HandleServerError(res http.ResponseWriter, err error) {
	type ErrorMessage struct {
		Message string `json:"message"`
	}
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
