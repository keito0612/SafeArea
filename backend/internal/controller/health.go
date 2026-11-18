package controller

import (
	"encoding/json"
	"net/http"

	"github.com/isobekeito/safearea/internal/infrastructure/persistence"
)

type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	dbStatus := "ok"
	if err := persistence.HealthCheck(); err != nil {
		dbStatus = "error"
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(HealthResponse{
		Status:   "ok",
		Database: dbStatus,
	})
}
