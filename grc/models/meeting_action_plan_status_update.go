package models

import (
	"time"
)

// MeetingActionPlanStatusUpdate represents the meeting_action_plan_status_updates table in the database.
type MeetingActionPlanStatusUpdate struct {
	ID                  uint       `json:"id" gorm:"primaryKey;autoIncrement" binding:"-"`
	MeetingActionPlanID uint       `json:"meeting_action_plan_id" gorm:"not null" binding:"required"`
	ActionPlanStatusID  uint       `json:"action_plan_status_id" gorm:"not null" binding:"required"`
	Comment             string     `json:"comment" gorm:"type:nvarchar(max);not null" binding:"required"`
	Percentage          uint       `json:"percentage" gorm:"default:0;not null" binding:"required"`
	CreatedAt           time.Time  `json:"created_at" gorm:"not null" binding:"required"`
	UpdatedAt           *time.Time `json:"updated_at,omitempty" gorm:"default:null" binding:"-"`
	CreatedByID         uint       `json:"created_by_id" gorm:"not null; column:created_by_id" binding:"required"`
	UpdatedByID         *uint      `json:"updated_by_id,omitempty" gorm:"default:null; column:updated_by_id" binding:"-"`

	// Relationships
	MeetingActionPlan *MeetingActionPlan                       `json:"meeting_action_plan" gorm:"foreignKey:MeetingActionPlanID"`
	ActionPlanStatus  *ActionPlanStatus                        `json:"action_plan_status" gorm:"foreignKey:ActionPlanStatusID"`
	CreateByUser      *User                                    `json:"create_by_user" gorm:"foreignKey:CreatedByID"`
	UpdateByUser      *User                                    `json:"update_by_user,omitempty" gorm:"foreignKey:UpdatedByID"`
	Documents         []*MeetingActionPlanStatusUpdateDocument `json:"documents,omitempty" gorm:"foreignKey:MeetingActionPlanStatusUpdateID"`
}

// TableName overrides the table name used by GORM
func (MeetingActionPlanStatusUpdate) TableName() string {
	return "meeting_action_plan_status_updates"
}
