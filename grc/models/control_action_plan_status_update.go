package models

import (
	"time"
)

// ControlActionPlanStatusUpdate represents the control_action_plan_status_updates table in the database.
type ControlActionPlanStatusUpdate struct {
	ID                  uint       `json:"id" gorm:"primaryKey;autoIncrement" binding:"-"`
	ControlActionPlanID uint       `json:"control_action_plan_id" gorm:"not null" binding:"required"`
	ActionPlanStatusID  uint       `json:"action_plan_status_id" gorm:"not null" binding:"required"`
	Comment             string     `json:"comment" gorm:"type:nvarchar(max);not null" binding:"required"`
	Percentage          uint       `json:"percentage" gorm:"default:0;not null" binding:"required"`
	CreatedAt           time.Time  `json:"created_at" gorm:"not null" binding:"required"`
	UpdatedAt           *time.Time `json:"updated_at" gorm:"default:null" binding:"-"`
	CreatedBy           uint       `json:"created_by_id" gorm:"not null; column:created_by_id" binding:"required"`
	UpdatedBy           *uint      `json:"updated_by_id" gorm:"default:null; column:updated_by_id" binding:"-"`

	// Relationships
	ControlActionPlan *ControlActionPlan                       `json:"control_action_plan" gorm:"foreignKey:ControlActionPlanID"`
	ActionPlanStatus  *ActionPlanStatus                        `json:"action_plan_status" gorm:"foreignKey:ActionPlanStatusID"`
	CreateByUser      *User                                    `json:"create_by_user" gorm:"foreignKey:CreatedBy"`
	UpdateByUser      *User                                    `json:"update_by_user" gorm:"foreignKey:UpdatedBy"`
	Documents         []*ControlActionPlanStatusUpdateDocument `json:"documents" gorm:"foreignKey:ControlActionPlanStatusUpdateID"`
}

// TableName overrides the table name used by GORM
func (ControlActionPlanStatusUpdate) TableName() string {
	return "control_action_plan_status_updates"
}
