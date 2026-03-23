package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/eryk-poradecki/sports-event-calendar/internal/competition"
	"github.com/eryk-poradecki/sports-event-calendar/internal/database"
	"github.com/eryk-poradecki/sports-event-calendar/internal/event"
	"github.com/eryk-poradecki/sports-event-calendar/internal/sport"
	"github.com/eryk-poradecki/sports-event-calendar/internal/team"
)

var indexTemplate = template.Must(template.ParseFiles("/web/templates/index.html"))
var eventDetailsTemplate = template.Must(template.ParseFiles("/web/templates/event_details.html"))

func renderIndex(w http.ResponseWriter, r *http.Request) {
	if err := indexTemplate.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func renderEventDetailsPage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "error parsing id", http.StatusBadRequest)
			return
		}
		eventDetails, err := event.GetByID(db, uint64(id))
		if err != nil {
			http.Error(w, "error getting event details", http.StatusNotFound)
			return
		}

		if err := eventDetailsTemplate.Execute(w, eventDetails); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}
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

	apiV1.HandleFunc("GET /sports", sport.HandleGetAllSports(db))
	apiV1.HandleFunc("GET /teams", team.HandleGetAllTeams(db))
	apiV1.HandleFunc("GET /competitions", competition.HandleGetAllCompetitions(db))

	router.Handle("/api/v1/", http.StripPrefix("/api/v1", apiV1))

	staticFiles := http.FileServer(http.Dir("web/static"))
	router.Handle("GET /static/", http.StripPrefix("/static/", staticFiles))
	router.HandleFunc("GET /{$}", renderIndex)
	router.HandleFunc("GET /events/{id}", renderEventDetailsPage(db))

	log.Printf("starting server on :%s", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	if err != nil {
		log.Fatalf("could not start the application: %v", err)
	}
}
