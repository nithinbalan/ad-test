package models

import "time"

type ProjectReportProject struct {
	ID                           uint      `json:"id" gorm:"primary_key;auto_increment"`
	ProjectReportID              uint      `json:"project_report_id" gorm:"not null"`
	ProjectID                    uint      `json:"project_id" gorm:"not null"`
	TitleEn                      string    `json:"title_en" gorm:"type:nvarchar(255);not null"`
	TitleAr                      string    `json:"title_ar" gorm:"type:nvarchar(255);not null"`
	DescriptionEn                *string   `json:"description_en" gorm:"type:nvarchar(max)"`
	DescriptionAr                *string   `json:"description_ar" gorm:"type:nvarchar(max)"`
	Budget                       *float64  `json:"budget" gorm:"type:decimal(18,2)"`
	Scope                        *string   `json:"scope" gorm:"type:nvarchar(max)"`
	StartDate                    time.Time `json:"start_date" gorm:"type:date;not null"`
	EndDate                      time.Time `json:"end_date" gorm:"type:date;not null"`
	ProjectTypeTitleEn           *string   `json:"project_type_title_en" gorm:"type:nvarchar(255)"`
	ProjectTypeTitleAr           *string   `json:"project_type_title_ar" gorm:"type:nvarchar(255)"`
	ProjectStatusTitleEn         *string   `json:"project_status_title_en" gorm:"type:nvarchar(255)"`
	ProjectStatusTitleAr         *string   `json:"project_status_title_ar" gorm:"type:nvarchar(255)"`
	ProjectPhaseTitleEn          *string   `json:"project_phase_title_en" gorm:"type:nvarchar(255)"`
	ProjectPhaseTitleAr          *string   `json:"project_phase_title_ar" gorm:"type:nvarchar(255)"`
	ProjectClassificationTitleEn *string   `json:"project_classification_title_en" gorm:"type:nvarchar(255)"`
	ProjectClassificationTitleAr *string   `json:"project_classification_title_ar" gorm:"type:nvarchar(255)"`
	ResponsibleUserNameEn        string    `json:"responsible_user_name_en" gorm:"type:nvarchar(255);not null"`
	ResponsibleUserNameAr        string    `json:"responsible_user_name_ar" gorm:"type:nvarchar(255);not null"`
	IncidentActionPlanID         *uint     `json:"incident_action_plan_id" gorm:"null"`
	IncidentActionPlanTitle      *string   `json:"incident_action_plan_title" gorm:"type:nvarchar(255);null"`
	CreatedAt                    time.Time `json:"created_at" gorm:"type:datetime;not null"`

	ProjectReport *ProjectReport           `gorm:"foreignKey:ProjectReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Project       *Project                 `gorm:"foreignKey:ProjectID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Documents     []*ProjectReportDocument `gorm:"foreignKey:ProjectReportProjectID"`
}

func (ProjectReportProject) TableName() string {
	return "project_report_projects"
}
