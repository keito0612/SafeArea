package repository

import (
	"context"

	"github.com/isobekeito/safearea/internal/domain/model"
)

// EarthquakeRiskRepository は地震危険度データへのアクセスを提供するインターフェース
type EarthquakeRiskRepository interface {
	// FindByMeshCode はメッシュコードで地震危険度を取得する
	FindByMeshCode(ctx context.Context, meshCode string) (*model.EarthquakeRisk, error)
}
