package models

import (
	"time"
)

type StrategyProfileActionPlanComment struct {
	ID                          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	StrategyProfileActionPlanID uint      `json:"strategy_profile_action_plan_id" gorm:"not null"`
	Comment                     string    `json:"comment" binding:"required"`
	CreatedBYID                 uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt                   time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	StrategyProfileActionPlan *StrategyProfileActionPlan                  `gorm:"foreignKey:StrategyProfileActionPlanID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedByUser             *User                                       `gorm:"foreignKey:CreatedBYID"`
	Documents                 []*StrategyProfileActionPlanCommentDocument `gorm:"foreignKey:StrategyProfileActionPlanCommentID"`
}

func (StrategyProfileActionPlanComment) TableName() string {
	return "strategy_profile_action_plan_comments"
}
