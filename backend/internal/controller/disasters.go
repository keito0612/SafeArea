package controller

import (
	"encoding/json"
	"net/http"
)

type DisastersResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// GetDisasters は過去の災害履歴エンドポイント
// GET /api/disasters?address=xxx
func GetDisasters(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	address := r.URL.Query().Get("address")
	if address == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(DisastersResponse{
			Success: false,
			Error: &ErrorInfo{
				Code:    "VALIDATION_ERROR",
				Message: "addressパラメータは必須です",
			},
		})
		return
	}

	// TODO: 実際の災害履歴検索ロジックを実装
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(DisastersResponse{
		Success: true,
		Data: map[string]interface{}{
			"disasters": []interface{}{},
			"message":   "Coming soon...",
		},
	})
}
