package models

import (
	"time"
)

type StrategicObjectiveProcess struct {
	ID                   uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ProcessID            uint      `json:"process_id" binding:"required"`
	StrategicObjectiveID uint      `json:"strategic_objective_id" binding:"required"`
	CreatedByID          uint      `json:"created_by_id" binding:"required"`
	CreatedAt            time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	Process            *Process            `gorm:"foreignKey:ProcessID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	StrategicObjective *StrategicObjective `gorm:"foreignKey:StrategicObjectiveID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser      *User               `gorm:"foreignKey:CreatedByID"`
}

func (StrategicObjectiveProcess) TableName() string {
	return "strategic_objective_processes"
}
