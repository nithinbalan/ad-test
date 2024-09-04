package models

import (
	"time"
)

type AuditDraftFindingDocument struct {
	ID                  uint      `json:"id" gorm:"primary_key;auto_increment"`
	AuditDraftFindingID uint      `json:"audit_draft_finding_id"`
	URL                 string    `json:"url"`
	Name                string    `json:"name" binding:"required"`
	Size                int       `json:"size"`
	Token               string    `json:"token"`
	OriginalName        string    `json:"original_name"`
	Type                string    `json:"type"`
	CreatedBYID         uint      `json:"created_by_id" gorm:"not null"`
	UpdatedBYID         uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt           time.Time `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt           time.Time `json:"updated_at" gorm:"type:datetime;not null"`

	// Relationships
	AuditDraftFinding *AuditDraftFinding `gorm:"foreignKey:AuditDraftFindingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser     *User              `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser     *User              `gorm:"foreignKey:UpdatedBYID"`
}

func (AuditDraftFindingDocument) TableName() string {
	return "audit_draft_finding_documents"
}
