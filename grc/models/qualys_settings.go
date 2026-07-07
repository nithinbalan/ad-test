package models

import "time"

type QualysSettings struct {
	ID           uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	PlatformURL  string `gorm:"type:varchar(255);not null" json:"platform_url"`
	Username     string `gorm:"type:varchar(300);not null" json:"username"`
	Password     string `gorm:"type:varchar(300);not null"`
	LookbackDays int    `gorm:"not null;default:30" json:"lookback_days"`
	Status       string `gorm:"type:varchar(10);not null;default:'active'" json:"status"`

	CreatedByID uint  `gorm:"not null;column:created_by_id" json:"created_by_id"`
	UpdatedByID *uint `gorm:"column:updated_by_id" json:"updated_by_id,omitempty"`

	CreatedAt time.Time  `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt *time.Time `gorm:"type:datetime" json:"updated_at,omitempty"`

	CreatedByUser *User `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (QualysSettings) TableName() string {
	return "qualys_settings"
}