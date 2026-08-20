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

	// Incident Register - additional descriptive and related fields
	Description                              *string    `json:"description" gorm:"type:nvarchar(max)"`
	QuickAction                              *string    `json:"quick_action" gorm:"type:nvarchar(max)"`
	IncidentTypeTitleEn                      *string    `json:"incident_type_title_en" gorm:"type:nvarchar(255)"`
	IncidentTypeTitleAr                      *string    `json:"incident_type_title_ar" gorm:"type:nvarchar(255)"`
	IncidentStatusColorCode                  *string    `json:"incident_status_color_code" gorm:"type:nvarchar(255)"`
	ResponsibleUnitTitleEn                   *string    `json:"responsible_unit_title_en" gorm:"type:nvarchar(255)"`
	ResponsibleUnitTitleAr                   *string    `json:"responsible_unit_title_ar" gorm:"type:nvarchar(255)"`
	ResponsibleUnitColorCode                 *string    `json:"responsible_unit_color_code" gorm:"type:nvarchar(255)"`
	BusinessOwnerNameEn                      *string    `json:"business_owner_name_en" gorm:"type:nvarchar(255)"`
	BusinessOwnerNameAr                      *string    `json:"business_owner_name_ar" gorm:"type:nvarchar(255)"`
	ResponseTeamLeadNameEn                   *string    `json:"response_team_lead_name_en" gorm:"type:nvarchar(255)"`
	ResponseTeamLeadNameAr                   *string    `json:"response_team_lead_name_ar" gorm:"type:nvarchar(255)"`
	ResponseTeamMembersEn                    *string    `json:"response_team_members_en" gorm:"type:nvarchar(max)"`
	ResponseTeamMembersAr                    *string    `json:"response_team_members_ar" gorm:"type:nvarchar(max)"`
	IncidentScopeTitleEn                     *string    `json:"incident_scope_title_en" gorm:"type:nvarchar(255)"`
	IncidentScopeTitleAr                     *string    `json:"incident_scope_title_ar" gorm:"type:nvarchar(255)"`
	IncidentSensitivityClassificationTitleEn *string    `json:"incident_sensitivity_classification_title_en" gorm:"type:nvarchar(255)"`
	IncidentSensitivityClassificationTitleAr *string    `json:"incident_sensitivity_classification_title_ar" gorm:"type:nvarchar(255)"`
	ConfidentialitySeverityLevelTitleEn      *string    `json:"confidentiality_severity_level_title_en" gorm:"type:nvarchar(255)"`
	ConfidentialitySeverityLevelTitleAr      *string    `json:"confidentiality_severity_level_title_ar" gorm:"type:nvarchar(255)"`
	ConfidentialitySeverityLevelColorCode    *string    `json:"confidentiality_severity_level_color_code" gorm:"type:nvarchar(255)"`
	IntegritySeverityLevelTitleEn            *string    `json:"integrity_severity_level_title_en" gorm:"type:nvarchar(255)"`
	IntegritySeverityLevelTitleAr            *string    `json:"integrity_severity_level_title_ar" gorm:"type:nvarchar(255)"`
	IntegritySeverityLevelColorCode          *string    `json:"integrity_severity_level_color_code" gorm:"type:nvarchar(255)"`
	AvailabilitySeverityLevelTitleEn         *string    `json:"availability_severity_level_title_en" gorm:"type:nvarchar(255)"`
	AvailabilitySeverityLevelTitleAr         *string    `json:"availability_severity_level_title_ar" gorm:"type:nvarchar(255)"`
	AvailabilitySeverityLevelColorCode       *string    `json:"availability_severity_level_color_code" gorm:"type:nvarchar(255)"`
	IsAcknowledge                            bool       `json:"is_acknowledge"`
	AcknowledgedAt                           *time.Time `json:"acknowledged_at" gorm:"type:datetime"`
	AcknowledgedByNameEn                     *string    `json:"acknowledged_by_name_en" gorm:"type:nvarchar(255)"`
	AcknowledgedByNameAr                     *string    `json:"acknowledged_by_name_ar" gorm:"type:nvarchar(255)"`
	IsSLATimeAchieved                        bool       `json:"is_sla_time_achieved"`
	IsSLATimeExceeded                        bool       `json:"is_sla_time_exceeded"`

	// Incident Closure details
	ClosureResolutionSummary       *string    `json:"closure_resolution_summary" gorm:"type:nvarchar(max)"`
	ClosureNotes                   *string    `json:"closure_notes" gorm:"type:nvarchar(max)"`
	ClosureSubmittedAt             *time.Time `json:"closure_submitted_at" gorm:"type:datetime"`
	ClosureSubmittedByNameEn       *string    `json:"closure_submitted_by_name_en" gorm:"type:nvarchar(255)"`
	ClosureSubmittedByNameAr       *string    `json:"closure_submitted_by_name_ar" gorm:"type:nvarchar(255)"`
	ClosureApprovedAt              *time.Time `json:"closure_approved_at" gorm:"type:datetime"`
	ClosureApprovedByNameEn        *string    `json:"closure_approved_by_name_en" gorm:"type:nvarchar(255)"`
	ClosureApprovedByNameAr        *string    `json:"closure_approved_by_name_ar" gorm:"type:nvarchar(255)"`
	ClosureWorkflowStatusTitleEn   *string    `json:"closure_workflow_status_title_en" gorm:"type:nvarchar(255)"`
	ClosureWorkflowStatusTitleAr   *string    `json:"closure_workflow_status_title_ar" gorm:"type:nvarchar(255)"`
	ClosureWorkflowStatusColorCode *string    `json:"closure_workflow_status_color_code" gorm:"type:nvarchar(255)"`

	IncidentReport    *IncidentReport                    `gorm:"foreignKey:IncidentReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	IncidentRegister  *IncidentRegister                  `gorm:"foreignKey:IncidentRegisterID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Documents         []*IncidentReportDocument          `gorm:"foreignKey:IncidentReportIncidentID"`
	Evidences         []*IncidentReportEvidence          `gorm:"foreignKey:IncidentReportIncidentID"`
	RootCauseAnalyses []*IncidentReportRootCauseAnalysis `gorm:"foreignKey:IncidentReportIncidentID"`
	ActionPlans       []*IncidentReportActionPlan        `gorm:"foreignKey:IncidentReportIncidentID"`
}

func (IncidentReportIncident) TableName() string {
	return "incident_report_incidents"
}
