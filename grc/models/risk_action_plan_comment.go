package models

import (
	"time"
)

type RiskActionPlanComment struct {
	ID                 uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	RiskActionPlanID   uint      `json:"risk_action_plan_id" gorm:"not null"`
	Comment            string    `json:"comment" binding:"required"`
	CreatedByID        uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt          time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	RiskActionPlan *RiskActionPlan                  `gorm:"foreignKey:RiskActionPlanID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedByUser  *User                            `gorm:"foreignKey:CreatedByID"`
	Documents      []*RiskActionPlanCommentDocument `gorm:"foreignKey:RiskActionPlanCommentID"`
}

func (RiskActionPlanComment) TableName() string {
	return "risk_action_plan_comments"
}
	