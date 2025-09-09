package models

import (
	"time"
)

type ControlActionPlanComment struct {
	ID                     uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ControlActionPlanID uint      `json:"control_action_plan_id" gorm:"not null"`
	Comment                string    `json:"comment" binding:"required"`
	CreatedBYID            uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt              time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	ControlActionPlan *ControlActionPlan                  `gorm:"foreignKey:ControlActionPlanID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedByUser        *User                                  `gorm:"foreignKey:CreatedBYID"`
	Documents            []*ControlActionPlanCommentDocument `gorm:"foreignKey:ControlActionPlanCommentID"`
}

func (ControlActionPlanComment) TableName() string {
	return "control_action_plan_comments"
}
