package models

import "time"

type IncidentStrategicObjective struct {
	ID                   uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	StrategicObjectiveID uint      `json:"strategic_objective_id" gorm:"not null"`
	IncidentRegisterID   uint      `json:"incident_register_id" gorm:"not null"`
	CreatedAt            time.Time `json:"created_at" gorm:"type:datetime;not null"`
	CreatedBYID          uint      `json:"created_by_id" gorm:"not null"`

	// Relationships
	StrategicObjective *StrategicObjective `gorm:"foreignKey:StrategicObjectiveID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	IncidentRegister   *IncidentRegister   `gorm:"foreignKey:IncidentRegisterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedByUser      *User               `gorm:"foreignKey:CreatedBYID"`
}

// TableName overrides the default table name
func (IncidentStrategicObjective) TableName() string {
	return "incident_strategic_objectives"
}
