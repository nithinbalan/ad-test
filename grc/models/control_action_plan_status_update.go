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
	CreatedByID         uint       `json:"created_by" gorm:"not null; column:created_by" binding:"required"`
	UpdatedByID         *uint      `json:"updated_by" gorm:"default:null; column:updated_by" binding:"-"`

	// Relationships
	ControlActionPlan *ControlActionPlan `json:"control_action_plan" gorm:"foreignKey:ControlActionPlanID"`
	ActionPlanStatus  *ActionPlanStatus  `json:"action_plan_status" gorm:"foreignKey:ActionPlanStatusID"`
	CreateByUser      *User              `json:"creator" gorm:"foreignKey:CreatedByID"`
	UpdateByUser      *User              `json:"updater" gorm:"foreignKey:UpdatedByID"`
}

// TableName overrides the table name used by GORM
func (ControlActionPlanStatusUpdate) TableName() string {
	return "control_action_plan_status_updates"
}
