package persistence

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/isobekeito/safearea/internal/domain/model"
	"github.com/isobekeito/safearea/internal/domain/repository"
)

// earthquakeRiskRepositoryImpl はEarthquakeRiskRepositoryの実装
type earthquakeRiskRepositoryImpl struct {
	db *gorm.DB
}

// NewEarthquakeRiskRepository はEarthquakeRiskRepositoryの新しいインスタンスを作成する
func NewEarthquakeRiskRepository(db *gorm.DB) repository.EarthquakeRiskRepository {
	return &earthquakeRiskRepositoryImpl{db: db}
}

// FindByMeshCode はメッシュコードで地震危険度を取得する
func (r *earthquakeRiskRepositoryImpl) FindByMeshCode(ctx context.Context, meshCode string) (*model.EarthquakeRisk, error) {
	var risk model.EarthquakeRisk
	if err := r.db.WithContext(ctx).
		Where("mesh_code = ?", meshCode).
		First(&risk).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find earthquake risk by mesh code: %w", err)
	}
	return &risk, nil
}
