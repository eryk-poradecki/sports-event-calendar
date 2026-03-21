package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/eryk-poradecki/sports-event-calendar/internal/database"
	"github.com/eryk-poradecki/sports-event-calendar/internal/event"
)

var indexTemplate = template.Must(template.ParseFiles("/web/templates/index.html"))

func renderIndex(w http.ResponseWriter, r *http.Request) {
	if err := indexTemplate.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

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

	apiV1 := http.NewServeMux()
	apiV1.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "healthy")
	})
	apiV1.HandleFunc("GET /events/{id}", event.HandleGetEventByID(db))
	apiV1.HandleFunc("GET /events", event.HandleGetAllEvents(db))
	apiV1.HandleFunc("POST /events", event.HandleCreateEvent(db))

	router.Handle("/api/v1", http.StripPrefix("/api/v1", apiV1))

	staticFiles := http.FileServer(http.Dir("web/static"))
	router.Handle("GET /static/", http.StripPrefix("/static/", staticFiles))
	router.HandleFunc("GET /", renderIndex)

	log.Printf("starting server on :%s", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	if err != nil {
		log.Fatalf("could not start the application: %v", err)
	}
}
