package models

import "time"

type StrategyProfileProject struct {
	ID                                  uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	StrategyProfileID                   uint       `json:"strategy_profile_id" gorm:"not null"`
	StrategyProfileStrategicObjectiveID *uint      `json:"strategy_profile_strategic_objective_id" gorm:"default:null"`
	StrategyProfileStrategicIntiativeID *uint      `json:"strategy_profile_strategic_intiative_id" gorm:"default:null"`
	ProjectID                           *uint      `json:"project_id" gorm:"default:null"`
	CreatedBYID                         uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID                         *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                           time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                           *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	StrategyProfile                   *StrategyProfile                   `gorm:"foreignKey:StrategyProfileID" json:"strategy_profile,omitempty"`
	Project                           *Project                           `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	StrategyProfileStrategicObjective *StrategyProfileStrategicObjective `gorm:"foreignKey:StrategyProfileStrategicObjectiveID" json:"strategy_profile_strategic_objective,omitempty"`
	StrategyProfileStrategicIntiative *StrategyProfileStrategicIntiative `gorm:"foreignKey:StrategyProfileStrategicIntiativeID" json:"strategy_profile_strategic_intiative,omitempty"`
	StrategyProfileProjectMilestones  []StrategyProfileProjectMilestone  `gorm:"foreignKey:StrategyProfileProjectID"`
	StrategyProfileProjectActionPlans []StrategyProfileProjectActionPlan `gorm:"foreignKey:StrategyProfileProjectID"`
	CreatedByUser                     *User                              `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	UpdatedByUser                     *User                              `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
}

func (StrategyProfileProject) TableName() string {
	return "strategy_profile_projects"
}
