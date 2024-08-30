package models

import (
	"time"
)

type AuditChecklist struct {
	ID                   uint       `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	AuditID              uint       `json:"audit_id" binding:"required" gorm:"not null"`
	AuditPlanChecklistID *uint      `json:"audit_plan_checklist_id"`
	Question             string     `json:"question" binding:"required" gorm:"type:nvarchar(max);not null"`
	ConformanceID        *uint      `json:"conformance_id"`
	Comment              *string    `json:"comment" gorm:"type:nvarchar(max)"`
	AnsweredByID         *uint      `json:"answered_by_id"`
	AnsweredAt           *time.Time `json:"answered_at" gorm:"type:datetime"`
	CreatedByID          uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	CreatedAt            time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedByID          *uint      `json:"updated_by_id"`
	UpdatedAt            *time.Time `json:"updated_at" gorm:"type:datetime"`
	StatusID             uint       `json:"status_id" binding:"required" gorm:"not null"`

	// Relationships
	Audit              *Audit              `gorm:"foreignKey:AuditID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	AuditPlanChecklist *AuditPlanChecklist `gorm:"foreignKey:AuditPlanChecklistID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	Conformance        *Conformance        `gorm:"foreignKey:ConformanceID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	AnsweredByUser     *User               `gorm:"foreignKey:AnsweredByID"` // Nullable relationship
	CreatedByUser      *User               `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser      *User               `gorm:"foreignKey:UpdatedByID"` // Nullable relationship
	Status             *Status             `gorm:"foreignKey:StatusID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
}

// TableName overrides the default table name used by GORM
func (AuditChecklist) TableName() string {
	return "audit_checklists"
}
