package models

import "time"

type StrategyProfileFocusArea struct {
	ID                uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	StrategyProfileID uint       `json:"strategy_profile_id" gorm:"not null"`
	Title             string     `json:"title" gorm:"size:500;not null"`
	Description       *string    `json:"description"`
	CreatedBYID       uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID       *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt         time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt         *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	StrategyProfile            *StrategyProfile            `gorm:"foreignKey:StrategyProfileID" json:"strategy_profile,omitempty"`
	StrategyProfileGoals       []StrategyProfileGoal       `gorm:"foreignKey:StrategyProfileFocusAreaID"`
	StrategyProfileObjectives  []StrategyProfileObjective  `gorm:"foreignKey:StrategyProfileFocusAreaID"`
	StrategyProfileInitiatives []StrategyProfileInitiative `gorm:"foreignKey:StrategyProfileFocusAreaID"`
	StrategyProfileProjects    []StrategyProfileProject    `gorm:"foreignKey:StrategyProfileFocusAreaID"`
	StrategyProfileMilestones  []StrategyProfileMilestone  `gorm:"foreignKey:StrategyProfileFocusAreaID"`
	StrategyProfileActionPlans []StrategyProfileActionPlan `gorm:"foreignKey:StrategyProfileFocusAreaID"`

	CreatedByUser *User `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
}

func (StrategyProfileFocusArea) TableName() string {
	return "strategy_profile_focus_areas"
}
