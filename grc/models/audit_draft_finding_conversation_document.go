package models

import (
	"time"
)

type AuditDraftFindingConversationDocument struct {
	ID                              uint      `json:"id" gorm:"primary_key;auto_increment"`
	AuditDraftFindingConversationID uint      `json:"audit_draft_finding_conversation_id"`
	URL                             string    `json:"url"`
	Name                            string    `json:"name" binding:"required"`
	Size                            int       `json:"size"`
	Token                           string    `json:"token"`
	OriginalName                    string    `json:"original_name"`
	Type                            string    `json:"type"`
	CreatedBYID                     uint      `json:"created_by_id" gorm:"not null"`
	UpdatedBYID                     uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                       time.Time `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                       time.Time `json:"updated_at" gorm:"type:datetime;not null"`

	// Relationships
	AuditDraftFindingConversation *AuditDraftFindingConversation           `gorm:"foreignKey:AuditDraftFindingConversationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Documents                     []*AuditDraftFindingConversationDocument `gorm:"foreignKey:AuditDraftFindingConversationID"`
	CreatedByUser                 *User                                    `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser                 *User                                    `gorm:"foreignKey:UpdatedBYID"`
}

func (AuditDraftFindingConversationDocument) TableName() string {
	return "audit_draft_finding_conversation_documents"
}
