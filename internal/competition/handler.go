package competition

import (
	"database/sql"
	"net/http"

	"github.com/eryk-poradecki/sports-event-calendar/internal/httpx"
)

func HandleGetAllCompetitions(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		competitions, err := GetAll(db)
		if err != nil {
			httpx.WriteError(w, http.StatusInternalServerError, "failed to fetch competitions", err)
			return
		}

		httpx.WriteJSON(w, http.StatusOK, competitions)
	}
}
