package sport

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func HandleGetAllSports(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sports, err := GetAll(db)
		if err != nil {
			http.Error(w, "failed to fetch sports", http.StatusInternalServerError)
			log.Printf("get sports failed: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(sports); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("response encoding failed: %v", err)
			return
		}
	}
}
