package models

import (
	"time"
)

type MeetingAssociatedDocument struct {
	ID          uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	MeetingID   uint      `json:"meeting_id" binding:"required"`
	DocumentID  uint      `json:"document_id" binding:"required"`
	CreatedByID uint      `json:"created_by_id" binding:"required"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	MeetingRegister *Meeting  `gorm:"foreignKey:MeetingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Document        *Document `gorm:"foreignKey:DocumentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser   *User     `gorm:"foreignKey:CreatedByID"`
}

func (MeetingAssociatedDocument) TableName() string {
	return "meeting_associated_documents"
}
