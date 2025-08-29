package models

import (
	"time"
)

type ProjectIncident struct {
	ID                 uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ProjectID          uint      `json:"project_id" binding:"required"`
	IncidentRegisterID uint      `json:"incident_register_id" binding:"required"`
	CreatedBYID        uint      `json:"created_by_id" binding:"required"`
	CreatedAt          time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	Project          *Project          `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	IncidentRegister *IncidentRegister `gorm:"foreignKey:IncidentRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser    *User             `gorm:"foreignKey:CreatedBYID"`
}

func (ProjectIncident) TableName() string {
	return "project_incidents"
}
