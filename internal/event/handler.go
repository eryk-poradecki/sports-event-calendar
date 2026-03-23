package event

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/eryk-poradecki/sports-event-calendar/internal/httpx"
)

func HandleGetEventByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "invalid event ID", err)
			return
		}

		ev, err := GetByID(db, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				httpx.WriteError(w, http.StatusNotFound, "event not found", err)
			} else {
				httpx.WriteError(w, http.StatusInternalServerError, "error getting event", err)
			}
			return
		}

		httpx.WriteJSON(w, http.StatusOK, ev)
	}
}

func HandleGetAllEvents(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.URL.Query().Get("page")
		pageSizeStr := r.URL.Query().Get("page_size")
		sport := r.URL.Query().Get("sport")
		dateFrom := r.URL.Query().Get("date_from")
		dateTo := r.URL.Query().Get("date_to")
		if pageStr == "" {
			pageStr = "1"
		}
		pageInt, err := strconv.Atoi(pageStr)
		if err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "invalid page", err)
			return
		}
		if pageSizeStr == "" {
			pageSizeStr = "10"
		}
		pageSizeInt, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "invalid page size", err)
			return
		}
		events, err := GetAllEvents(db, pageInt, pageSizeInt, sport, dateFrom, dateTo)
		if err != nil {
			if errors.Is(err, ErrSportNotFound) {
				httpx.WriteError(w, http.StatusNotFound, "sport not found", err)
				return
			}
			if errors.Is(err, ErrInvalidDate) {
				httpx.WriteError(w, http.StatusBadRequest, err.Error(), err)
				return
			}
			httpx.WriteError(w, http.StatusInternalServerError, "failed to fetch events", err)
			return
		}

		httpx.WriteJSON(w, http.StatusOK, events)
	}
}

func HandleCreateEvent(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var ev Event
		if err := json.NewDecoder(r.Body).Decode(&ev); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "failed to decode request body", err)
			return
		}

		if err := CreateEvent(db, &ev); err != nil {
			if errors.Is(err, ErrInvalidEvent) {
				httpx.WriteError(w, http.StatusBadRequest, err.Error(), err)
				return
			}
			httpx.WriteError(w, http.StatusInternalServerError, "failed to create event", err)
			return
		}

		httpx.WriteJSON(w, http.StatusCreated, ev)
	}
}
