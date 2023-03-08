package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"govtech-onecv-assessment/src/api"
	"govtech-onecv-assessment/src/database"
	"govtech-onecv-assessment/src/utils"
	"log"
	"net/http"
)

func registerAPIHandlers() {
	http.HandleFunc("/api/register", api.Register)
	http.HandleFunc("/api/commonstudents.go", api.CommonStudents)
	http.HandleFunc("/api/suspend", api.Suspend)
	http.HandleFunc("/api/retrievefornotifications", api.RetrieveForNotifications)
	http.HandleFunc("/api/reset", api.Reset)
}

func connectDatabase() {
	cfg := mysql.Config{
		User:   utils.GetEnvVariable("DATABASE_USER"),
		Passwd: utils.GetEnvVariable("DATABASE_PASSWORD"),
		Net:    "tcp",
		Addr:   utils.GetEnvVariable("DATABASE_URL"),
		DBName: utils.GetEnvVariable("DATABASE_NAME"),
	}
	// Get the database handle.
	var err error
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	database.SetDB(db)
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Printf("info: connected to database '%s'.\n", cfg.DBName)
}

func startServer() {
	// Start the server at the port specified in .env.
	port := utils.GetEnvVariable("PORT")
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
