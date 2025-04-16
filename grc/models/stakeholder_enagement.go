package models

import (
	"time"
)

type StakeholderEngagement struct {
	ID                 uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	StakeholderID      uint       `json:"stakeholder_id" binding:"required" gorm:"not null"`
	Title              string     `json:"title" binding:"required,max=500" gorm:"type:varchar(500);not null"`
	ReviewFrequencyID  *uint      `json:"review_frequency_id"`
	AssignedOwnerID    *uint      `json:"assigned_owner_id"`
	EngagementStrategy *string    `json:"engagement_strategy"`
	Description        *string    `json:"description"`
	Score              *int       `json:"score" binding:"omitempty,min=0,max=100" gorm:"not null;default:0"`
	EvaluatedAt        *time.Time `json:"evaluated_at" binding:"omitempty" gorm:"type:date,default:null"`
	EvaluatedBYID      *uint      `json:"evaluated_by_id"`
	CreatedAt          time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`
	CreatedBYID        uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID        *uint      `json:"updated_by_id,omitempty" gorm:"default:null"`

	// Relationships
	Stakeholder     *Stakeholder     `gorm:"foreignKey:StakeholderID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ReviewFrequency *ReviewFrequency `gorm:"foreignKey:ReviewFrequencyID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser   *User            `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser   *User            `gorm:"foreignKey:UpdatedBYID"`
	EvaluatedByUser *User            `gorm:"foreignKey:EvaluatedBYID"`
	AssignedOwner   *User            `gorm:"foreignKey:AssignedOwnerID"`
}

// TableName overrides the table name used by StakeholderNeedsAndExpectation to `stakeholder_engagements`
func (StakeholderEngagement) TableName() string {
	return "stakeholder_engagements"
}
