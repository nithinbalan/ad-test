package models

import (
	"jethur/grc/masters/models"
	"time"
)

type SystemPreference struct {
	ID                  uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	CurrencyID          uint       `json:"currency_id" binding:"required" gorm:"not null"`
	TimezoneID          uint       `json:"timezone_id" binding:"required" gorm:"not null"`
	DateFormat          *string    `json:"date_format" binding:"required" gorm:"size:45"`
	TimeFormat          string     `json:"time_format" binding:"required" gorm:"type:varchar(10);not null;default:12;check:time_format IN ('12', '24')"`
	DocUploadSize       int        `json:"doc_upload_size" binding:"required" gorm:"not null"`
	LogoUploadSize      int        `json:"logo_upload_size" binding:"required" gorm:"not null"`
	IdleTimeout         int        `json:"idle_timeout" binding:"required" gorm:"not null"`
	IsMaintenanceModeOn *bool      `json:"is_maintenance_mode_on" gorm:"type:boolean"`
	CreatedBy           uint       `json:"created_by" gorm:"not null"`
	UpdatedBy           *uint      `json:"updated_by" gorm:"default:null"`
	CreatedAt           time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt           *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	Currency      *Currency `gorm:"foreignkey:CurrencyID;constraint:OnDelete:NO ACTION,OnUpdate:NO ACTION"`
	Timezone      *models.Timezone `gorm:"foreignkey:TimezoneID;constraint:OnDelete:NO ACTION,OnUpdate:NO ACTION"`
	CreatedByUser *User     `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser *User     `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by SystemPreference to `system_preferences`
func (SystemPreference) TableName() string {
	return "system_preferences"
}
