package venue

import (
	"database/sql"
	"net/http"

	"github.com/eryk-poradecki/sports-event-calendar/internal/httpx"
)

func HandleGetAllVenues(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		venues, err := GetAll(db)
		if err != nil {
			httpx.WriteError(w, http.StatusInternalServerError, "failed to fetch venues", err)
			return
		}

		httpx.WriteJSON(w, http.StatusOK, venues)
	}
}
