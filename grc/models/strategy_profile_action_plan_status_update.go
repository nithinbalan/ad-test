package models

import (
	"time"
)

type StrategyProfileActionPlanStatusUpdate struct {
	ID                          uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	StrategyProfileActionPlanID uint       `json:"strategy_profile_action_plan_id" gorm:"not null"`
	ActionPlanStatusID          uint       `json:"action_plan_status_id" gorm:"not null"`
	Percentage                  *int       `json:"percentage" binding:"omitempty,min=0,max=100" gorm:"not null;default:0"`
	Comment                     string     `json:"comment" binding:"required" gorm:"type:nvarchar(max);not null"`
	CreatedAt                   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                   *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedByID                 uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID                 *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	StrategyProfileActionPlan *StrategyProfileActionPlan                       `gorm:"foreignKey:StrategyProfileActionPlanID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ActionPlanStatus          *ActionPlanStatus                                `gorm:"foreignKey:ActionPlanStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser             *User                                            `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser             *User                                            `gorm:"foreignKey:UpdatedByID"`
	Documents                 []*StrategyProfileActionPlanStatusUpdateDocument `gorm:"foreignKey:StrategyProfileActionPlanStatusUpdateID"`
}

// TableName overrides the table name used by StrategyProfileActionPlanStatusUpdate to `sla_action_plan_status_updates`
func (StrategyProfileActionPlanStatusUpdate) TableName() string {
	return "strategy_profile_action_plan_status_updates"
}
