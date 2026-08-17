package models

import (
	"time"
)

type IncidentRegisterResponseTeam struct {
	ID                 uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	IncidentRegisterID uint      `json:"incident_register_id" binding:"required"`
	UserID             uint      `json:"user_id" binding:"required"`
	CreatedByID        uint      `json:"created_by_id" binding:"required"`
	CreatedAt          time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	IncidentRegister *IncidentRegister `gorm:"foreignKey:IncidentRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	User             *User             `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser    *User             `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the table name used by IncidentRegisterResponseTeam to `incident_register_response_teams`
func (IncidentRegisterResponseTeam) TableName() string {
	return "incident_register_response_teams"
}
