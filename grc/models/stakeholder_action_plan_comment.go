package models

import (
	"time"
)

type StakeholderActionPlanComment struct {
	ID                      uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	StakeholderActionPlanID uint      `json:"stakeholder_action_plan_id" gorm:"not null"`
	Comment                 string    `json:"comment" binding:"required"`
	CreatedBYID             uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt               time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	StakeholderActionPlan *StakeholderActionPlan                  `gorm:"foreignKey:StakeholderActionPlanID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedByUser         *User                                   `gorm:"foreignKey:CreatedBYID"`
	Documents             []*StakeholderActionPlanCommentDocument `gorm:"foreignKey:StakeholderActionPlanCommentID"`
}

func (StakeholderActionPlanComment) TableName() string {
	return "stakeholder_action_plan_comments"
}
