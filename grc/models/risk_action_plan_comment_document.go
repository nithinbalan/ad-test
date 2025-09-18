package models

import (
	"time"
)

type RiskActionPlanCommentDocument struct {
	ID                        uint      `json:"id" gorm:"primary_key;auto_increment"`
	RiskActionPlanCommentID   uint      `json:"risk_action_plan_comment_id" gorm:"not null"`
	URL                       string    `json:"url"`
	Name                      string    `json:"name" binding:"required"`
	Size                      int       `json:"size"`
	Token                     string    `json:"token"`
	OriginalName              string    `json:"original_name"`
	Type                      string    `json:"type"`
	CreatedByID               uint      `json:"created_by_id" gorm:"not null"`
	UpdatedByID               *uint     `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                 time.Time `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                 time.Time `json:"updated_at" gorm:"type:datetime;not null"`

	// Relationships
	RiskActionPlanComment *RiskActionPlanComment `gorm:"foreignKey:RiskActionPlanCommentID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE;"`
	CreatedByUser         *User                  `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser         *User                  `gorm:"foreignKey:UpdatedByID"`
}

func (RiskActionPlanCommentDocument) TableName() string {
	return "risk_action_plan_comment_documents"
}
