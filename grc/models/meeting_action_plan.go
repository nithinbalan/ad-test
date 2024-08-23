package models

import (
	"time"
)

type MeetingActionPlan struct {
	ID                 uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	MeetingID          uint       `json:"meeting_id" gorm:"not null" binding:"required"`
	MeetingMinuteID    *uint      `json:"meeting_minute_id" gorm:"null"`
	Title              string     `json:"title" gorm:"type:nvarchar(500);not null" binding:"required"`
	Description        *string    `json:"description" gorm:"type:nvarchar(max)" binding:"omitempty"`
	SeverityLevelID    *uint      `json:"severity_level_id" gorm:"default:null" binding:"omitempty"`
	TargetDate         *time.Time `json:"target_date" binding:"omitempty" gorm:"type:date,default:null"`
	ResponsibleUserID  uint       `json:"responsible_user_id" gorm:"not null" binding:"required"`
	ActionPlanStatusID uint       `json:"action_plan_status_id" gorm:"not null" binding:"required"`
	Percentage         uint       `json:"percentage" gorm:"default:0;not null" binding:"required"`
	Comments           *string    `json:"comments" gorm:"type:nvarchar(max);null"`
	CreatedAt          time.Time  `json:"created_at" gorm:"not null" binding:"required"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty" gorm:"default:null"`
	CreatedByID        uint       `json:"created_by_id" gorm:"not null; column:created_by_id" binding:"required"`
	UpdatedByID        *uint      `json:"updated_by_id,omitempty" gorm:"default:null; column:updated_by_id"`

	// Relationships
	Meeting                      *Meeting                         `gorm:"foreignKey:MeetingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	MeetingMinute                *MeetingMinute                   `gorm:"foreignKey:MeetingMinuteID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SeverityLevel                *SeverityLevel                   `gorm:"foreignkey:SeverityLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUser              *User                            `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ActionPlanStatus             *ActionPlanStatus                `gorm:"foreignKey:ActionPlanStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ActionPlanStatusUpdates      []*MeetingActionPlanStatusUpdate `gorm:"foreignKey:MeetingActionPlanID"`
	LatestActionPlanStatusUpdate *MeetingActionPlanStatusUpdate   `gorm:"foreignKey:MeetingActionPlanID"`
	CreatedByUser                *User                            `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser                *User                            `gorm:"foreignKey:UpdatedByID"`
}

func (MeetingActionPlan) TableName() string {
	return "meeting_action_plans"
}
