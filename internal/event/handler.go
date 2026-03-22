package event

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

func HandleGetEventByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid event ID", http.StatusBadRequest)
			log.Printf("get event failed: %v", err)
			return
		}

		ev, err := GetByID(db, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "event not found", http.StatusNotFound)
			} else {
				http.Error(w, "error getting event", http.StatusInternalServerError)
			}
			log.Printf("get event failed: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(ev); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			log.Printf("response encoding failed: %v", err)
			return
		}
	}
}

func HandleGetAllEvents(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.URL.Query().Get("page")
		pageSizeStr := r.URL.Query().Get("page_size")
		if pageStr == "" {
			pageStr = "0"
		}
		pageInt, err := strconv.Atoi(pageStr)
		if err != nil {
			http.Error(w, "invalid page", http.StatusBadRequest)
			return
		}
		if pageSizeStr == "" {
			pageSizeStr = "0"
		}
		pageSizeInt, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			http.Error(w, "invalid page size", http.StatusBadRequest)
			return
		}
		events, err := GetAllEvents(db, pageInt, pageSizeInt)
		if err != nil {
			http.Error(w, "failed to fetch events", http.StatusInternalServerError)
			log.Printf("get events failed: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(events); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			log.Printf("response encoding failed: %v", err)
			return
		}
	}
}

func HandleCreateEvent(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var ev Event
		if err := json.NewDecoder(r.Body).Decode(&ev); err != nil {
			http.Error(w, "failed to decode request body", http.StatusBadRequest)
			log.Printf("body decode failed: %v", err)
			return
		}

		if err := CreateEvent(db, &ev); err != nil {
			http.Error(w, "failed to create event", http.StatusBadRequest)
			log.Printf("create event failed: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(ev); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			log.Printf("response encoding failed: %v", err)
			return
		}
	}
}
