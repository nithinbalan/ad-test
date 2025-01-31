package models

import (
	"time"
)

type ComplianceActionPlanStatusUpdate struct {
	ID                     uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ComplianceActionPlanID uint       `json:"compliance_action_plan_id" gorm:"not null"`
	ActionPlanStatusID     uint       `json:"action_plan_status_id" gorm:"not null"`
	Percentage             *int       `json:"percentage" binding:"omitempty,min=0,max=100" gorm:"not null;default:0"`
	Comment                string     `json:"comment" binding:"required" gorm:"type:nvarchar(max);not null"`
	CreatedAt              time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt              *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedByID            uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID            *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	ComplianceActionPlan *ComplianceActionPlan                       `gorm:"foreignKey:ComplianceActionPlanID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ActionPlanStatus     *ActionPlanStatus                           `gorm:"foreignKey:ActionPlanStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser        *User                                       `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser        *User                                       `gorm:"foreignKey:UpdatedByID"`
	Documents            []*ComplianceActionPlanStatusUpdateDocument `gorm:"foreignKey:ComplianceActionPlanStatusUpdateID"`
}

// TableName overrides the table name used by ComplianceActionPlanStatusUpdate to `compliance_action_plan_status_updates`
func (ComplianceActionPlanStatusUpdate) TableName() string {
	return "compliance_action_plan_status_updates"
}
