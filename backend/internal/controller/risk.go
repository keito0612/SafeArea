package controller

import (
	"encoding/json"
	"net/http"
)

type RiskResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// GetRisk は災害リスク検索エンドポイント
// GET /api/risk?address=xxx
func GetRisk(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	address := r.URL.Query().Get("address")
	if address == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RiskResponse{
			Success: false,
			Error: &ErrorInfo{
				Code:    "VALIDATION_ERROR",
				Message: "addressパラメータは必須です",
			},
		})
		return
	}

	// TODO: 実際のリスク検索ロジックを実装
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(RiskResponse{
		Success: true,
		Data: map[string]interface{}{
			"address":      address,
			"overall_risk": "medium",
			"message":      "Coming soon...",
		},
	})
}
