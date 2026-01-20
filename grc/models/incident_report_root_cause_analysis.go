package models

type IncidentReportRootCauseAnalysis struct {
	ID                       uint    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	IncidentReportIncidentID uint    `json:"incident_report_incident_id" binding:"required" gorm:"not null"`
	IncidentReportID         uint    `json:"incident_report_id" binding:"required" gorm:"not null"`
	Why                      uint    `json:"why" binding:"required" gorm:"not null"`
	Question                 string  `json:"question" binding:"required" gorm:"type:nvarchar(500);not null"`
	Answer                   *string `json:"answer" gorm:"type:nvarchar(1000)"`

	// Relationships
	IncidentReportIncident *IncidentReportIncident `gorm:"foreignKey:IncidentReportIncidentID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	IncidentReport         *IncidentReport         `gorm:"foreignKey:IncidentReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (IncidentReportRootCauseAnalysis) TableName() string {
	return "incident_report_root_cause_analyses"
}
