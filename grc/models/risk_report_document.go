package models

import (
	"time"
)

type RiskReportDocument struct {
	ID                         uint       `json:"id" gorm:"primary_key;auto_increment"`
	RiskReportID               uint       `json:"risk_register_id"`
	RiskReportIdentificationID uint       `json:"risk_report_identification_id" gorm:"not null"`
	URL                        string     `json:"url"`
	Name                       string     `json:"name" binding:"required"`
	Size                       int        `json:"size"`
	Token                      string     `json:"token"`
	OriginalName               string     `json:"original_name"`
	Type                       string     `json:"type"`
	CreatedByID                uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID                *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                  time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt                  *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`

	// Relationships
	RiskReport               *RiskReport               `gorm:"foreignKey:RiskReportID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskReportIdentification *RiskReportIdentification `gorm:"foreignKey:RiskReportIdentificationID"`
	CreatedByUser            *User                     `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser            *User                     `gorm:"foreignKey:UpdatedByID"`
}

func (RiskReportDocument) TableName() string {
	return "risk_report_documents"
}
