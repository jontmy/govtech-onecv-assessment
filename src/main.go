package main

import (
	"errors"
	"fmt"
	"net/http"
)

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
