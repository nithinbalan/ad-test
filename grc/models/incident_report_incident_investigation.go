package models

import "time"

type IncidentReportIncidentInvestigation struct {
	ID                       uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	IncidentReportID         uint       `json:"incident_report_id" gorm:"not null"`
	IncidentReportIncidentID uint       `json:"incident_report_incident_id" gorm:"not null"`
	IncidentInvestigationID  *uint      `json:"incident_investigation_id" gorm:"type:int"`
	InitiatedAt              *time.Time `json:"initiated_at" gorm:"type:datetime"`
	Description              *string    `json:"description" gorm:"type:nvarchar(max)"`
	TeamLeadID               *uint      `json:"team_lead_id" gorm:"type:int"`
	TeamLeadFirstNameEn      *string    `json:"team_lead_first_name_en" gorm:"type:nvarchar(255)"`
	TeamLeadFirstNameAr      *string    `json:"team_lead_first_name_ar" gorm:"type:nvarchar(255)"`
	TeamLeadLastNameEn       *string    `json:"team_lead_last_name_en" gorm:"type:nvarchar(255)"`
	TeamLeadLastNameAr       *string    `json:"team_lead_last_name_ar" gorm:"type:nvarchar(255)"`
	ActionRecommendation     *string    `json:"action_recommendation" gorm:"type:nvarchar(max)"`
	CreatedAt                time.Time  `json:"created_at" gorm:"type:datetime;not null;default:getdate()"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	IncidentReport         *IncidentReport         `gorm:"foreignKey:IncidentReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	IncidentReportIncident *IncidentReportIncident `gorm:"foreignKey:IncidentReportIncidentID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	IncidentInvestigation  *IncidentInvestigation  `gorm:"foreignKey:IncidentInvestigationID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	TeamLead               *User                   `gorm:"foreignKey:TeamLeadID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (IncidentReportIncidentInvestigation) TableName() string {
	return "incident_report_incident_investigations"
}
