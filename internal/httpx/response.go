package httpx

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	var buf bytes.Buffer

	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(true)

	if err := encoder.Encode(payload); err != nil {
		log.Printf("failed to encode JSON response: %v", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		_, writeErr := w.Write([]byte(`{"error": "internal server error"}` + "\n"))
		if writeErr != nil {
			log.Printf("failed to write fallback JSON error response: %v", writeErr)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err := w.Write(buf.Bytes()); err != nil {
		log.Printf("failed to write JSON response: %v", err)
	}
}

func WriteError(w http.ResponseWriter, status int, message string, err error) {
	if err != nil {
		if status >= http.StatusInternalServerError {
			log.Printf("HTTP %d: %s: %v", status, message, err)
		} else {
			log.Printf("HTTP %d: %s", status, message)
		}
	}

	WriteJSON(w, status, struct {
		Error string `json:"error"`
	}{
		Error: message,
	})
}
