package main

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

// POST /api/register
func register(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	res.WriteHeader(http.StatusNoContent)
}

// GET /api/commonstudents
func commonStudents(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	res.WriteHeader(http.StatusOK)
}

// POST /api/suspend
func suspend(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	res.WriteHeader(http.StatusNoContent)
}

// POST /api/retrievefornotifications
func retrieveForNotifications(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func getEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("error loading .env: %s\n", err)
	}
	return os.Getenv(key)
}

func main() {
	// Register the API handlers.
	http.HandleFunc("/api/register", register)
	http.HandleFunc("/api/commonstudents", commonStudents)
	http.HandleFunc("/api/suspend", suspend)
	http.HandleFunc("/api/retrievefornotifications", retrieveForNotifications)

	// Start the server at the port specified in .env.
	port := getEnvVariable("PORT")
	var addr string
	if port == "" {
		addr = ":3000"
	} else {
		addr = ":" + port
	}
	err := http.ListenAndServe(addr, nil)
	fmt.Printf("info: server listening on port %s\n", port)

	// Handle server initialization errors.
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("error: server closed.\n")
	} else if err != nil {
		fmt.Printf("error while starting server: %s\n", err)
	}
}
