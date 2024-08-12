package models

import (
	"time"
)

type AuditPlanAuditMember struct {
	ID          uint      `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	AuditPlanID uint      `json:"audit_plan_id" gorm:"not null"`
	UserID      uint      `json:"user_id" gorm:"not null"`
	CreatedBYID uint      `json:"created_by_id"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	AuditPlan     *AuditPlan `gorm:"foreignKey:AuditPlanID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	User          *User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User      `gorm:"foreignKey:CreatedBYID"`
}

// TableName overrides the table name used by AuditPlanAuditMember to `audit_plan_audit_members`
func (AuditPlanAuditMember) TableName() string {
	return "audit_plan_audit_members"
}
