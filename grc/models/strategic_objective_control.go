package models

import (
	"time"
)

type StrategicObjectiveControl struct {
	ID                   uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	StrategicObjectiveID uint      `json:"strategic_objective_id" binding:"required"`
	ControlID            uint      `json:"control_id" binding:"required"`
	CreatedBYID          uint      `json:"created_by_id" binding:"required"`
	CreatedAt            time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	StrategicObjective *StrategicObjective `gorm:"foreignKey:StrategicObjectiveID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Control            *Control            `gorm:"foreignKey:ControlID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser      *User               `gorm:"foreignKey:CreatedBYID"`
}

func (StrategicObjectiveControl) TableName() string {
	return "strategic_objective_controls"
}
