package models

import (
	"time"
)

type AuditDocumentRequestConversation struct {
	ID                     uint      `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	AuditDocumentRequestID uint      `json:"audit_document_request_id" gorm:"not null"`
	Comment                string    `json:"comment" gorm:"type:nvarchar(MAX);not null"`
	CommentedAt            time.Time `json:"commented_at" gorm:"type:datetime;not null"`
	CommentedBYID          *uint     `json:"commented_by_id" gorm:"null"`

	// Relationships
	AuditDocumentRequest *AuditDocumentRequest                       `gorm:"foreignKey:AuditDocumentRequestID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	Documents            []*AuditDocumentRequestConversationDocument `gorm:"foreignKey:AuditDocumentRequestConversationID"`
	CommentedByUser      *User                                       `gorm:"foreignKey:CommentedBYID"`
}

// TableName overrides the table name used by AuditDocumentRequestConversation to `audit_document_request_conversations`
func (AuditDocumentRequestConversation) TableName() string {
	return "audit_document_request_conversations"
}
