package models

import (
	"time"
)

type AuditPlanAuditableItem struct {
	ID              uint      `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	AuditPlanID     uint      `json:"audit_plan_id" gorm:"not null"`
	AuditableItemID uint      `json:"auditable_item_id" gorm:"not null"`
	CreatedBYID     uint      `json:"created_by_id"`
	CreatedAt       time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	AuditPlan     *AuditPlan     `gorm:"foreignKey:AuditPlanID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	AuditableItem *AuditableItem `gorm:"foreignKey:AuditableItemID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	CreatedByUser *User          `gorm:"foreignKey:CreatedBYID"`
}

// TableName overrides the table name used by AuditPlanAuditableItem to `audit_plan_auditable_items`
func (AuditPlanAuditableItem) TableName() string {
	return "audit_plan_auditable_items"
}
