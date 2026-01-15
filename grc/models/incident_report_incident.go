package models

import "time"

type IncidentReportIncident struct {
	ID                         uint       `json:"id" gorm:"primary_key;auto_increment"`
	IncidentReportID           uint       `json:"incident_report_id" gorm:"not null"`
	IncidentRegisterID         uint       `json:"incident_register_id" gorm:"not null"`
	ReferenceNumber            *string    `json:"reference_number" gorm:"type:nvarchar(50)"`
	IncidentTitle              string     `json:"incident_title" gorm:"type:nvarchar(255);not null"`
	IncidentCategoryTitleEn    *string    `json:"incident_category_title_en" gorm:"type:nvarchar(255)"`
	IncidentCategoryTitleAr    *string    `json:"incident_category_title_ar" gorm:"type:nvarchar(255)"`
	IncidentSubCategoryTitleEn *string    `json:"incident_sub_category_title_en" gorm:"type:nvarchar(255)"`
	IncidentSubCategoryTitleAr *string    `json:"incident_sub_category_title_ar" gorm:"type:nvarchar(255)"`
	SeverityLevelTitleEn       *string    `json:"severity_level_title_en" gorm:"type:nvarchar(255)"`
	SeverityLevelTitleAr       *string    `json:"severity_level_title_ar" gorm:"type:nvarchar(255)"`
	SeverityLevelColorCode     *string    `json:"severity_level_color_code" gorm:"type:nvarchar(255)"`
	IncidentAt                 *time.Time `json:"incident_at" gorm:"type:datetime"`
	ReportedAt                 *time.Time `json:"reported_at" gorm:"type:datetime"`
	IncidentLocation           string     `json:"incident_location" gorm:"type:nvarchar(255);not null"`
	IncidentStatusTitleEn      *string    `json:"incident_status_title_en" gorm:"type:nvarchar(255)"`
	IncidentStatusTitleAr      *string    `json:"incident_status_title_ar" gorm:"type:nvarchar(255)"`
	ResponsibleUserNameEn      string     `json:"responsible_user_name_en" gorm:"type:nvarchar(255);not null"`
	ResponsibleUserNameAr      string     `json:"responsible_user_name_ar" gorm:"type:nvarchar(255);not null"`
	ReportedUserNameEn         string     `json:"reported_user_name_en" gorm:"type:nvarchar(255);not null"`
	ReportedUserNameAr         string     `json:"reported_user_name_ar" gorm:"type:nvarchar(255);not null"`
	CreatedAt                  *time.Time `json:"created_at" gorm:"type:datetime"`
	UpdatedAt                  *time.Time `json:"updated_at" gorm:"type:datetime"`
	ClosedDate                 *time.Time `json:"closed_date" gorm:"type:datetime"`

	IncidentReport    *IncidentReport                    `gorm:"foreignKey:IncidentReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	IncidentRegister  *IncidentRegister                  `gorm:"foreignKey:IncidentRegisterID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Documents         []*IncidentReportDocument          `gorm:"foreignKey:IncidentReportIncidentID"`
	RootCauseAnalyses []*IncidentReportRootCauseAnalysis `gorm:"foreignKey:IncidentReportIncidentID"`
	ActionPlans       []*IncidentReportActionPlan        `gorm:"foreignKey:IncidentReportIncidentID"`
}

func (IncidentReportIncident) TableName() string {
	return "incident_report_incidents"
}
