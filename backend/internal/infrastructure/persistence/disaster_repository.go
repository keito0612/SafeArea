package persistence

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/isobekeito/safearea/internal/domain/model"
	"github.com/isobekeito/safearea/internal/domain/repository"
)

// disasterRepositoryImpl はDisasterRepositoryの実装
type disasterRepositoryImpl struct {
	db *gorm.DB
}

// NewDisasterRepository はDisasterRepositoryの新しいインスタンスを作成する
func NewDisasterRepository(db *gorm.DB) repository.DisasterRepository {
	return &disasterRepositoryImpl{db: db}
}

// FindByPrefecture は都道府県コードで災害を検索する
func (r *disasterRepositoryImpl) FindByPrefecture(ctx context.Context, prefCode string) ([]*model.Disaster, error) {
	var disasters []*model.Disaster
	if err := r.db.WithContext(ctx).
		Where("prefecture_code = ?", prefCode).
		Order("year DESC").
		Find(&disasters).Error; err != nil {
		return nil, fmt.Errorf("failed to find disasters by prefecture: %w", err)
	}
	return disasters, nil
}

// FindByCity は都道府県コードと市区町村名で災害を検索する
func (r *disasterRepositoryImpl) FindByCity(ctx context.Context, prefCode, cityName string) ([]*model.Disaster, error) {
	var disasters []*model.Disaster

	// disastersテーブルから直接検索
	if err := r.db.WithContext(ctx).
		Where("prefecture_code = ? AND city_name LIKE ?", prefCode, "%"+cityName+"%").
		Order("year DESC").
		Find(&disasters).Error; err != nil {
		return nil, fmt.Errorf("failed to find disasters by city: %w", err)
	}

	// disaster_areasテーブルからも検索して結合
	var areaDisasters []*model.Disaster
	if err := r.db.WithContext(ctx).
		Joins("JOIN disaster_areas ON disaster_areas.disaster_id = disasters.id").
		Where("disasters.prefecture_code = ? AND disaster_areas.city_name LIKE ?", prefCode, "%"+cityName+"%").
		Order("disasters.year DESC").
		Find(&areaDisasters).Error; err != nil {
		return nil, fmt.Errorf("failed to find disasters by city from areas: %w", err)
	}

	// 重複を排除して結合
	disasterMap := make(map[int64]*model.Disaster)
	for _, d := range disasters {
		disasterMap[d.ID] = d
	}
	for _, d := range areaDisasters {
		disasterMap[d.ID] = d
	}

	result := make([]*model.Disaster, 0, len(disasterMap))
	for _, d := range disasterMap {
		result = append(result, d)
	}

	return result, nil
}

// FindByID はIDで災害を取得する
func (r *disasterRepositoryImpl) FindByID(ctx context.Context, id int64) (*model.Disaster, error) {
	var disaster model.Disaster
	if err := r.db.WithContext(ctx).
		Preload("DisasterAreas").
		First(&disaster, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find disaster by id: %w", err)
	}
	return &disaster, nil
}
