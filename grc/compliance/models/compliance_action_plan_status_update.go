package models

import (
	"time"

	masterModels "jethur/grc/masters/models"
	"jethur/grc/models"
)

type ComplianceActionPlanStatusUpdate struct {
	ID                     uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ComplianceActionPlanID uint       `json:"compliance_action_plan_id" gorm:"not null"`
	ActionPlanStatusID     uint       `json:"action_plan_status_id" gorm:"not null"`
	Percentage             *uint      `json:"percentage" binding:"omitempty,min=0,max=100" gorm:"not null;default:0"`
	Comment                string     `json:"comment" binding:"required" gorm:"type:nvarchar(max);not null"`
	CreatedAt              time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt              *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBy              uint       `json:"created_by" gorm:"not null"`
	UpdatedBy              *uint      `json:"updated_by" gorm:"default:null"`

	// Relationships
	ComplianceActionPlan *ComplianceActionPlan          `gorm:"foreignkey:ComplianceActionPlanID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ActionPlanStatus     *masterModels.ActionPlanStatus `gorm:"foreignkey:ActionPlanStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser        *models.User                   `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser        *models.User                   `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by ComplianceActionPlanStatusUpdate to `compliance_action_plan_status_updates`
func (ComplianceActionPlanStatusUpdate) TableName() string {
	return "compliance_action_plan_status_updates"
}
