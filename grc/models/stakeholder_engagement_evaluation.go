package models

import (
	"time"
)

type StakeholderEngagementEvaluation struct {
	ID                      uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	StakeholderEngagementID uint       `json:"stakeholder_engagement_id" binding:"required" gorm:"not null"`
	Score                   *int       `json:"score" binding:"omitempty,min=0,max=100" gorm:"not null;default:0"`
	Comment                 *string    `json:"comment" binding:"required" gorm:"type:nvarchar(max);not null"`
	EvaluatedAt             *time.Time `json:"evaluated_at" binding:"omitempty" gorm:"type:date,default:null"`
	EvaluatedBYID           uint       `json:"evaluated_by_id"`
	CreatedAt               time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt               *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`
	CreatedBYID             uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID             *uint      `json:"updated_by_id,omitempty" gorm:"default:null"`

	// Relationships
	StakeholderEngagement *StakeholderEngagement `gorm:"foreignKey:StakeholderEngagementID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser         *User                  `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser         *User                  `gorm:"foreignKey:UpdatedBYID"`
	EvaluatedByUser       *User                  `gorm:"foreignKey:EvaluatedBYID"`
}

// TableName overrides the table name used by StakeholderNeedsAndExpectation to `stakeholder_engagement_evaluations`
func (StakeholderEngagementEvaluation) TableName() string {
	return "stakeholder_engagement_evaluations"
}
