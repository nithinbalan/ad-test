package models

import "time"

type StrategyProfileProjectMilestone struct {
	ID                       uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	StrategyProfileID        uint       `json:"strategy_profile_id" binding:"required" gorm:"not null"`
	StrategyProfileProjectID *uint      `json:"strategy_profile_project_id" gorm:"null"`
	ProjectMilestoneID       *uint      `json:"project_milestone_id" gorm:"null"`
	CreatedBYID              uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID              *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	StrategyProfile                   *StrategyProfile                   `gorm:"foreignKey:StrategyProfileID" json:"strategy_profile,omitempty"`
	StrategyProfileProject            *StrategyProfileProject            `gorm:"foreignKey:StrategyProfileProjectID" json:"strategy_profile_project,omitempty"`
	ProjectMilestone                  *ProjectMilestone                  `gorm:"foreignKey:ProjectMilestoneID" json:"project_milestone,omitempty"`
	StrategyProfileProjectActionPlans []StrategyProfileProjectActionPlan `gorm:"foreignKey:StrategyProfileProjectMilestoneID"`
	CreatedByUser                     *User                              `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	UpdatedByUser                     *User                              `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
}

func (StrategyProfileProjectMilestone) TableName() string {
	return "strategy_profile_project_milestones"
}
