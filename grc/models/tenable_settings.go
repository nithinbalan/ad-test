package models

import (
	"time"
)

type TenableSettings struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	BaseURL            string    `gorm:"type:varchar(255);not null" json:"base_url"`
	AccessKey          string    `gorm:"type:varchar(300);not null" json:"access_key"`
	SecretKey          string    `gorm:"type:varchar(300);not null" json:"secret_key"`
	SyncFrequencyHours int       `gorm:"not null;default:24" json:"sync_frequency_hours"`
	LookbackDays       int       `gorm:"not null;default:30" json:"lookback_days"`
	Status             string    `gorm:"type:varchar(10);not null;default:'active'" json:"status"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (TenableSettings) TableName() string {
	return "tenable_settings"
}
