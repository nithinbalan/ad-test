package models

import (
	"time"
)

type UserNotification struct {
	ID                   uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	UserID               uint      `json:"user_id" binding:"required" gorm:"not null"`
	ModuleID             uint      `json:"module_id" binding:"required" gorm:"not null"`
	NotificationType     string    `json:"notification_type" binding:"required" gorm:"size:255;not null"`
	ReferenceNumber      string    `json:"reference_number" binding:"required" gorm:"size:45;not null"`
	HeadingEn            string    `json:"heading_en" gorm:"size:255;not null" binding:"required"`
	HeadingAr            string    `json:"heading_ar" gorm:"size:255;not null" binding:"required"`
	TitleEn              string    `json:"title_en" gorm:"size:500;not null" binding:"required"`
	TitleAr              string    `json:"title_ar" gorm:"size:500;not null" binding:"required"`
	Entity               *string   `json:"entity" gorm:"size:255"`
	Action               *string   `json:"action" gorm:"size:255"`
	URLParam             string    `json:"url_param" gorm:"type:nvarchar(max)"`
	QueryParam           string    `json:"query_param" gorm:"type:nvarchar(max)"`
	IsNew                bool      `json:"is_new" gorm:"default:true"`
	IsRead               bool      `json:"is_read" gorm:"default:false"`
	IsSystemNotification bool      `json:"is_system_notification" gorm:"default:false"`
	CreatedByID          *uint     `json:"created_by_id"`
	CreatedAt            time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	User          *User   `gorm:"foreignKey:UserID;constraint:OnDelete:NO ACTION,OnUpdate:NO ACTION"`
	Module        *Module `gorm:"foreignKey:ModuleID;constraint:OnDelete:NO ACTION,OnUpdate:NO ACTION"`
	CreatedByUser *User   `gorm:"foreignKey:CreatedByID"`
}

// TableName overrides the table name used by UserNotification to `user_notifications`
func (UserNotification) TableName() string {
	return "user_notifications"
}
