package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/pulsecity/services/team-service/internal/domain"
)

func RegisterHealthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", handleHealth)
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(domain.NewHealth()); err != nil {
		http.Error(w, "encode health response", http.StatusInternalServerError)
	}
}
