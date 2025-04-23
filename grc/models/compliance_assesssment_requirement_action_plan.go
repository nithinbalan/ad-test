package models

import (
	"time"
)

type ComplianceAssessmentRequirementActionPlan struct {
	ID                                uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ComplianceAssessmentRequirementID *uint      `json:"compliance_requirement_id" binding:"omitempty"  gorm:"default:null"`
	Title                             string     `json:"title" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	Description                       *string    `json:"description" gorm:"type:nvarchar(max);null"`
	ResponsibleUserID                 uint       `json:"responsible_user_id" binding:"required" gorm:"not null"`
	ResponsibleUnitID                 uint       `json:"responsible_unit_id" binding:"required" gorm:"not null"`
	TargetDate                        *time.Time `json:"target_date" binding:"omitempty" gorm:"type:date,default:null"`
	ActionPlanStatusID                uint       `json:"action_plan_status_id" binding:"required" gorm:"not null"`
	Percentage                        *int       `json:"percentage" binding:"omitempty,min=0,max=100" gorm:"not null;default:0"`
	CreatedAt                         time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                         *time.Time `json:"updated_at" gorm:"type:datetime"`
	CreatedByID                       uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID                       *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	ComplianceAssessmentRequirement *ComplianceAssessmentRequirement `gorm:"foreignKey:ComplianceAssessmentRequirementID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUser                 *User                            `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUnit                 *OrganizationUnit                `gorm:"foreignKey:ResponsibleUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ActionPlanStatus                *ActionPlanStatus                `gorm:"foreignKey:ActionPlanStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser                   *User                            `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser                   *User                            `gorm:"foreignKey:UpdatedByID"`
}

// TableName overrides the table name used by ComplianceAssessmentRequirementActionPlan to `compliance_action_plans`
func (ComplianceAssessmentRequirementActionPlan) TableName() string {
	return "compliance_assessment_requirement_action_plans"
}
