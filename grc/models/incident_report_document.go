package models

import "time"

type IncidentReportDocument struct {
	ID                       uint       `json:"id" gorm:"primary_key;auto_increment"`
	IncidentDocumentID       uint       `json:"incident_document_id" gorm:"not null"`
	IncidentReportIncidentID uint       `json:"incident_report_incident_id" gorm:"not null"`
	IncidentInvestigationID  *uint      `json:"incident_investigation_id" gorm:"type:int"`
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
	IncidentDocument       *IncidentDocument       `gorm:"foreignKey:IncidentDocumentID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	IncidentReportIncident *IncidentReportIncident `gorm:"foreignKey:IncidentReportIncidentID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	IncidentInvestigation  *IncidentInvestigation  `gorm:"foreignKey:IncidentInvestigationID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser          *User                   `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (IncidentReportDocument) TableName() string {
	return "incident_report_documents"
}
