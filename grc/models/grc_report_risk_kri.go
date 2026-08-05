package models

import "time"

type GRCReportRiskKRI struct {
	ID                           uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	GRCReportID                  uint      `json:"grc_report_id" gorm:"not null"`
	RiskRegisterTitle            string    `json:"risk_register_title" gorm:"type:nvarchar(255);not null"`
	KRIRegisterTitle             string    `json:"kri_register_title" gorm:"type:nvarchar(255);not null"`
	KRIRegisterDescription       *string   `json:"kri_register_description" gorm:"type:nvarchar(max)"`
	AnnualTarget                 *float64  `json:"annual_target" gorm:"type:decimal(18,2);null"`
	MinimumSusceptibilityPercent *float64  `json:"minimum_susceptibility_percent" gorm:"type:decimal(18,2);null"`
	MaximumSusceptibilityPercent *float64  `json:"maximum_susceptibility_percent" gorm:"type:decimal(18,2);null"`
	MaximumTolerablePercent      *float64  `json:"maximum_tolerable_percent" gorm:"type:decimal(18,2);null"`
	IsLowerBetter                *bool     `json:"is_lower_better" gorm:"type:bit;default:0"`
	IsHigherBetter               *bool     `json:"is_higher_better" gorm:"type:bit;default:0"`
	CurrentValue                 *float64  `json:"current_value" gorm:"type:decimal(18,2);null"`
	RiskStatusTitleEn            string    `json:"risk_status_title_en" gorm:"type:nvarchar(255);not null"`
	RiskStatusTitleAr            string    `json:"risk_status_title_ar" gorm:"type:nvarchar(255);not null"`
	RiskStatusTitleColorCode     string    `json:"risk_status_title_color_code" gorm:"type:varchar(7);not null"`
	CreatedBYID                  uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt                    time.Time `json:"created_at" gorm:"type:datetime;not null"`

	CreatedByUser *User `gorm:"foreignKey:CreatedBYID"`
}

func (GRCReportRiskKRI) TableName() string {
	return "grc_report_risk_kris"
}
