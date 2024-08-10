package models

import (
	"time"
)

type AuditPlanOrganizationUnit struct {
	ID                 uint      `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	AuditPlanID        uint      `json:"audit_plan_id" gorm:"not null"`
	OrganizationUnitID uint      `json:"organization_unit_id" gorm:"not null"`
	CreatedBYID        uint      `json:"created_by_id"`
	CreatedAt          time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	AuditPlan        *AuditPlan        `gorm:"foreignKey:AuditPlanID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	OrganizationUnit *OrganizationUnit `gorm:"foreignKey:OrganizationUnitID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	CreatedByUser    *User             `gorm:"foreignKey:CreatedBYID"`
}

// TableName overrides the table name used by AuditPlanOrganizationUnit to `audit_plan_organization_units`
func (AuditPlanOrganizationUnit) TableName() string {
	return "audit_plan_organization_units"
}
