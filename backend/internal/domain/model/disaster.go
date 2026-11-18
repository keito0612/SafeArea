package model

import (
	"time"

	"github.com/google/uuid"
)

// DisasterType は災害種別を表す
type DisasterType string

const (
	DisasterTypeEarthquake DisasterType = "earthquake" // 地震
	DisasterTypeFlood      DisasterType = "flood"      // 洪水
	DisasterTypeTsunami    DisasterType = "tsunami"    // 津波
	DisasterTypeLandslide  DisasterType = "landslide"  // 土砂災害
	DisasterTypeTyphoon    DisasterType = "typhoon"    // 台風
	DisasterTypeOther      DisasterType = "other"      // その他
)

// Disaster は過去の災害データを表すエンティティ
type Disaster struct {
	ID             uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Type           DisasterType `gorm:"type:varchar(20);not null;index" json:"type"`
	Name           string       `gorm:"type:varchar(255);not null" json:"name"`
	Year           int          `gorm:"not null;index" json:"year"`
	Damage         string       `gorm:"type:text" json:"damage"`
	PrefectureCode string       `gorm:"type:varchar(2);not null;index" json:"prefecture_code"`
	CityName       string       `gorm:"type:varchar(255)" json:"city_name"`
	CreatedAt      time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time    `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	DisasterAreas []DisasterArea `gorm:"foreignKey:DisasterID" json:"disaster_areas,omitempty"`
}

// TableName はGORMのテーブル名を指定する
func (Disaster) TableName() string {
	return "disasters"
}
