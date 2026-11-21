package usecase

import (
	"context"
	"fmt"

	"github.com/isobekeito/safearea/internal/domain/model"
	"github.com/isobekeito/safearea/internal/domain/repository"
	"github.com/isobekeito/safearea/internal/dto/response"
)

// disasterHistoryUseCaseImpl はDisasterHistoryUseCaseの実装
type disasterHistoryUseCaseImpl struct {
	disasterRepo repository.DisasterRepository
}

// NewDisasterHistoryUseCase はDisasterHistoryUseCaseの新しいインスタンスを作成する
func NewDisasterHistoryUseCase(disasterRepo repository.DisasterRepository) DisasterHistoryUseCase {
	return &disasterHistoryUseCaseImpl{
		disasterRepo: disasterRepo,
	}
}

// GetDisasterHistory は住所から災害履歴を取得する
func (u *disasterHistoryUseCaseImpl) GetDisasterHistory(ctx context.Context, address string) (*response.DisasterHistoryResponse, error) {
	// 住所から都道府県・市区町村を抽出
	addrInfo := ParseAddress(address)
	if addrInfo.PrefectureCode == "" {
		return nil, fmt.Errorf("invalid address: prefecture not found")
	}

	var disasterModels []*model.Disaster
	var err error

	if addrInfo.City != "" {
		// 市区町村が特定できる場合は市区町村で検索
		disasterModels, err = u.disasterRepo.FindByCity(ctx, addrInfo.PrefectureCode, addrInfo.City)
	} else {
		// 都道府県のみの場合は都道府県で検索
		disasterModels, err = u.disasterRepo.FindByPrefecture(ctx, addrInfo.PrefectureCode)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find disasters: %w", err)
	}

	// レスポンスに変換
	disasters := make([]response.DisasterItem, 0, len(disasterModels))
	for _, d := range disasterModels {
		disasters = append(disasters, response.DisasterItem{
			ID:     d.ID,
			Type:   string(d.Type),
			Name:   d.Name,
			Year:   d.Year,
			Damage: d.Damage,
		})
	}

	return &response.DisasterHistoryResponse{
		Disasters: disasters,
	}, nil
}
