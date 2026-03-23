package team

import (
	"database/sql"
	"net/http"

	"github.com/eryk-poradecki/sports-event-calendar/internal/httpx"
)

func HandleGetAllTeams(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teams, err := GetAll(db)
		if err != nil {
			httpx.WriteError(w, http.StatusInternalServerError, "failed to fetch teams", err)
			return
		}

		httpx.WriteJSON(w, http.StatusOK, teams)
	}
}
