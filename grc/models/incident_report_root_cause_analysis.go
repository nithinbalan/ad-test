package models

import (
	"time"
)

type IncidentReportRootCauseAnalysis struct {
	ID                       uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	IncidentReportIncidentID uint       `json:"incident_report_incident_id" binding:"required" gorm:"not null"`
	IncidentReportID         uint       `json:"incident_report_id" binding:"required" gorm:"not null"`
	Why                      uint       `json:"why" binding:"required" gorm:"not null"`
	Question                 string     `json:"question" binding:"required" gorm:"not null"`
	Answer                   string     `json:"answer,omitempty" gorm:"default:null"`
	CreatedByID              uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID              *uint      `json:"updated_by_id,omitempty" gorm:"default:null"`
	CreatedAt                time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`

	// Relationships
	IncidentReportIncident *IncidentReportIncident `gorm:"foreignKey:IncidentReportIncidentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	IncidentReport         *IncidentReport         `gorm:"foreignKey:IncidentReportID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser          *User                   `gorm:"foreignKey:UpdatedByID"`
}

func (IncidentReportRootCauseAnalysis) TableName() string {
	return "incident_report_root_cause_analyses"
}
