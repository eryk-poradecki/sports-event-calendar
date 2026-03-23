package web

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/eryk-poradecki/sports-event-calendar/internal/competition"
	"github.com/eryk-poradecki/sports-event-calendar/internal/event"
	"github.com/eryk-poradecki/sports-event-calendar/internal/sport"
	"github.com/eryk-poradecki/sports-event-calendar/internal/team"
	"github.com/eryk-poradecki/sports-event-calendar/internal/venue"
)

func NewRouter(db *sql.DB) http.Handler {
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
	apiV1.HandleFunc("GET /venues", venue.HandleGetAllVenues(db))

	router.Handle("/api/v1/", http.StripPrefix("/api/v1", apiV1))

	staticFiles := http.FileServer(http.Dir("web/static"))
	router.Handle("GET /static/", http.StripPrefix("/static/", staticFiles))
	router.HandleFunc("GET /{$}", RenderIndex)
	router.HandleFunc("GET /events/{id}", RenderEventDetailsPage(db))

	return router
}
