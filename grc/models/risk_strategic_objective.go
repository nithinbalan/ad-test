package models

import (
	"time"
)

type RiskStrategicObjective struct {
	ID                   uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskRegisterID       uint      `json:"risk_register_id" binding:"required"`
	StrategicObjectiveID uint      `json:"strategic_objective_id" binding:"required"`
	CreatedByID          uint      `json:"created_by_id" binding:"required"`
	CreatedAt            time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	RiskRegister       *RiskRegister       `gorm:"foreignKey:RiskRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	StrategicObjective *StrategicObjective `gorm:"foreignKey:StrategicObjectiveID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser      *User               `gorm:"foreignKey:CreatedByID"`
}

func (RiskStrategicObjective) TableName() string {
	return "risk_strategic_objectives"
}
