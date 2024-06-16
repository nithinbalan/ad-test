package models

import (
	"time"
)

type SLAActionPlanStatusUpdate struct {
	ID                 uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	SLAActionPlanID    uint       `json:"sla_action_plan_id" gorm:"not null"`
	ActionPlanStatusID uint       `json:"action_plan_status_id" gorm:"not null"`
	Percentage         *uint      `json:"percentage" binding:"omitempty,min=0,max=100" gorm:"not null;default:0"`
	Comment            string     `json:"comment" binding:"required" gorm:"type:nvarchar(max);not null"`
	CreatedAt          time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt          *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBy          uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBy          *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	SLAActionPlan    *SLAActionPlan                       `gorm:"foreignkey:SLAActionPlanID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ActionPlanStatus *ActionPlanStatus                    `gorm:"foreignkey:ActionPlanStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser    *User                                `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser    *User                                `gorm:"foreignKey:UpdatedBy"`
	Documents        []*SLAActionPlanStatusUpdateDocument `gorm:"foreignKey:SLAActionPlanStatusUpdateID"`
}

// TableName overrides the table name used by SLAActionPlanStatusUpdate to `sla_action_plan_status_updates`
func (SLAActionPlanStatusUpdate) TableName() string {
	return "sla_action_plan_status_updates"
}
