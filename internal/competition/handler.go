package competition

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func HandleGetAllCompetitions(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		competitions, err := GetAll(db)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "competitions not found", http.StatusNotFound)
				log.Printf("failed to fetch competitions: %v", err)
				return
			}
			http.Error(w, "failed to fetch competitions", http.StatusInternalServerError)
			log.Printf("failed to fetch competitions: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(competitions); err != nil {
			http.Error(w, "failed to encode competitions", http.StatusInternalServerError)
			log.Printf("failed to encode competitions: %v", err)
			return
		}
	}
}
