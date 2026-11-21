//go:generate mockgen -source=$GOFILE -destination=../mock/mock_disaster_history_usecase.go -package=mock

package usecase

import (
	"context"

	"github.com/isobekeito/safearea/internal/dto/response"
)

// DisasterHistoryUseCase は災害履歴検索のユースケースインターフェース
type DisasterHistoryUseCase interface {
	// GetDisasterHistory は住所から災害履歴を取得する
	GetDisasterHistory(ctx context.Context, address string) (*response.DisasterHistoryResponse, error)
}
