package models

import "time"

type IncidentReportIncident struct {
	ID                         uint       `json:"id" gorm:"primary_key;auto_increment"`
	IncidentReportID           uint       `json:"incident_report_id" gorm:"not null"`
	IncidentRegisterID         uint       `json:"incident_register_id" gorm:"not null"`
	ReferenceNumber            *string    `json:"reference_number"`
	IncidentCategoryTitleEn    string     `json:"incident_category_title_en"`
	IncidentCategoryTitleAr    string     `json:"incident_category_title_ar"`
	IncidentSubCategoryTitleEn *string    `json:"incident_sub_category_title_en"`
	IncidentSubCategoryTitleAr *string    `json:"incident_sub_category_title_ar"`
	SeverityTitleEn            string     `json:"severity_title_en"`
	SeverityTitleAr            string     `json:"severity_title_ar"`
	IncidentAt                 time.Time  `json:"incident_at"`
	ReportedAt                 time.Time  `json:"reported_at"`
	IndicatorsOfCompromise     *string    `json:"indicators_of_compromise" gorm:"type:nvarchar(max)"`
	IncidentStatusTitleEn      string     `json:"incident_status_title_en"`
	IncidentStatusTitleAr      string     `json:"incident_status_title_ar"`
	ContainmentStrategy        *string    `json:"containment_strategy" gorm:"type:nvarchar(max)"`
	NextSteps                  *string    `json:"next_steps" gorm:"type:nvarchar(max)"`
	AssignedUserInitial        *string    `json:"assigned_user_initial"`
	AssignedUserFirstNameEn    *string    `json:"assigned_user_first_name_en"`
	AssignedUserLastNameEn     *string    `json:"assigned_user_last_name_en"`
	AssignedUserFirstNameAr    *string    `json:"assigned_user_first_name_ar"`
	AssignedUserLastNameAr     *string    `json:"assigned_user_last_name_ar"`
	ResolvedDate               *time.Time `json:"resolved_date"`
	RecoveryMeasures           *string    `json:"recovery_measures" gorm:"type:nvarchar(max)"`

	IncidentReport *IncidentReport `gorm:"foreignKey:IncidentReportID"`
}

func (IncidentReportIncident) TableName() string {
	return "incident_report_incidents"
}
