package models

import "time"

type StrategyProfileGoal struct {
	ID                         uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	StrategyProfileID          uint       `json:"strategy_profile_id" binding:"required" gorm:"not null"`
	StrategyProfileFocusAreaID *uint      `json:"strategy_profile_focus_area_id" gorm:"null"`
	Title                      string     `json:"title" gorm:"size:500;not null" binding:"required"`
	Description                *string    `json:"description" binding:"omitempty"`
	TargetDate                 *time.Time `json:"target_date" gorm:"type:date;default:null"`
	ResponsibleUserID          *uint      `json:"responsible_user_id"`
	CreatedBYID                uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID                *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                  time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                  *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	StrategyProfile          *StrategyProfile          `gorm:"foreignKey:StrategyProfileID" json:"strategy_profile,omitempty"`
	StrategyProfileFocusArea *StrategyProfileFocusArea `gorm:"foreignKey:StrategyProfileFocusAreaID" json:"strategy_profile_focus_area,omitempty"`
	ResponsibleUser          *User                     `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser            *User                     `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	UpdatedByUser            *User                     `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
}

func (StrategyProfileGoal) TableName() string {
	return "strategy_profile_goals"
}
