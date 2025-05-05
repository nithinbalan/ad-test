package models

type GRCReportRiskKRI struct {
	ID                       uint    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	GRCReportID              uint    `json:"grc_report_id" gorm:"not null"`
	RiskRegisterTitle        string  `json:"risk_register_title" gorm:"type:nvarchar(255);not null"`
	KRIRegisterTitle         string  `json:"kri_register_title" gorm:"type:nvarchar(255);not null"`
	KRIRegisterDescription   *string `json:"kri_register_description" gorm:"type:nvarchar(max)"`
	LowThreshold             *int    `json:"low_threshold"`
	MediumThreshold          *int    `json:"medium_threshold"`
	CurrentValue             *int    `json:"current_value"`
	RiskStatusTitleEn        string  `json:"risk_status_title_en" gorm:"type:nvarchar(255);not null"`
	RiskStatusTitleAr        string  `json:"risk_status_title_ar" gorm:"type:nvarchar(255);not null"`
	RiskStatusTitleColorCode string  `json:"risk_status_title_color_code" gorm:"type:varchar(7);not null"`
}

func (GRCReportRiskKRI) TableName() string {
	return "grc_report_risk_kris"
}
