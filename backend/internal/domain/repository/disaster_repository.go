//go:generate mockgen -source=$GOFILE -destination=../../mock/mock_disaster_repository.go -package=mock

package repository

import (
	"context"

	"github.com/isobekeito/safearea/internal/domain/model"
)

// DisasterRepository は災害データへのアクセスを提供するインターフェース
type DisasterRepository interface {
	// FindByPrefecture は都道府県コードで災害を検索する
	FindByPrefecture(ctx context.Context, prefCode string) ([]*model.Disaster, error)

	// FindByCity は都道府県コードと市区町村名で災害を検索する
	FindByCity(ctx context.Context, prefCode, cityName string) ([]*model.Disaster, error)

	// FindByID はIDで災害を取得する
	FindByID(ctx context.Context, id int64) (*model.Disaster, error)
}
