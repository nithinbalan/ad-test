package models

import "time"

type StrategyProfileProject struct {
	ID                uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	StrategyProfileID uint       `json:"strategy_profile_id" gorm:"not null"`
	ProjectID         *uint      `json:"project_id" gorm:"default:null"`
	CreatedBYID       uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID       *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt         time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt         *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	StrategyProfile *StrategyProfile `gorm:"foreignKey:StrategyProfileID" json:"strategy_profile,omitempty"`
	Project         *Project         `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	// StrategyProfileFocusArea *StrategyProfileFocusArea `gorm:"foreignKey:StrategyProfileFocusAreaID" json:"strategy_profile_focus_area,omitempty"`
	// StrategyProfileGoal      *StrategyProfileGoal      `gorm:"foreignKey:StrategyProfileGoalID" json:"strategy_profile_goal,omitempty"`
	//StrategyProfileStrategicObjective *StrategyProfileStrategicObjective `gorm:"foreignKey:StrategyProfileStrategicObjectiveID" json:"strategy_profile_strategic_objective,omitempty"`
	//StrategyProfileInitiative  *StrategyProfileInitiative  `gorm:"foreignKey:StrategyProfileInitiativeID" json:"strategy_profile_initiative,omitempty"`
	// StrategyProfileMilestones  []StrategyProfileMilestone  `gorm:"foreignKey:StrategyProfileProjectID"`
	// StrategyProfileActionPlans []StrategyProfileActionPlan `gorm:"foreignKey:StrategyProfileProjectID"`
	// ResponsibleUser            *User                       `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
}

func (StrategyProfileProject) TableName() string {
	return "strategy_profile_projects"
}
