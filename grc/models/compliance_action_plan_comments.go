package models

import (
	"time"
)

type ComplianceActionPlanComment struct {
	ID                     uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ComplianceActionPlanID uint      `json:"compliance_action_plan_id" gorm:"not null"`
	Comment                string    `json:"comment" binding:"required"`
	CreatedBYID            uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt              time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	ComplianceActionPlan *ComplianceActionPlan                  `gorm:"foreignKey:ComplianceActionPlanID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedByUser        *User                                  `gorm:"foreignKey:CreatedBYID"`
	Documents            []*ComplianceActionPlanCommentDocument `gorm:"foreignKey:ComplianceActionPlanCommentID"`
}

func (ComplianceActionPlanComment) TableName() string {
	return "compliance_action_plan_comments"
}
