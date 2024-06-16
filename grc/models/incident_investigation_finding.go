package models

import "time"

type IncidentInvestigationFinding struct {
	ID                      uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	IncidentInvestigationID uint       `json:"incident_investigation_id" binding:"required" gorm:"not null"`
	Title                   string     `json:"title" binding:"required" gorm:"type:nvarchar(500);not null"`
	CreatedByID             uint       `json:"created_by_id" binding:"required"`
	UpdatedBy               *uint      `json:"updated_by_id"`
	CreatedAt               time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt               *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	IncidentInvestigation IncidentInvestigation `gorm:"foreignKey:IncidentInvestigationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser         User                  `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser         *User                 `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by IncidentInvestigationFinding to `incident_investigation_findings`
func (IncidentInvestigationFinding) TableName() string {
	return "incident_investigation_findings"
}
