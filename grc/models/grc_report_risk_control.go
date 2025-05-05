package models

import "time"

type GRCReportRiskControl struct {
	ID                                   uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	GRCReportID                          uint       `json:"grc_report_id" gorm:"not null"`
	ControlCode                          string     `json:"control_code" gorm:"type:nvarchar(255);not null"`
	ControlCategoryTitleEn               string     `json:"control_category_title_en" gorm:"type:nvarchar(255);not null"`
	ControlCategoryTitleAr               string     `json:"control_category_title_ar" gorm:"type:nvarchar(255);not null"`
	Title                                string     `json:"title" gorm:"type:nvarchar(255);not null"`
	Description                          *string    `json:"description" gorm:"type:nvarchar(max)"`
	ResponsibleUserNameEn                string     `json:"responsible_user_name_en" gorm:"type:nvarchar(255);not null"`
	ResponsibleUserNameAr                string     `json:"responsible_user_name_ar" gorm:"type:nvarchar(255);not null"`
	ResponsibleUnitTitleEn               string     `json:"responsible_unit_title_en" gorm:"type:nvarchar(255);not null"`
	ResponsibleUnitTitleAr               string     `json:"responsible_unit_title_ar" gorm:"type:nvarchar(255);not null"`
	ControlImplementationStatusTitleEn   string     `json:"control_implementation_status_title_en" gorm:"type:nvarchar(255);not null"`
	ControlImplementationStatusTitleAr   string     `json:"control_implementation_status_title_ar" gorm:"type:nvarchar(255);not null"`
	ControlImplementationStatusColorCode string     `json:"control_implementation_status_color_code" gorm:"type:nvarchar(50);not null"`
	ControlEfficiencyStatusTitleEn       string     `json:"control_efficiency_status_title_en" gorm:"type:nvarchar(255);not null"`
	ControlEfficiencyStatusTitleAr       string     `json:"control_efficiency_status_title_ar" gorm:"type:nvarchar(255);not null"`
	ControlEfficiencyStatusColorCode     string     `json:"control_efficiency_status_color_code" gorm:"type:nvarchar(50);not null"`
	TestFrequencyTitleEn                 *string    `json:"test_frequency_title_en" gorm:"type:nvarchar(255)"`
	TestFrequencyTitleAr                 *string    `json:"test_frequency_title_ar" gorm:"type:nvarchar(255)"`
	NextTestDate                         *time.Time `json:"next_test_date" gorm:"type:date"`
	TestResponsibleUserNameEn            *string    `json:"test_responsible_user_name_en" gorm:"type:nvarchar(255)"`
	TestResponsibleUserNameAr            *string    `json:"test_responsible_user_name_ar" gorm:"type:nvarchar(255)"`

	GRCReport *GRCReport `gorm:"foreignKey:GRCReportID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE"`
}

func (GRCReportRiskControl) TableName() string {
	return "grc_report_risk_controls"
}
