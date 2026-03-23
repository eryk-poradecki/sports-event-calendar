package sport

import (
	"database/sql"
	"net/http"

	"github.com/eryk-poradecki/sports-event-calendar/internal/httpx"
)

func HandleGetAllSports(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sports, err := GetAll(db)
		if err != nil {
			httpx.WriteError(w, http.StatusInternalServerError, "failed to fetch sports", err)
			return
		}

		httpx.WriteJSON(w, http.StatusOK, sports)
	}
}
