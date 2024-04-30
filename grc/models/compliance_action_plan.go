package models

import (
	"time"

	masterModels "jethur/grc/masters/models"
)

type ComplianceActionPlan struct {
	ID                      uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ComplianceRegisterID    uint       `json:"compliance_register_id" binding:"required" gorm:"not null"`
	ComplianceRequirementID *uint      `json:"compliance_requirement_id" binding:"omitempty"  gorm:"default:null"`
	Title                   string     `json:"title" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	Description             *string    `json:"description" gorm:"type:nvarchar(max);null"`
	SeverityLevelID         *uint      `json:"severity_level_id" binding:"omitempty" gorm:"not null"`
	ResponsibleUserID       uint       `json:"responsible_user_id" binding:"required" gorm:"not null"`
	Comments                *string    `json:"comments" gorm:"type:nvarchar(max);null"`
	TargetDate              *time.Time `json:"target_date" binding:"omitempty" gorm:"type:date,default:null"`
	ActionPlanStatusID      uint       `json:"action_plan_status_id" binding:"required" gorm:"not null"`
	Percentage              *uint      `json:"percentage" binding:"omitempty,min=0,max=100" gorm:"not null;default:0"`
	CreatedAt               time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt               *time.Time `json:"updated_at" gorm:"type:datetime"`
	CreatedBy               uint       `json:"created_by" gorm:"not null"`
	UpdatedBy               *uint      `json:"updated_by" gorm:"default:null"`

	// Relationships
	ComplianceRegister    *ComplianceRegister            `gorm:"foreignkey:ComplianceRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceRequirement *ComplianceRequirement         `gorm:"foreignkey:ComplianceRequirementID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SeverityLevel         *masterModels.SeverityLevel    `gorm:"foreignkey:SeverityLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUser       *User                          `gorm:"foreignkey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ActionPlanStatus      *masterModels.ActionPlanStatus `gorm:"foreignkey:ActionPlanStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser         *User                          `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser         *User                          `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by ComplianceActionPlan to `compliance_action_plans`
func (ComplianceActionPlan) TableName() string {
	return "compliance_action_plans"
}
