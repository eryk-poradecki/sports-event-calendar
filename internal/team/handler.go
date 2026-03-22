package team

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func HandleGetAllTeams(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teams, err := GetAll(db)
		if err != nil {
			http.Error(w, "failed to fetch teams", http.StatusInternalServerError)
			log.Printf("failed to fetch teams: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(teams); err != nil {
			http.Error(w, "failed to encode teams", http.StatusInternalServerError)
			log.Printf("failed to encode teams: %v", err)
			return
		}
	}
}
