package models

import "time"

type RiskReportKRI struct {
	ID                         uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskReportID               uint      `json:"risk_report_id" gorm:"not null"`
	RiskReportIdentificationID uint      `json:"risk_report_identification_id" gorm:"not null"`
	Title                      string    `json:"title" gorm:"type:nvarchar(500);not null"`
	Description                *string   `json:"description" gorm:"type:nvarchar(max)"`
	LowThreshold               *int      `json:"low_threshold" gorm:"type:int;null"`
	MediumThreshold            *int      `json:"medium_threshold" gorm:"type:int;null"`
	CurrentValue               *int      `json:"current_value" gorm:"type:int;null"`
	UnitTitleEn                *string   `json:"unit_title_en" gorm:"type:nvarchar(255);null"`
	UnitTitleAr                *string   `json:"unit_title_ar" gorm:"type:nvarchar(255);null"`
	CreatedAt                  time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	RiskReportIdentification *RiskReportIdentification `gorm:"foreignKey:RiskReportIdentificationID"`
	RiskReport               *RiskReport               `gorm:"foreignKey:RiskReportID"`
}

func (RiskReportKRI) TableName() string {
	return "risk_report_kris"
}
