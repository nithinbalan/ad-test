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
	UpdatedBy            *uint     `json:"updated_by" gorm:"default:null" binding:"omitempty"`
}

type NotificationMaster struct {
	ID         uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	ModuleID   uint   `json:"module_id" gorm:"not null"  binding:"required"`
	Identifier string `json:"identifier" gorm:"size:100;not null" binding:"required"`
	TitleEn    string `json:"title_en" gorm:"size:100;not null" binding:"required"`
	TitleAr    string `json:"title_ar" gorm:"size:100;not null" binding:"required"`

	Module Module `gorm:"foreignKey:ModuleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	// One-to-many relationship with UserNotificationSetting
	UserNotificationSettings *[]UserNotificationSetting `gorm:"foreignKey:NotificationMasterID"`
}
