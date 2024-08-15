package models

import (
	"time"
)

type UserNotificationSetting struct {
	ID                   uint      `json:"id" gorm:"primaryKey;autoIncrement" binding:"-"`
	UserID               uint      `json:"users_id" gorm:"not null;index" binding:"required"`
	NotificationMasterID uint      `json:"notification_master_id" gorm:"not null" binding:"required"`
	IsEnable             bool      `json:"is_enable" gorm:"not null;type:tinyint;default:0" binding:"required"`
	UpdatedAt            time.Time `json:"updated_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP" binding:"-"`
	UpdatedByID          *uint     `json:"updated_by_id" gorm:"default:null" binding:"omitempty"`

	// Relationships
	User               *User               `gorm:"foreignKey:UserID"`
	NotificationMaster *NotificationMaster `gorm:"foreignKey:NotificationMasterID"`
}

func (UserNotificationSetting) TableName() string {
	return "user_notification_settings"
}
