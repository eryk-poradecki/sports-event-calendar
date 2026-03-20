package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/eryk-poradecki/sports-event-calendar/internal/database"
	"github.com/eryk-poradecki/sports-event-calendar/internal/event"
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

	router := http.NewServeMux()

	router.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "healthy")
	})

	router.HandleFunc("GET /events/{id}", event.HandleGetEventByID(db))
	router.HandleFunc("GET /events", event.HandleGetAllEvents(db))
	router.HandleFunc("POST /events", event.HandleCreateEvent(db))

	log.Printf("starting server on :%s", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	if err != nil {
		log.Fatalf("could not start the application: %v", err)
	}
}
