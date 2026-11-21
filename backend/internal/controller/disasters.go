package controller

import (
	"net/http"
	"strings"

	"github.com/isobekeito/safearea/internal/dto/response"
	"github.com/isobekeito/safearea/internal/usecase"
)

// DisastersController は災害履歴APIのコントローラー
type DisastersController struct {
	disasterHistoryUseCase usecase.DisasterHistoryUseCase
}

// NewDisastersController はDisastersControllerの新しいインスタンスを作成する
func NewDisastersController(disasterHistoryUseCase usecase.DisasterHistoryUseCase) *DisastersController {
	return &DisastersController{
		disasterHistoryUseCase: disasterHistoryUseCase,
	}
}

// GetDisasters は過去の災害履歴エンドポイント
// GET /api/search/disasters?address=xxx
func (c *DisastersController) GetDisasters(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	address := strings.TrimSpace(r.URL.Query().Get("address"))
	if address == "" {
		response.WriteValidationError(w, []response.ErrorDetail{
			{Field: "address", Message: "住所は必須です"},
		})
		return
	}

	result, err := c.disasterHistoryUseCase.GetDisasterHistory(ctx, address)
	if err != nil {
		if strings.Contains(err.Error(), "prefecture not found") {
			response.WriteValidationError(w, []response.ErrorDetail{
				{Field: "address", Message: "有効な都道府県が含まれていません"},
			})
			return
		}
		response.WriteInternalError(w)
		return
	}

	response.WriteSuccess(w, result)
}
