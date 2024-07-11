package models

import (
	"time"
)

type IncidentInvestigation struct {
	ID                            uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	IncidentRegisterID            uint       `json:"incident_register_id" binding:"required" gorm:"not null"`
	InitiatedAt                   *time.Time `json:"initiated_at" binding:"required" gorm:"type:datetime;not null"`
	Description                   *string    `json:"description" gorm:"type:nvarchar(max)"`
	TeamLeadID                    uint       `json:"team_lead_id" binding:"required" gorm:"not null"`
	ActionRecommendation          *string    `json:"action_recommendation" gorm:"type:nvarchar(max)"`
	IncidentInvestigationStatusID int        `json:"incident_investigation_status_id" binding:"required" gorm:"not null"`
	CreatedByID                   uint       `json:"created_by_id" binding:"required"`
	UpdatedByID                   *uint      `json:"updated_by_id"`
	CreatedAt                     time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                     *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	IncidentRegister            IncidentRegister                 `gorm:"foreignKey:IncidentRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	TeamLead                    User                             `gorm:"foreignKey:TeamLeadID"`
	IncidentInvestigationStatus IncidentInvestigationStatus      `gorm:"foreignKey:IncidentInvestigationStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Documents                   []*IncidentInvestigationDocument `gorm:"foreignKey:IncidentInvestigationID"`
	Users                       []*IncidentInvestigationUser     `gorm:"foreignKey:IncidentInvestigationID"`
	Findings                    []*IncidentInvestigationFinding  `gorm:"foreignKey:IncidentInvestigationID"`
	CreatedByUser               User                             `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser               *User                            `gorm:"foreignKey:UpdatedByID"`
}

// TableName overrides the table name used by IncidentInvestigation to `incident_investigations`
func (IncidentInvestigation) TableName() string {
	return "incident_investigations"
}
