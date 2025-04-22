package models

import "time"

type StrategyProfileMilestone struct {
	ID                          uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	StrategyProfileID           uint       `json:"strategy_profile_id" binding:"required" gorm:"not null"`
	StrategyProfileFocusAreaID  *uint      `json:"strategy_profile_focus_area_id" gorm:"null"`
	StrategyProfileGoalID       *uint      `json:"strategy_profile_goal_id" gorm:"null"`
	StrategyProfileObjectiveID  *uint      `json:"strategy_profile_objective_id" gorm:"null"`
	StrategyProfileInitiativeID *uint      `json:"strategy_profile_initiative_id" gorm:"null"`
	StrategyProfileProjectID    *uint      `json:"strategy_profile_project_id" gorm:"null"`
	Title                       string     `json:"title" gorm:"size:500;not null" binding:"required"`
	CreatedBYID                 uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID                 *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                   *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	StrategyProfile           *StrategyProfile           `gorm:"foreignKey:StrategyProfileID" json:"strategy_profile,omitempty"`
	StrategyProfileFocusArea  *StrategyProfileFocusArea  `gorm:"foreignKey:StrategyProfileFocusAreaID" json:"strategy_profile_focus_area,omitempty"`
	StrategyProfileGoal       *StrategyProfileGoal       `gorm:"foreignKey:StrategyProfileGoalID" json:"strategy_profile_goal,omitempty"`
	StrategyProfileObjective  *StrategyProfileObjective  `gorm:"foreignKey:StrategyProfileObjectiveID" json:"strategy_profile_objective,omitempty"`
	StrategyProfileInitiative *StrategyProfileInitiative `gorm:"foreignKey:StrategyProfileInitiativeID" json:"strategy_profile_initiative,omitempty"`
	StrategyProfileProject    *StrategyProfileProject    `gorm:"foreignKey:StrategyProfileProjectID" json:"strategy_profile_project,omitempty"`
	CreatedByUser             *User                      `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	UpdatedByUser             *User                      `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
}

func (StrategyProfileMilestone) TableName() string {
	return "strategy_profile_milestones"
}
