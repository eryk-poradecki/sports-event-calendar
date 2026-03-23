package web

import (
	"database/sql"
	"errors"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/eryk-poradecki/sports-event-calendar/internal/event"
)

var indexTemplate = template.Must(template.ParseFiles("/web/templates/index.html"))
var eventDetailsTemplate = template.Must(template.ParseFiles("/web/templates/event_details.html"))

func RenderIndex(w http.ResponseWriter, r *http.Request) {
	if err := indexTemplate.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func RenderEventDetailsPage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "error parsing id", http.StatusBadRequest)
			return
		}
		eventDetails, err := event.GetByID(db, uint64(id))
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "event with given id not found", http.StatusNotFound)
				log.Printf("event with given id not found: %v", err)
				return
			}
			http.Error(w, "error getting event details", http.StatusInternalServerError)
			log.Printf("error getting event details: %v", err)
			return
		}

		if err := eventDetailsTemplate.Execute(w, eventDetails); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}
}
