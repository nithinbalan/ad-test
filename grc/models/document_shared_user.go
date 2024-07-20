package models

import (
	"time"
)

type DocumentSharedUser struct {
	ID          uint      `json:"id" binding:"required" gorm:"primaryKey;autoIncrement"`
	DocumentID  uint      `json:"document_id" binding:"required" gorm:"not null"`
	UserID      uint      `json:"user_id" binding:"required" gorm:"not null"`
	Type        string    `json:"type" binding:"required,oneof=created-user review-user responsible-user shared-user" gorm:"type:varchar(20);not null;default:'shared-user';check:type IN ('created-user', 'review-user', 'responsible-user', 'shared-user')"`
	CreatedByID uint      `json:"created_by_id" binding:"required" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" binding:"required" gorm:"type:datetime;not null"`

	// Relationships
	Document      *Document `gorm:"foreignKey:DocumentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	User          *User     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User     `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (DocumentSharedUser) TableName() string {
	return "document_shared_users"
}
