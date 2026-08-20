package models

import "time"

type IncidentReportEvidence struct {
	ID                       uint       `json:"id" gorm:"primary_key;auto_increment"`
	IncidentEvidenceID       uint       `json:"incident_evidence_id" gorm:"not null"`
	IncidentReportIncidentID uint       `json:"incident_report_incident_id" gorm:"not null"`
	URL                      *string    `json:"url" gorm:"type:varchar(1000)"`
	Name                     *string    `json:"name" gorm:"type:nvarchar(255)"`
	Size                     *int       `json:"size" gorm:"type:int"`
	Token                    *string    `json:"token" gorm:"type:varchar(255)"`
	OriginalName             *string    `json:"original_name" gorm:"type:nvarchar(255)"`
	Type                     *string    `json:"type" gorm:"type:varchar(255)"`
	CreatedByID              *uint      `json:"created_by_id" gorm:"type:int"`
	UpdatedByID              *uint      `json:"updated_by_id" gorm:"type:int"`
	CreatedAt                time.Time  `json:"created_at" gorm:"type:datetime;not null;default:getdate()"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	IncidentEvidence       *IncidentEvidence       `gorm:"foreignKey:IncidentEvidenceID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	IncidentReportIncident *IncidentReportIncident `gorm:"foreignKey:IncidentReportIncidentID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser          *User                   `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (IncidentReportEvidence) TableName() string {
	return "incident_report_evidences"
}
