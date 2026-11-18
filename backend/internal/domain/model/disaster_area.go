package model

import (
	"time"

	"github.com/google/uuid"
)

// DisasterArea は災害の被災地域を表すエンティティ
type DisasterArea struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	DisasterID    uuid.UUID `gorm:"type:uuid;not null;index" json:"disaster_id"`
	CityName      string    `gorm:"type:varchar(255);not null;index" json:"city_name"`
	DamageSummary string    `gorm:"type:text" json:"damage_summary"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	Disaster Disaster `gorm:"foreignKey:DisasterID;constraint:OnDelete:CASCADE" json:"disaster,omitempty"`
}

// TableName はGORMのテーブル名を指定する
func (DisasterArea) TableName() string {
	return "disaster_areas"
}
