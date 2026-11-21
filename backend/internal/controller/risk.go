package controller

import (
	"net/http"
	"strings"

	"github.com/isobekeito/safearea/internal/dto/response"
	"github.com/isobekeito/safearea/internal/usecase"
)

// SearchController は検索APIのコントローラー
type SearchController struct {
	searchUseCase usecase.SearchUseCase
}

// NewSearchController はSearchControllerの新しいインスタンスを作成する
func NewSearchController(searchUseCase usecase.SearchUseCase) *SearchController {
	return &SearchController{
		searchUseCase: searchUseCase,
	}
}

// Search は災害リスク検索エンドポイント
// GET /api/search?address=xxx
func (c *SearchController) Search(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	address := strings.TrimSpace(r.URL.Query().Get("address"))
	if address == "" {
		response.WriteValidationError(w, []response.ErrorDetail{
			{Field: "address", Message: "住所は必須です"},
		})
		return
	}

	result, err := c.searchUseCase.Search(ctx, address)
	if err != nil {
		if strings.Contains(err.Error(), "geocode") || strings.Contains(err.Error(), "no results") {
			response.WriteAddressNotFoundError(w)
			return
		}
		response.WriteInternalError(w)
		return
	}

	response.WriteSuccess(w, result)
}
