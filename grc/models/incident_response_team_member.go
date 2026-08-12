package models

import (
	"time"
)

type IncidentResponseTeamMember struct {
	ID                       uint      `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	IncidentResponseTeamID   uint      `json:"incident_response_team_id" gorm:"not null"`
	UserID                   uint      `json:"user_id" gorm:"not null"`
	CreatedBYID              uint      `json:"created_by_id"`
	CreatedAt                time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	IncidentResponseTeam     *IncidentResponseTeam `gorm:"foreignKey:IncidentResponseTeamID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	User                     *User                  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser            *User                  `gorm:"foreignKey:CreatedBYID"`
}

// TableName overrides the table name used by IncidentResponseTeamMember to `incident_response_team_members`
func (IncidentResponseTeamMember) TableName() string {
	return "incident_response_team_members"
}
