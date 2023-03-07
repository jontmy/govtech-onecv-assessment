package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func getEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("error loading .env: %s\n", err)
	}
	return os.Getenv(key)
}

func registerAPIHandlers() {
	http.HandleFunc("/api/register", register)
	http.HandleFunc("/api/commonstudents", commonStudents)
	http.HandleFunc("/api/suspend", suspend)
	http.HandleFunc("/api/retrievefornotifications", retrieveForNotifications)
}

func connectDatabase() {
	cfg := mysql.Config{
		User:   getEnvVariable("DATABASE_USER"),
		Passwd: getEnvVariable("DATABASE_PASSWORD"),
		Net:    "tcp",
		Addr:   getEnvVariable("DATABASE_URL"),
		DBName: getEnvVariable("DATABASE_NAME"),
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Printf("info: connected to database '%s'.\n", cfg.DBName)
}

func startServer() {
	// Start the server at the port specified in .env.
	port := getEnvVariable("PORT")
	var addr string
	if port == "" {
		addr = ":3000"
	} else {
		addr = ":" + port
	}
	err := http.ListenAndServe(addr, nil)
	fmt.Printf("info: server listening on port %s.\n", port)

	// Handle server initialization errors.
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("error: server closed.\n")
	} else if err != nil {
		fmt.Printf("error while starting server: %s.\n", err)
	}
}

func main() {
	registerAPIHandlers()
	connectDatabase()
	startServer()
}
