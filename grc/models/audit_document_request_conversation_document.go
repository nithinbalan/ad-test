package models

import (
	"time"
)

type AuditDocumentRequestConversationDocument struct {
	ID                                 uint      `json:"id" gorm:"primary_key;auto_increment"`
	AuditDocumentRequestConversationID uint      `json:"audit_document_request_conversation_id"`
	URL                                string    `json:"url"`
	Name                               string    `json:"name" binding:"required"`
	Size                               int       `json:"size"`
	Token                              string    `json:"token"`
	OriginalName                       string    `json:"original_name"`
	Type                               string    `json:"type"`
	CreatedByID                        uint      `json:"created_by_id" gorm:"not null"`
	UpdatedByID                        uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                          time.Time `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                          time.Time `json:"updated_at" gorm:"type:datetime;not null"`

	// Relationships
	AuditDocumentRequestConversation *AuditDocumentRequestConversation `gorm:"foreignKey:AuditDocumentRequestConversationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser                    *User                             `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser                    *User                             `gorm:"foreignKey:UpdatedByID"`
}

func (AuditDocumentRequestConversationDocument) TableName() string {
	return "audit_document_request_conversation_documents"
}
