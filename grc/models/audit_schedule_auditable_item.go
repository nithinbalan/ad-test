package models

import (
	"time"
)

type AuditScheduleAuditableItem struct {
	ID              uint      `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	AuditScheduleID uint      `json:"audit_schedule_id" gorm:"not null"`
	AuditableItemID uint      `json:"auditable_item_id" gorm:"not null"`
	CreatedBYID     uint      `json:"created_by_id"`
	CreatedAt       time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	AuditSchedule *AuditSchedule `gorm:"foreignKey:AuditScheduleID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	AuditableItem *AuditableItem `gorm:"foreignKey:AuditableItemID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	CreatedByUser *User          `gorm:"foreignKey:CreatedBYID"`
}

// TableName overrides the table name used by AuditPlanScheduleAuditableItem to `audit_plan_schedule_auditable_items`
func (AuditScheduleAuditableItem) TableName() string {
	return "audit_schedule_auditable_items"
}
