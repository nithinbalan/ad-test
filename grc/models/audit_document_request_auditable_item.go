package models

import (
	"time"
)

type AuditDocumentRequestAuditableItem struct {
	ID                     uint      `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	AuditDocumentRequestID uint      `json:"audit_document_request_id" gorm:"not null"`
	AuditableItemID        uint      `json:"auditable_item_id" gorm:"not null"`
	CreatedBYID            uint      `json:"created_by_id"`
	CreatedAt              time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	AuditDocumentRequest *AuditDocumentRequest `gorm:"foreignKey:AuditDocumentRequestID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	AuditableItem        *AuditableItem        `gorm:"foreignKey:AuditableItemID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	CreatedByUser        *User                 `gorm:"foreignKey:CreatedBYID"`
}

func (AuditDocumentRequestAuditableItem) TableName() string {
	return "audit_document_request_auditable_items"
}
