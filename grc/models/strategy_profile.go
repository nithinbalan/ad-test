package models

import "time"

type StrategyProfile struct {
	ID                uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ReferenceNumber   string     `json:"reference_number" gorm:"type:varchar(45);not null;unique" binding:"required"`
	Title             string     `json:"title" gorm:"size:500;not null" binding:"required"`
	Description       *string    `json:"description"`
	Vision            *string    `json:"vision"`
	Mission           *string    `json:"mission"`
	ResponsibleUserID *uint      `json:"responsible_user_id"`
	StartDate         *time.Time `json:"start_date" gorm:"type:date;default:null"`
	EndDate           *time.Time `json:"end_date" gorm:"type:date;default:null"`
	CreatedBYID       uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID       *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt         time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt         *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	StrategyProfileFocusAreas  []StrategyProfileFocusArea  `gorm:"foreignKey:StrategyProfileID"`
	StrategyProfileGoals       []StrategyProfileGoal       `gorm:"foreignKey:StrategyProfileID"`
	StrategyProfileObjectives  []StrategyProfileObjective  `gorm:"foreignKey:StrategyProfileID"`
	StrategyProfileInitiatives []StrategyProfileInitiative `gorm:"foreignKey:StrategyProfileID"`
	StrategyProfileProjects    []StrategyProfileProject    `gorm:"foreignKey:StrategyProfileID"`
	StrategyProfileMilestones  []StrategyProfileMilestone  `gorm:"foreignKey:StrategyProfileID"`
	StrategyProfileActionPlans []StrategyProfileActionPlan `gorm:"foreignKey:StrategyProfileID"`

	ResponsibleUser *User `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser   *User `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	UpdatedByUser   *User `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
}

func (StrategyProfile) TableName() string {
	return "strategy_profiles"
}
