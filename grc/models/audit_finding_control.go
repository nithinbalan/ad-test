package models

import (
	"time"
)

type AuditFindingControl struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	AuditFindingID uint      `json:"audit_finding_id" binding:"required"`
	ControlID      uint      `json:"control_id" binding:"required"`
	CreatedBYID    uint      `json:"created_by_id" binding:"required"`
	CreatedAt      time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	AuditFinding  *AuditFinding `gorm:"foreignKey:AuditFindingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Control       *Control      `gorm:"foreignKey:ControlID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User         `gorm:"foreignKey:CreatedBYID"`
}

func (AuditFindingControl) TableName() string {
	return "audit_finding_controls"
}
