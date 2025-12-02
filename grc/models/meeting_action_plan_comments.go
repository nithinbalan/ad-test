package models

import (
	"time"
)

type MeetingActionPlanComment struct {
	ID                     uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	MeetingActionPlanID uint      `json:"meeting_action_plan_id" gorm:"not null"`
	Comment                string    `json:"comment" binding:"required"`
	CreatedBYID            uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt              time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	MeetingActionPlan *MeetingActionPlan                  `gorm:"foreignKey:MeetingActionPlanID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedByUser        *User                                  `gorm:"foreignKey:CreatedBYID"`
	Documents            []*MeetingActionPlanCommentDocument `gorm:"foreignKey:MeetingActionPlanCommentID"`
}

func (MeetingActionPlanComment) TableName() string {
	return "meeting_action_plan_comments"
}
