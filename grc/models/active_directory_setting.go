package models

import (
	"time"
)

type ActiveDirectorySetting struct {
	ID             uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	Server         string     `json:"server" gorm:"size:255;not null"`
	Port           int        `json:"port" gorm:"not null"`
	UseSSL         bool       `json:"use_ssl" gorm:"not null"`
	BaseDN         string     `json:"base_dn" gorm:"type:nvarchar(500);not null"`
	BindDN         string     `json:"bind_dn" gorm:"type:nvarchar(500);not null"`
	BindPassword   string     `json:"bind_password" gorm:"type:nvarchar(255);not null"`
	UserSearchAttr string     `json:"user_search_attr" gorm:"size:100;not null"`
	Username       string     `json:"username" gorm:"size:100;not null"`
	Password       string     `json:"password" gorm:"type:nvarchar(255);not null"`
	CreatedBYID    uint       `json:"created_by_id" binding:"required"`
	CreatedAt      time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedBYID    *uint      `json:"updated_by_id,omitempty" gorm:"default:null"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`

	// Relationships
	CreatedByUser *User `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (ActiveDirectorySetting) TableName() string {
	return "active_directory_settings"
}
