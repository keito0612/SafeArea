package model

import (
	"time"
)

// RiskLevel はリスクレベルを表す
type RiskLevel string

const (
	RiskLevelHigh   RiskLevel = "high"
	RiskLevelMedium RiskLevel = "medium"
	RiskLevelLow    RiskLevel = "low"
)

// EarthquakeRisk は地震危険度データを表すエンティティ（メッシュ単位）
type EarthquakeRisk struct {
	MeshCode          string    `gorm:"type:varchar(20);primaryKey" json:"mesh_code"`
	RiskLevel         RiskLevel `gorm:"type:varchar(10);not null" json:"risk_level"`
	Probability30Year float64   `gorm:"type:decimal(5,2)" json:"probability_30year"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName はGORMのテーブル名を指定する
func (EarthquakeRisk) TableName() string {
	return "earthquake_risks"
}
