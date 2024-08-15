package models

import (
	"time"
)

type DocumentAcknowledgeUser struct {
	ID          uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	DocumentID  uint      `gorm:"not null" json:"document_id"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	CreatedBYID uint      `gorm:"not null" json:"created_by_id"`
	CreatedAt   time.Time `gorm:"type:datetime;not null" json:"created_at"`

	// Relationships
	Document      Document `gorm:"foreignKey:DocumentID;constraint:OnDelete:NO ACTION,OnUpdate:NO ACTION"`
	User          User     `gorm:"foreignKey:UserID;constraint:OnDelete:NO ACTION,OnUpdate:NO ACTION"`
	CreatedByUser User     `gorm:"foreignKey:CreatedBYID;"`
}

func (DocumentAcknowledgeUser) TableName() string {
	return "document_acknowledge_users"
}
