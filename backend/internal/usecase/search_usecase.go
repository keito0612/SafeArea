//go:generate mockgen -source=$GOFILE -destination=../mock/mock_search_usecase.go -package=mock

package usecase

import (
	"context"

	"github.com/isobekeito/safearea/internal/dto/response"
)

// SearchUseCase は災害リスク検索のユースケースインターフェース
type SearchUseCase interface {
	// Search は住所から災害リスクを検索する
	Search(ctx context.Context, address string) (*response.SearchResponse, error)
}
