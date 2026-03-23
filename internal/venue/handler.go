package venue

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func HandleGetAllVenues(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		venues, err := GetAll(db)
		if err != nil {
			http.Error(w, "failed to fetch venues", http.StatusInternalServerError)
			log.Printf("failed to fetch venues: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(venues); err != nil {
			http.Error(w, "failed to fetch venues", http.StatusInternalServerError)
			log.Printf("failed to fetch venues: %v", err)
			return
		}
	}
}
