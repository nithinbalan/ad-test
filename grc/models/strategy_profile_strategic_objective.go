package models

import "time"

type StrategyProfileStrategicObjective struct {
	ID                         uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	StrategyProfileID          uint       `json:"strategy_profile_id" binding:"required" gorm:"not null"`
	StrategyProfileFocusAreaID *uint      `json:"strategy_profile_focus_area_id" gorm:"null"`
	StrategyProfileGoalID      *uint      `json:"strategy_profile_goal_id" gorm:"null"`
	StrategicObjectiveID       *uint      `json:"strategic_objective_id" gorm:"null"`
	Score                      *float64   `json:"score" gorm:"default:null"`
	CreatedBYID                uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID                *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                  time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                  *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	StrategyProfile          *StrategyProfile          `gorm:"foreignKey:StrategyProfileID" json:"strategy_profile,omitempty"`
	StrategyProfileFocusArea *StrategyProfileFocusArea `gorm:"foreignKey:StrategyProfileFocusAreaID" json:"strategy_profile_focus_area,omitempty"`
	StrategyProfileGoal      *StrategyProfileGoal      `gorm:"foreignKey:StrategyProfileGoalID" json:"strategy_profile_goal,omitempty"`
	StrategicObjective       *StrategicObjective       `gorm:"foreignKey:StrategicObjectiveID" json:"strategic_objective,omitempty"`
	// StrategyProfileInitiatives []StrategyProfileInitiative `gorm:"foreignKey:StrategyProfileObjectiveID"`
	// StrategyProfileProjects    []StrategyProfileProject    `gorm:"foreignKey:StrategyProfileObjectiveID"`
	// StrategyProfileMilestones  []StrategyProfileMilestone  `gorm:"foreignKey:StrategyProfileObjectiveID"`
	// StrategyProfileActionPlans []StrategyProfileActionPlan `gorm:"foreignKey:StrategyProfileObjectiveID"`
	CreatedByUser *User `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
}

func (StrategyProfileStrategicObjective) TableName() string {
	return "strategy_profile_strategic_objectives"
}
