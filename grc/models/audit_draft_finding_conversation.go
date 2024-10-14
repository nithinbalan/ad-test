package models

import (
	"time"
)

type AuditDraftFindingConversation struct {
	ID                  uint       `gorm:"primaryKey;autoIncrement"`
	AuditDraftFindingID uint       `json:"audit_draft_finding_id" binding:"required" gorm:"not null"`
	Comment             string     `json:"comment" binding:"required" gorm:"type:nvarchar(MAX);not null"`
	Type                string     `gorm:"type:varchar(20);not null;check:type IN ('DR', 'AC')"`
	IsAccepted          bool       `json:"is_accepted" binding:"required" gorm:"type:tinyint(1);not null;default:0"`
	IsDismissed         bool       `json:"is_dismissed" binding:"required" gorm:"type:tinyint(1);not null;default:0"`
	CreatedAt           time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt           *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBYID         uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID         *uint      `json:"updated_by_id" gorm:"default:null"`

	AuditDraftFinding AuditDraftFinding                        `gorm:"foreignKey:AuditDraftFindingID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	Documents         []*AuditDraftFindingConversationDocument `gorm:"foreignKey:AuditDraftFindingConversationID"`
	CreatedByUser     User                                     `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser     *User                                    `gorm:"foreignKey:UpdatedBYID"`
}

// TableName overrides the table name used by AuditDraftFindingConversation to `audit_draft_finding_conversations`
func (AuditDraftFindingConversation) TableName() string {
	return "audit_draft_finding_conversations"
}
