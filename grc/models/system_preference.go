package models

import (
	"time"
)

type SystemPreference struct {
	ID                  uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	CurrencyID          uint       `json:"currency_id" binding:"required" gorm:"not null"`
	TimezoneID          uint       `json:"timezone_id" binding:"required" gorm:"not null"`
	DateFormat          string     `json:"date_format" binding:"required" gorm:"size:45;not null"`
	TimeFormat          string     `json:"time_format" binding:"required" gorm:"type:varchar(10);not null;default:12;check:time_format IN ('12', '24')"`
	DocUploadSize       int        `json:"doc_upload_size" binding:"required" gorm:"not null"`
	LogoUploadSize      int        `json:"logo_upload_size" binding:"required" gorm:"not null"`
	IdleTimeout         int        `json:"idle_timeout" binding:"required" gorm:"not null"`
	IsMaintenanceModeOn *bool      `json:"is_maintenance_mode_on" gorm:"type:boolean"`
	CreatedByID         uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID         *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt           time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt           *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	Currency      *Currency `gorm:"foreignKey:CurrencyID;constraint:OnDelete:NO ACTION,OnUpdate:NO ACTION"`
	Timezone      *Timezone `gorm:"foreignKey:TimezoneID;constraint:OnDelete:NO ACTION,OnUpdate:NO ACTION"`
	CreatedByUser *User     `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser *User     `gorm:"foreignKey:UpdatedByID"`
}

// TableName overrides the table name used by SystemPreference to `system_preferences`
func (SystemPreference) TableName() string {
	return "system_preferences"
}
