package models

import (
	"time"
)

type AuditFindingIncident struct {
	ID                 uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	AuditFindingID     uint      `json:"audit_finding_id" binding:"required"`
	IncidentRegisterID uint      `json:"incident_register_id" binding:"required"`
	CreatedBYID        uint      `json:"created_by_id" binding:"required"`
	CreatedAt          time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	AuditFinding     *AuditFinding     `gorm:"foreignKey:AuditFindingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	IncidentRegister *IncidentRegister `gorm:"foreignKey:IncidentRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser    *User             `gorm:"foreignKey:CreatedBYID"`
}

func (AuditFindingIncident) TableName() string {
	return "audit_finding_incidents"
}
