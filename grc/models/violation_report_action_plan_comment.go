package models

import (
	"time"
)

type ViolationReportActionPlanComment struct {
	ID                          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ViolationReportActionPlanID uint      `json:"violation_report_action_plan_id" gorm:"not null"`
	Comment                     string    `json:"comment" binding:"required"`
	CreatedBYID                 uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt                   time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	ViolationReportActionPlan *ViolationReportActionPlan                  `gorm:"foreignKey:ViolationReportActionPlanID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedByUser             *User                                       `gorm:"foreignKey:CreatedBYID"`
	Documents                 []*ViolationReportActionPlanCommentDocument `gorm:"foreignKey:ViolationReportActionPlanCommentID"`
}

func (ViolationReportActionPlanComment) TableName() string {
	return "violation_report_action_plan_comments"
}
