package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/eryk-poradecki/sports-event-calendar/internal/database"
	"github.com/eryk-poradecki/sports-event-calendar/internal/web"
)

func main() {
	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8080"
	}
	connectionString := os.Getenv("DB_CONNECTION_STRING")
	if connectionString == "" {
		log.Fatal("DB_CONNECTION_STRING environment variable not set")
	}

	db, err := database.ConnectDatabase(connectionString)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	defer db.Close()

	log.Printf("starting server on :%s", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), web.NewRouter(db))
	if err != nil {
		log.Fatalf("could not start the application: %v", err)
	}
}
