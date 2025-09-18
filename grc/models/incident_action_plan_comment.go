package models

import (
	"time"
)

type IncidentActionPlanComment struct {
	ID                     uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	IncidentActionPlanID uint      `json:"incident_action_plan_id" gorm:"not null"`
	Comment                string    `json:"comment" binding:"required"`
	CreatedBYID            uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt              time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	IncidentActionPlan *IncidentActionPlan                  `gorm:"foreignKey:IncidentActionPlanID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedByUser        *User                                  `gorm:"foreignKey:CreatedBYID"`
	Documents            []*IncidentActionPlanCommentDocument `gorm:"foreignKey:IncidentActionPlanCommentID"`
}

func (IncidentActionPlanComment) TableName() string {
	return "incident_action_plan_comments"
}
