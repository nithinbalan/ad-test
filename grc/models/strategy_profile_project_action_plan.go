package models

import "time"

type StrategyProfileProjectActionPlan struct {
	ID                                uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	StrategyProfileID                 uint       `json:"strategy_profile_id" binding:"required" gorm:"not null"`
	StrategyProfileProjectID          *uint      `json:"strategy_profile_project_id" gorm:"null"`
	StrategyProfileProjectMilestoneID *uint      `json:"strategy_profile_project_milestone_id" gorm:"null"`
	ProjectActionPlanID               *uint      `json:"project_action_plan_id" gorm:"null"`
	CreatedBYID                       uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID                       *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                         time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                         *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	StrategyProfile                 *StrategyProfile                 `gorm:"foreignKey:StrategyProfileID" json:"strategy_profile,omitempty"`
	StrategyProfileProject          *StrategyProfileProject          `gorm:"foreignKey:StrategyProfileProjectID" json:"strategy_profile_project,omitempty"`
	StrategyProfileProjectMilestone *StrategyProfileProjectMilestone `gorm:"foreignKey:StrategyProfileProjectMilestoneID" json:"strategy_profile_project_milestone,omitempty"`
	ProjectActionPlan               *ProjectActionPlan               `gorm:"foreignKey:ProjectActionPlanID" json:"project_action_plan,omitempty"`
	CreatedByUser                   *User                            `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	UpdatedByUser                   *User                            `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
}

func (StrategyProfileProjectActionPlan) TableName() string {
	return "strategy_profile_project_action_plans"
}
