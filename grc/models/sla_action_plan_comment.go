package models

import (
	"time"
)

type SLAActionPlanComment struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	SLAActionPlanID uint     `json:"sla_action_plan_id" gorm:"not null"`
	Comment        string    `json:"comment" binding:"required"`
	CreatedByID    uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt      time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	SLAActionPlan *SLAActionPlan                  `gorm:"foreignKey:SLAActionPlanID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedByUser *User                           `gorm:"foreignKey:CreatedByID"`
	Documents     []*SLAActionPlanCommentDocument `gorm:"foreignKey:SLAActionPlanCommentID"`
}

func (SLAActionPlanComment) TableName() string {
	return "sla_action_plan_comments"
}
