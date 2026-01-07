package models

import "time"

type IncidentReportActionPlan struct {
	ID                       uint `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	IncidentReportID         uint `json:"incident_report_id" gorm:"not null"`
	IncidentReportIncidentID uint `json:"incident_report_incident_id" gorm:"not null"`

	Title       string     `json:"title" gorm:"type:nvarchar(255);not null"`
	Description *string     `json:"description" gorm:"type:nvarchar(max)"`
	TargetDate  *time.Time `json:"target_date"`
	StatusTitle string     `json:"status_title" gorm:"type:nvarchar(255)"`

	CreatedByID *uint     `json:"created_by_id"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null;default:getdate()"`

	// Relations
	IncidentReport         *IncidentReport         `gorm:"foreignKey:IncidentReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	IncidentReportIncident *IncidentReportIncident `gorm:"foreignKey:IncidentReportIncidentID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedBy              *User                   `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (IncidentReportActionPlan) TableName() string {
	return "incident_report_action_plans"
}
