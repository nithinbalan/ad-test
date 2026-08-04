package models

import "time"

type RiskReportKRI struct {
	ID                           uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskReportID                 uint      `json:"risk_report_id" gorm:"not null"`
	RiskReportIdentificationID   uint      `json:"risk_report_identification_id" gorm:"not null"`
	Title                        string    `json:"title" gorm:"type:nvarchar(500);not null"`
	Description                  *string   `json:"description" gorm:"type:nvarchar(max)"`
	AnnualTarget                 *float64  `json:"annual_target" gorm:"type:decimal(18,2);null"`
	MinimumSusceptibilityPercent *float64  `json:"minimum_susceptibility_percent" gorm:"type:decimal(18,2);null"`
	MaximumSusceptibilityPercent *float64  `json:"maximum_susceptibility_percent" gorm:"type:decimal(18,2);null"`
	MaximumTolerablePercent      *float64  `json:"maximum_tolerable_percent" gorm:"type:decimal(18,2);null"`
	IsLowerBetter                *bool     `json:"is_lower_better" gorm:"type:bit;default:0"`
	IsHigherBetter               *bool     `json:"is_higher_better" gorm:"type:bit;default:0"`
	CurrentValue                 *float64  `json:"current_value" gorm:"type:decimal(18,2);null"`
	UnitTitleEn                  *string   `json:"unit_title_en" gorm:"type:nvarchar(255);null"`
	UnitTitleAr                  *string   `json:"unit_title_ar" gorm:"type:nvarchar(255);null"`
	CreatedAt                    time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	RiskReportIdentification *RiskReportIdentification `gorm:"foreignKey:RiskReportIdentificationID"`
	RiskReport               *RiskReport               `gorm:"foreignKey:RiskReportID"`
}

func (RiskReportKRI) TableName() string {
	return "risk_report_kris"
}
