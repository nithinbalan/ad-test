package models

import "time"

type IncidentInvestigationFinding struct {
	ID                      uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	IncidentInvestigationID uint       `json:"incident_investigation_id" binding:"required" gorm:"not null"`
	Title                   string     `json:"title" binding:"required" gorm:"type:nvarchar(500);not null"`
	CreatedBy               uint       `json:"created_by" binding:"required"`
	UpdatedBy               *uint      `json:"updated_by"`
	CreatedAt               time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt               *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	IncidentInvestigation IncidentInvestigation `gorm:"foreignkey:IncidentInvestigationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser         User                  `gorm:"foreignkey:CreatedBy"`
	UpdatedByUser         *User                 `gorm:"foreignkey:UpdatedBy"`
}

// TableName overrides the table name used by IncidentInvestigationFinding to `incident_investigation_findings`
func (IncidentInvestigationFinding) TableName() string {
	return "incident_investigation_findings"
}
