package models

import "time"

type StrategyProfileActionPlan struct {
	ID                          uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	StrategyProfileID           uint       `json:"strategy_profile_id" binding:"required" gorm:"not null"`
	StrategyProfileFocusAreaID  *uint      `json:"strategy_profile_focus_area_id" gorm:"null"`
	StrategyProfileGoalID       *uint      `json:"strategy_profile_goal_id" gorm:"null"`
	StrategyProfileObjectiveID  *uint      `json:"strategy_profile_objective_id" gorm:"null"`
	StrategyProfileInitiativeID *uint      `json:"strategy_profile_initiative_id" gorm:"null"`
	StrategyProfileProjectID    *uint      `json:"strategy_profile_project_id" gorm:"null"`
	StrategyProfileMilestoneID  *uint      `json:"strategy_profile_milestone_id" gorm:"null"`
	SeverityLevelID             *uint      `json:"severity_level_id"`
	Title                       string     `json:"title"`
	Description                 *string    `json:"description" gorm:"type:nvarchar(max);null"`
	ResponsibleUserID           *uint      `json:"responsible_user_id" gorm:"not null" binding:"required"`
	ResponsibleUnitID           *uint      `json:"responsible_unit_id" gorm:"not null" binding:"required"`
	TargetDate                  *time.Time `json:"target_date"`
	ActionPlanStatusId          *uint      `json:"action_plan_status_id"`
	Percentage                  *int       `json:"percentage"`
	ReviewUserID                *uint      `json:"review_user_id" gorm:"default:null"`
	CreatedBYID                 uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID                 *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                   *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	StrategyProfile           *StrategyProfile                         `gorm:"foreignKey:StrategyProfileID" json:"strategy_profile,omitempty"`
	StrategyProfileFocusArea  *StrategyProfileFocusArea                `gorm:"foreignKey:StrategyProfileFocusAreaID" json:"strategy_profile_focus_area,omitempty"`
	StrategyProfileGoal       *StrategyProfileGoal                     `gorm:"foreignKey:StrategyProfileGoalID" json:"strategy_profile_goal,omitempty"`
	StrategyProfileObjective  *StrategyProfileObjective                `gorm:"foreignKey:StrategyProfileObjectiveID" json:"strategy_profile_objective,omitempty"`
	StrategyProfileInitiative *StrategyProfileInitiative               `gorm:"foreignKey:StrategyProfileInitiativeID" json:"strategy_profile_initiative,omitempty"`
	StrategyProfileProject    *StrategyProfileProject                  `gorm:"foreignKey:StrategyProfileProjectID" json:"strategy_profile_project,omitempty"`
	StrategyProfileMilestone  *StrategyProfileMilestone                `gorm:"foreignKey:StrategyProfileMilestoneID" json:"strategy_profile_milestone,omitempty"`
	ResponsibleUser           *User                                    `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SeverityLevel             *SeverityLevel                           `gorm:"foreignKey:SeverityLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ActionPlanStatus          *ActionPlanStatus                        `gorm:"foreignKey:ActionPlanStatusId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUnit           *OrganizationUnit                        `gorm:"foreignKey:ResponsibleUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ReviewUser                *User                                    `gorm:"foreignKey:ReviewUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"` // New foreign key relationship
	ActionPlanStatusUpdates   []*StrategyProfileActionPlanStatusUpdate `gorm:"foreignKey:StrategyProfileActionPlanID"`
	CreatedByUser             *User                                    `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	UpdatedByUser             *User                                    `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
}

func (StrategyProfileActionPlan) TableName() string {
	return "strategy_profile_action_plans"
}
