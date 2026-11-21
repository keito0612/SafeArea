package usecase

import (
	"context"
	"fmt"

	"github.com/isobekeito/safearea/internal/domain/repository"
	"github.com/isobekeito/safearea/internal/dto/response"
	"github.com/isobekeito/safearea/internal/infrastructure/external"
)

// searchUseCaseImpl はSearchUseCaseの実装
type searchUseCaseImpl struct {
	geocodingClient external.GeocodingClient
	disasterRepo    repository.DisasterRepository
}

// NewSearchUseCase はSearchUseCaseの新しいインスタンスを作成する
func NewSearchUseCase(
	geocodingClient external.GeocodingClient,
	disasterRepo repository.DisasterRepository,
) SearchUseCase {
	return &searchUseCaseImpl{
		geocodingClient: geocodingClient,
		disasterRepo:    disasterRepo,
	}
}

// Search は住所から災害リスクを検索する
func (u *searchUseCaseImpl) Search(ctx context.Context, address string) (*response.SearchResponse, error) {
	// 1. ジオコーディング（住所→緯度経度）
	geocodeResult, err := u.geocodingClient.Geocode(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("failed to geocode address: %w", err)
	}

	// 2. 住所から都道府県・市区町村を抽出
	addrInfo := ParseAddress(address)

	// 3. 災害履歴件数を取得
	var disasterCount int
	if addrInfo.PrefectureCode != "" {
		var disasters []*struct{}
		if addrInfo.City != "" {
			disasterModels, err := u.disasterRepo.FindByCity(ctx, addrInfo.PrefectureCode, addrInfo.City)
			if err == nil {
				disasterCount = len(disasterModels)
			}
			_ = disasters // unused
		} else {
			disasterModels, err := u.disasterRepo.FindByPrefecture(ctx, addrInfo.PrefectureCode)
			if err == nil {
				disasterCount = len(disasterModels)
			}
		}
	}

	// 4. リスク情報を構築（現時点ではプレースホルダー）
	// TODO: ハザードマップAPI、J-SHIS APIからリスク情報を取得
	risks := []response.RiskItem{}
	tags := []string{}
	overallRisk := "unknown"

	return &response.SearchResponse{
		Address: address,
		Coordinates: response.Coordinates{
			Lat: geocodeResult.Coordinates.Lat,
			Lng: geocodeResult.Coordinates.Lng,
		},
		OverallRisk:          overallRisk,
		Risks:                risks,
		Tags:                 tags,
		DisasterHistoryCount: disasterCount,
	}, nil
}
