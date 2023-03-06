package main

import (
	"errors"
	"fmt"
	"net/http"
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

func main() {
	http.HandleFunc("/api/register", register)
	http.HandleFunc("/api/commonstudents", commonStudents)
	http.HandleFunc("/api/suspend", suspend)
	http.HandleFunc("/api/retrievefornotifications", retrieveForNotifications)
	port := ":3333"
	err := http.ListenAndServe(port, nil)
	fmt.Printf("info: server listening on port %s\n", port)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("error: server closed.\n")
	} else if err != nil {
		fmt.Printf("error: %s\n", err)
	}
}
