package models

import "time"

type StrategyProfileKpi struct {
	ID                                  uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	StrategyProfileID                   uint       `json:"strategy_profile_id" binding:"required" gorm:"not null"`
	StrategyProfileStrategicObjectiveID *uint      `json:"strategy_profile_strategic_objective_id" gorm:"null"`
	KpiID                               *uint      `json:"kpi_id" gorm:"null"`
	Weightage                           *float64   `json:"weightage" gorm:"default:null"`
	CreatedBYID                         uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID                         *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                           time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                           *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	StrategyProfile                   *StrategyProfile                   `gorm:"foreignKey:StrategyProfileID" json:"strategy_profile,omitempty"`
	StrategyProfileStrategicObjective *StrategyProfileStrategicObjective `gorm:"foreignKey:StrategyProfileStrategicObjectiveID" json:"strategy_profile_strategic_objective,omitempty"`
	Kpi                               *KPIRegister                       `gorm:"foreignKey:KpiID" json:"kpi,omitempty"`
	CreatedByUser                     *User                              `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	UpdatedByUser                     *User                              `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
}

func (StrategyProfileKpi) TableName() string {
	return "strategy_profile_kpis"
}
