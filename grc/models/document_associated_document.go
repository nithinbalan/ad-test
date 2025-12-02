package models

import (
	"time"
)

type DocumentAssociatedDocument struct {
	ID                   uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	DocumentID           uint      `json:"document_id" binding:"required"`
	AssociatedDocumentID uint      `json:"associated_document_id" binding:"required"`
	CreatedByID          uint      `json:"created_by_id" binding:"required"`
	CreatedAt            time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	Document           *Document `gorm:"foreignKey:DocumentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AssociatedDocument *Document `gorm:"foreignKey:AssociatedDocumentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser      *User     `gorm:"foreignKey:CreatedByID"`
}

func (DocumentAssociatedDocument) TableName() string {
	return "document_associated_documents"
}
