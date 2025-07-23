package models

import "time"

type StrategyProfileStrategicIntiative struct {
	ID                                  uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	StrategyProfileID                   uint       `json:"strategy_profile_id" binding:"required" gorm:"not null"`
	StrategyProfileStrategicObjectiveID *uint      `json:"strategy_profile_strategic_objective_id" gorm:"null"`
	StrategicIntiativeID                *uint      `json:"strategic_intiative_id" gorm:"null"`
	CreatedBYID                         uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID                         *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                           time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                           *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	StrategyProfile                   *StrategyProfile                   `gorm:"foreignKey:StrategyProfileID" json:"strategy_profile,omitempty"`
	StrategyProfileStrategicObjective *StrategyProfileStrategicObjective `gorm:"foreignKey:StrategyProfileStrategicObjectiveID" json:"strategy_profile_strategic_objective,omitempty"`
	StrategicIntiative                *StrategicIntiative                `gorm:"foreignKey:StrategicIntiativeID" json:"strategic_intiative,omitempty"`
	CreatedByUser                     *User                              `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	UpdatedByUser                     *User                              `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
}

func (StrategyProfileStrategicIntiative) TableName() string {
	return "strategy_profile_strategic_intiatives"
}
