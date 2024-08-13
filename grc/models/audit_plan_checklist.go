package models

import (
	"time"
)

type AuditPlanChecklist struct {
	ID              uint       `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	AuditPlanID     uint       `json:"audit_plan_id" gorm:"not null"`
	AuditableItemID uint       `json:"auditable_item_id" gorm:"not null"`
	Question        string     `json:"question" gorm:"type:nvarchar(255);not null"`
	CreatedAt       time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt       *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBYID     uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID     *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	AuditPlan     *AuditPlan     `gorm:"foreignKey:AuditPlanID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	AuditableItem *AuditableItem `gorm:"foreignKey:AuditableItemID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	CreatedByUser *User          `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser *User          `gorm:"foreignKey:UpdatedBYID"`
}

// TableName overrides the table name used by AuditPlan to `audit_plans`
func (AuditPlanChecklist) TableName() string {
	return "audit_plan_checklists"
}
