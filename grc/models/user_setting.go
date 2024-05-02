package models

import (
	"time"
)

// UserSetting represents a user's settings in the system.
type UserSetting struct {
	ID             uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         uint       `gorm:"not null;index" json:"user_id"`
	TimezoneID     uint       `gorm:"not null" json:"timezone_id"`
	Language       string     `gorm:"type:varchar(2);not null;default:'en'" json:"language"`
	Timeout        int        `gorm:"not null" json:"timeout"`
	IsNotification bool       `gorm:"type:tinyint;not null; default:0" json:"is_notification"`
	UpdatedAt      *time.Time `gorm:"type:timestamp; not null; default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
	UpdatedBy      *uint      `gorm:"default:null" json:"updated_by"`

	// Foreign key relationships
	User          User     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Timezone      Timezone `gorm:"foreignKey:TimezoneID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	UpdatedByUser User     `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by UserSetting to `user_settings`
func (UserSetting) TableName() string {
	return "user_settings"
}
