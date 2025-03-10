package models

import "time"

type RiskReportActionPlan struct {
	ID                         uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskReportID               uint      `json:"risk_report_id" gorm:"not null"`
	RiskReportIdentificationID uint      `json:"risk_report_identification_id" gorm:"not null"`
	ActionPlanTitle            string    `json:"action_plan_title" gorm:"type:varchar(255);not null"`
	ActionPlanDescription      *string   `json:"action_plan_description" gorm:"type:text"`
	SeverityLevelTitleEn       string    `json:"severity_level_title_en" gorm:"type:varchar(255);not null"`
	SeverityLevelTitleAr       string    `json:"severity_level_title_ar" gorm:"type:varchar(255);not null"`
	ResponsibleUserNameEn      string    `json:"responsible_user_name_en" gorm:"type:varchar(255);not null"`
	ResponsibleUserNameAr      string    `json:"responsible_user_name_ar" gorm:"type:varchar(255);not null"`
	TargetDate                 time.Time `json:"target_date" gorm:"type:date;not null"`
	ActionPlanStatusTitleEn    string    `json:"action_plan_status_title_en" gorm:"type:varchar(255);not null"`
	ActionPlanStatusTitleAr    string    `json:"action_plan_status_title_ar" gorm:"type:varchar(255);not null"`
	ActionPlanStatusColorCode  string    `json:"action_plan_status_color_code" gorm:"type:varchar(7);not null"`

	// Relationships
	RiskReportIdentification *RiskReportIdentification `gorm:"foreignKey:RiskReportIdentificationID"`
	RiskReport               *RiskReport               `gorm:"foreignKey:RiskReportID"`
}

func (RiskReportActionPlan) TableName() string {
	return "risk_report_action_plans"
}
