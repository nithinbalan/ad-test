package models

import "time"

type IncidentReportActionPlan struct {
	ID                       uint `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	IncidentReportID         uint `json:"incident_report_id" gorm:"not null"`
	IncidentReportIncidentID uint `json:"incident_report_incident_id" gorm:"not null"`

	ActionPlanTitle       string     `json:"action_plan_title" gorm:"type:nvarchar(255);not null"`
	ActionPlanDescription *string    `json:"action_plan_description" gorm:"type:nvarchar(max)"`
	TargetDate            *time.Time `json:"target_date" gorm:"type:date;null"`
	ResponsibleUserNameEn string     `json:"responsible_user_name_en" gorm:"type:nvarchar(255);not null"`
	ResponsibleUserNameAr string     `json:"responsible_user_name_ar" gorm:"type:nvarchar(255);not null"`

	SeverityLevelTitleEn *string `json:"severity_level_title_en" gorm:"type:nvarchar(255)"`
	SeverityLevelTitleAr *string `json:"severity_level_title_ar" gorm:"type:nvarchar(255)"`

	ActionPlanStatusTitleEn   string `json:"action_plan_status_title_en" gorm:"type:nvarchar(255);not null"`
	ActionPlanStatusTitleAr   string `json:"action_plan_status_title_ar" gorm:"type:nvarchar(255);not null"`
	ActionPlanStatusColorCode string `json:"action_plan_status_color_code" gorm:"type:nvarchar(7);not null"`

	// Relations
	IncidentReport         *IncidentReport         `gorm:"foreignKey:IncidentReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	IncidentReportIncident *IncidentReportIncident `gorm:"foreignKey:IncidentReportIncidentID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (IncidentReportActionPlan) TableName() string {
	return "incident_report_action_plans"
}
