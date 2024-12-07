package models

import (
	"time"
)

type AuditChecklistDocument struct {
	ID               uint      `json:"id" gorm:"primary_key;auto_increment"`
	AuditChecklistID uint      `json:"audit_checklist_id"`
	URL              string    `json:"url"`
	Name             string    `json:"name" binding:"required"`
	Size             int       `json:"size"`
	Token            string    `json:"token"`
	OriginalName     string    `json:"original_name"`
	Type             string    `json:"type"`
	CreatedBYID      uint      `json:"created_by_id" gorm:"not null"`
	UpdatedBYID      uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt        time.Time `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"type:datetime;not null"`

	// Relationships
	AuditChecklist *AuditChecklist `gorm:"foreignKey:AuditChecklistID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser  *User           `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser  *User           `gorm:"foreignKey:UpdatedBYID"`
}

func (AuditChecklistDocument) TableName() string {
	return "audit_checklist_documents"
}
