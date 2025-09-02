package models

import (
	"time"
)

type ThreatActionPlanComment struct {
	ID                 uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ThreatActionPlanID uint      `json:"threat_action_plan_id" gorm:"not null"`
	Comment            string    `json:"comment" binding:"required"`
	CreatedBYID        uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt          time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	ThreatActionPlan *ThreatActionPlan                  `gorm:"foreignKey:ThreatActionPlanID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedByUser    *User                              `gorm:"foreignKey:CreatedBYID"`
	Documents        []*ThreatActionPlanCommentDocument `gorm:"foreignKey:ThreatActionPlanCommentID"`
}

func (ThreatActionPlanComment) TableName() string {
	return "threat_action_plan_comments"
}
