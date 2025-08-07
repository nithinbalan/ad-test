package models

import (
	"time"
)

type AuditFindingAssociatedDocument struct {
	ID             uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	AuditFindingID uint      `json:"audit_finding_id" binding:"required"`
	DocumentID     uint      `json:"document_id" binding:"required"`
	CreatedByID    uint      `json:"created_by_id" binding:"required"`
	CreatedAt      time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	AuditFinding  *AuditFinding `gorm:"foreignKey:AuditFindingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Document      *Document     `gorm:"foreignKey:DocumentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User         `gorm:"foreignKey:CreatedByID"`
}

func (AuditFindingAssociatedDocument) TableName() string {
	return "audit_finding_associated_documents"
}
