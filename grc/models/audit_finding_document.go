package models

import (
	"time"
)

type AuditFindingDocument struct {
	ID             uint      `json:"id" gorm:"primary_key;auto_increment"`
	AuditFindingID uint      `json:"audit_finding_id"`
	URL            string    `json:"url"`
	Name           string    `json:"name" binding:"required"`
	Size           int       `json:"size"`
	Token          string    `json:"token"`
	OriginalName   string    `json:"original_name"`
	Type           string    `json:"type"`
	CreatedByID    uint      `json:"created_by_id" gorm:"not null"`
	UpdatedByID    uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt      time.Time `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"type:datetime;not null"`

	// Relationships
	AuditFinding  *AuditFinding `gorm:"foreignKey:AuditFindingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User         `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser *User         `gorm:"foreignKey:UpdatedByID"`
}

func (AuditFindingDocument) TableName() string {
	return "audit_finding_documents"
}
