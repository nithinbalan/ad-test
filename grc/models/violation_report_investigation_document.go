package models

import (
	"time"
)

type ViolationReportInvestigationDocument struct {
	ID                             uint      `json:"id" gorm:"primary_key;auto_increment"`
	ViolationReportInvestigationID uint      `json:"violation_report_investigation_id"`
	URL                            string    `json:"url"`
	Name                           string    `json:"name" binding:"required"`
	Size                           int       `json:"size"`
	Token                          string    `json:"token"`
	OriginalName                   string    `json:"original_name"`
	Type                           string    `json:"type"`
	CreatedByID                    uint      `json:"created_by_id" gorm:"not null"`
	UpdatedByID                    uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                      time.Time `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                      time.Time `json:"updated_at" gorm:"type:datetime;not null"`

	// Relationships
	ViolationReportInvestigation *ViolationReportInvestigation `gorm:"foreignKey:ViolationReportInvestigationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser                *User                         `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser                *User                         `gorm:"foreignKey:UpdatedByID"`
}

func (ViolationReportInvestigationDocument) TableName() string {
	return "violation_report_investigation_documents"
}
