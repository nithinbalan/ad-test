package models

import (
	"time"
)

type DocumentEditHistory struct {
	ID                uint      `json:"id" gorm:"primaryKey;autoIncrement" binding:"required"`
	DocumentID        uint      `json:"document_id" gorm:"not null;index" binding:"required"`
	DocumentVersionID uint      `json:"document_version_id" gorm:"not null;index" binding:"required"`
	ActionType        string    `json:"action_type" gorm:"type:varchar(20);not null;default:'created'" binding:"required,oneof=create update check-in check-out"`
	Comment           *string   `json:"comment" gorm:"type:nvarchar(max);default:null"`
	URL               *string   `json:"url" gorm:"type:varchar(1000)" binding:"omitempty"`
	Name              *string   `json:"name" gorm:"type:nvarchar(255)" binding:"omitempty"`
	Size              *int      `json:"size" gorm:"default:null" binding:"omitempty"`
	Token             *string   `json:"token" gorm:"type:varchar(255)" binding:"omitempty"`
	OriginalName      *string   `json:"original_name" gorm:"type:nvarchar(255)" binding:"omitempty"`
	Type              *string   `json:"type" binding:"required" gorm:"type:varchar(255);not null"`
	CreatedBYID       uint      `json:"created_by_id" gorm:"not null" binding:"required"`
	CreatedAt         time.Time `json:"created_at" gorm:"type:datetime;not null" binding:"required"`

	// Relations
	Document  *Document `gorm:"foreignKey:DocumentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedBy *User     `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (DocumentEditHistory) TableName() string {
	return "document_edit_histories"
}
