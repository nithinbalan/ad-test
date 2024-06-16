package models

import (
	"time"
)

type UserNotification struct {
	ID          uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	UserID      uint      `json:"user_id" binding:"required" gorm:"not null"`
	ModuleID    uint      `json:"module_id" binding:"required" gorm:"not null"`
	Icon        Icon      `json:"icon" binding:"required" gorm:"size:45;not null"`
	MessageEn   string    `json:"message_en" gorm:"size:255;not null" binding:"required"`
	MessageAr   string    `json:"message_ar" gorm:"size:255;not null" binding:"required"`
	RedirectURL string    `json:"redirect_url" gorm:"size:255;not null" binding:"required"`
	IsNew       bool      `json:"is_new" gorm:"type:boolean;default:true"`
	IsRead      bool      `json:"is_read" gorm:"type:boolean;default:false"`
	CreatedByID uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	User          *User   `gorm:"foreignKey:UserID;constraint:OnDelete:NO ACTION,OnUpdate:NO ACTION"`
	Module        *Module `gorm:"foreignKey:ModuleID;constraint:OnDelete:NO ACTION,OnUpdate:NO ACTION"`
	CreatedByUser *User   `gorm:"foreignKey:CreatedByID"`
}

// TableName overrides the table name used by UserNotification to `user_notifications`
func (UserNotification) TableName() string {
	return "user_notifications"
}

type Icon string

const (
	CreateIcon Icon = "create"
	UpdateIcon Icon = "update"
	DeleteIcon Icon = "delete"
)
