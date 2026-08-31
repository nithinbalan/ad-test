package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"time"
)

type IncidentReportDetailDTO struct {
	ID                                uint                                              `json:"id"`
	ReferenceNumber                   string                                            `json:"reference_number"`
	Title                             string                                            `json:"title"`
	Description                       string                                            `json:"description"`
	IncidentAtFormatted               string                                            `json:"incident_at_formatted"`
	ReportedAtFormatted               string                                            `json:"reported_at_formatted"`
	Location                          *string                                           `json:"location"`
	QuickAction                       *string                                           `json:"quick_action"`
	IncidentCategory                  IncidentCategoryDTO                               `json:"incident_category"`
	IncidentSubCategory               *IncidentSubCategoryDTO                           `json:"incident_sub_category"`
	IncidentType                      IncidentTypeDTO                                   `json:"incident_type"`
	SeverityLevel                     *masterDTOs.SeverityLevelDTO                      `json:"severity_level"`
	IncidentStatus                    IncidentStatusDTO                                 `json:"incident_status"`
	ReportedByUser                    organizationDTOs.UserMinimalDTO                   `json:"reported_by"`
	ResponsibleUnit                   settingsDTOs.OrganizationUnitDTO                  `json:"responsible_unit"`
	ResponsibleUser                   organizationDTOs.UserMinimalDTO                   `json:"responsible_user"`
	ResponsibleUserID                 uint                                              `json:"responsible_user_id" binding:"required" gorm:"not null"`
	BusinessOwner                     *organizationDTOs.UserMinimalDTO                  `json:"business_owner"`
	IncidentResponseTeam              *IncidentResponseTeamDetailDTO                    `json:"incident_response_team"`
	ResponseTeamLead                  *organizationDTOs.UserMinimalDTO                  `json:"response_team_lead"`
	ResponseTeamMembers               []*IncidentRegisterResponseTeamMemberDTO          `json:"response_team_members"`
	IncidentScope                     *IncidentScopeDTO                                 `json:"incident_scope"`
	IncidentSensitivityClassification *IncidentSensitivityClassificationDTO             `json:"incident_sensitivity_classification"`
	ConfidentialitySeverityLevel      *masterDTOs.SeverityLevelDTO                      `json:"confidentiality_severity_level"`
	IntegritySeverityLevel            *masterDTOs.SeverityLevelDTO                      `json:"integrity_severity_level"`
	AvailabilitySeverityLevel         *masterDTOs.SeverityLevelDTO                      `json:"availability_severity_level"`
	IsAcknowledge                     bool                                              `json:"is_acknowledge"`
	AcknowledgedAt                    *time.Time                                        `json:"acknowledged_at"`
	AcknowledgedByUser                *organizationDTOs.UserMinimalDTO                  `json:"acknowledged_by"`
	IsSLATimeAchieved                 bool                                              `json:"is_sla_time_achieved"`
	IsSLATimeExceeded                 bool                                              `json:"is_sla_time_exceeded"`
	Documents                         []*IncidentDocumentDTO                            `json:"documents"`
	Controls                          []*commonDTOs.AssociatedControlListDTO            `json:"controls"`
	Risks                             []*commonDTOs.AssociatedRiskListDTO               `json:"risks"`
	Assets                            []*commonDTOs.AssociatedAssetListDTO              `json:"assets"`
	Processes                         []*commonDTOs.AssociatedProcessListDTO            `json:"processes"`
	Timelines                         []*IncidentTimelineListDTO                        `json:"timelines"`
	SLARegisters                      []*IncidentSLARegisterDTO                         `json:"sla_registers"`
	AssociatedDocuments               []*commonDTOs.AssociatedDocumentListDTO           `json:"associated_documents"`
	CreatedAt                         time.Time                                         `json:"created_at"`
	UpdatedAt                         *time.Time                                        `json:"updated_at"`
	CreatedByUser                     organizationDTOs.UserMinimalDTO                   `json:"created_by"`
	UpdatedByUser                     organizationDTOs.UserMinimalDTO                   `json:"updated_by"`
	GeneratedAt                       string                                            `json:"generated_at"`
	CommunicationLogs                 []*IncidentCommunicationLogDetailDTO              `json:"communication_log"`
	ActionPlans                       []*IncidentActionPlanListDTO                      `json:"action_plan"`
	IncidentClosure                   *IncidentClosureDetailDTO                         `json:"incident_closure"`
	IncidentPostReview                *IncidentPostReviewDetailDTO                      `json:"incident_post_review"`
	RootCauseAnalyses                 []*IncidentRootCauseAnalysisListDTO               `json:"root_cause_analysis"`
	StrategicObjectives               []*commonDTOs.AssociatedStrategicObjectiveListDTO `json:"strategic_objectives"`
	KPIs                              []*commonDTOs.AssociatedKPIRegisterListDTO        `json:"kpi"`
	ComplianceRequirements            []*commonDTOs.AssociatedClauseListDTO             `json:"compliance_clause"`
	Investigation                     *IncidentInvestigationDetailDTO                   `json:"investigation"`
}

type IncidentReportRequestDTO struct {
	Type        string  `json:"type" binding:"required,oneof=individual date-range"`
	IncidentIDs *[]uint `json:"incident_ids"`
	Title       string  `json:"title" binding:"required"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date"`
}

type IncidentReportParsedDates struct {
	StartDate *time.Time
	EndDate   *time.Time
}

type IncidentReportIncidentDTO struct {
	ID              uint   `json:"id"`
	ReferenceNumber string `json:"reference_number"`
	IncidentTitle   string `json:"incident_title"`

	IncidentCategoryTitleEn string `json:"incident_category_title_en"`
	IncidentCategoryTitleAr string `json:"incident_category_title_ar"`

	IncidentSubCategoryTitleEn string `json:"incident_sub_category_title_en"`
	IncidentSubCategoryTitleAr string `json:"incident_sub_category_title_ar"`

	IncidentTypeTitleEn string `json:"incident_type_title_en"`
	IncidentTypeTitleAr string `json:"incident_type_title_ar"`

	SeverityTitleEn   string  `json:"severity_title_en"`
	SeverityTitleAr   string  `json:"severity_title_ar"`
	SeverityColorCode *string `json:"severity_color_code"`

	RealTime         time.Time `json:"real_time"`
	DetectedTime     time.Time `json:"detected_time"`
	IncidentLocation string    `json:"incident_location"`

	IncidentStatusTitleEn   string `json:"incident_status_title_en"`
	IncidentStatusTitleAr   string `json:"incident_status_title_ar"`
	IncidentStatusColorCode string `json:"incident_status_color_code"`

	ResponsibleUserNameEn string `json:"responsible_user_name_en"`
	ResponsibleUserNameAr string `json:"responsible_user_name_ar"`
	ReportedUserNameEn    string `json:"reported_user_name_en"`
	ReportedUserNameAr    string `json:"reported_user_name_ar"`

	ResponsibleUnitTitleEn   string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr   string `json:"responsible_unit_title_ar"`
	ResponsibleUnitColorCode string `json:"responsible_unit_color_code"`

	BusinessOwnerNameEn string `json:"business_owner_name_en"`
	BusinessOwnerNameAr string `json:"business_owner_name_ar"`

	ResponseTeamLeadNameEn string `json:"response_team_lead_name_en"`
	ResponseTeamLeadNameAr string `json:"response_team_lead_name_ar"`

	ResponseTeamMembersEn string `json:"response_team_members_en"`
	ResponseTeamMembersAr string `json:"response_team_members_ar"`

	IncidentScopeTitleEn string `json:"incident_scope_title_en"`
	IncidentScopeTitleAr string `json:"incident_scope_title_ar"`

	IncidentSensitivityClassificationTitleEn string `json:"incident_sensitivity_classification_title_en"`
	IncidentSensitivityClassificationTitleAr string `json:"incident_sensitivity_classification_title_ar"`

	ConfidentialitySeverityLevelTitleEn   string `json:"confidentiality_severity_level_title_en"`
	ConfidentialitySeverityLevelTitleAr   string `json:"confidentiality_severity_level_title_ar"`
	ConfidentialitySeverityLevelColorCode string `json:"confidentiality_severity_level_color_code"`

	IntegritySeverityLevelTitleEn   string `json:"integrity_severity_level_title_en"`
	IntegritySeverityLevelTitleAr   string `json:"integrity_severity_level_title_ar"`
	IntegritySeverityLevelColorCode string `json:"integrity_severity_level_color_code"`

	AvailabilitySeverityLevelTitleEn   string `json:"availability_severity_level_title_en"`
	AvailabilitySeverityLevelTitleAr   string `json:"availability_severity_level_title_ar"`
	AvailabilitySeverityLevelColorCode string `json:"availability_severity_level_color_code"`

	Description string `json:"description"`
	QuickAction string `json:"quick_action"`

	IsAcknowledge        bool       `json:"is_acknowledge"`
	AcknowledgedAt       *time.Time `json:"acknowledged_at"`
	AcknowledgedByNameEn string     `json:"acknowledged_by_name_en"`
	AcknowledgedByNameAr string     `json:"acknowledged_by_name_ar"`

	IsSLATimeAchieved bool `json:"is_sla_time_achieved"`
	IsSLATimeExceeded bool `json:"is_sla_time_exceeded"`

	ContainmentStrategy *string `json:"containment_strategy"`
	NextSteps           *string `json:"next_steps"`

	ResolvedDate     *time.Time `json:"resolved_date"`
	RecoveryMeasures *string    `json:"recovery_measures"`

	ClosureResolutionSummary       string     `json:"closure_resolution_summary"`
	ClosureNotes                   string     `json:"closure_notes"`
	ClosureSubmittedAt             *time.Time `json:"closure_submitted_at"`
	ClosureSubmittedByNameEn       string     `json:"closure_submitted_by_name_en"`
	ClosureSubmittedByNameAr       string     `json:"closure_submitted_by_name_ar"`
	ClosureApprovedAt              *time.Time `json:"closure_approved_at"`
	ClosureApprovedByNameEn        string     `json:"closure_approved_by_name_en"`
	ClosureApprovedByNameAr        string     `json:"closure_approved_by_name_ar"`
	ClosureWorkflowStatusTitleEn   string     `json:"closure_workflow_status_title_en"`
	ClosureWorkflowStatusTitleAr   string     `json:"closure_workflow_status_title_ar"`
	ClosureWorkflowStatusColorCode string     `json:"closure_workflow_status_color_code"`

	AssignedUserInitial     *string `json:"assigned_user_initial"`
	AssignedUserFirstNameEn *string `json:"assigned_user_first_name_en"`
	AssignedUserLastNameEn  *string `json:"assigned_user_last_name_en"`
	AssignedUserFirstNameAr *string `json:"assigned_user_first_name_ar"`
	AssignedUserLastNameAr  *string `json:"assigned_user_last_name_ar"`
}

type IncidentReportDocumentDTO struct {
	ID                       uint    `json:"id" mapstructure:"ID"`
	IncidentDocumentID       uint    `json:"incident_document_id" mapstructure:"IncidentDocumentID"`
	IncidentReportIncidentID uint    `json:"incident_report_incident_id" mapstructure:"IncidentReportIncidentID"`
	URL                      *string `json:"url" mapstructure:"URL"`
	Name                     *string `json:"name" mapstructure:"Name"`
	Size                     *int    `json:"size" mapstructure:"Size"`
	Token                    *string `json:"token" mapstructure:"Token"`
	OriginalName             *string `json:"original_name" mapstructure:"OriginalName"`
	Type                     *string `json:"type" mapstructure:"Type"`
}

type IncidentReportEvidenceDTO struct {
	ID                       uint    `json:"id" mapstructure:"ID"`
	IncidentEvidenceID       uint    `json:"incident_evidence_id" mapstructure:"IncidentEvidenceID"`
	IncidentReportIncidentID uint    `json:"incident_report_incident_id" mapstructure:"IncidentReportIncidentID"`
	URL                      *string `json:"url" mapstructure:"URL"`
	Name                     *string `json:"name" mapstructure:"Name"`
	Size                     *int    `json:"size" mapstructure:"Size"`
	Token                    *string `json:"token" mapstructure:"Token"`
	OriginalName             *string `json:"original_name" mapstructure:"OriginalName"`
	Type                     *string `json:"type" mapstructure:"Type"`
}

type IncidentReportListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
	Filters  map[string][]string
}

type IncidentReportListDTO struct {
	ID                          uint   `json:"id"`
	Title                       string `json:"title"`
	StartDate                   string `json:"start_date"`
	EndDate                     string `json:"end_date"`
	CreatedAt                   string `json:"created_at"`
	CreatedUserID               uint   `json:"created_user_id"`
	CreatedUserFirstNameEn      string `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string `json:"created_user_last_name_ar"`
	CreatedUserInitial          string `json:"created_user_initial"`
	CreatedUserColor            string `json:"created_user_color"`
	CreatedUserEmail            string `json:"created_user_email"`
	CreatedUserProfileImagePath string `json:"created_user_profile_image_path"`
}

type IncidentReportListResponseDTO struct {
	Items      []*IncidentReportListDTO     `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type UserInternalDTO struct {
	ID               uint    `json:"id" mapstructure:"ID"`
	FirstNameEn      string  `json:"first_name_en" mapstructure:"FirstNameEn"`
	LastNameEn       string  `json:"last_name_en" mapstructure:"LastNameEn"`
	FirstNameAr      string  `json:"first_name_ar" mapstructure:"FirstNameAr"`
	LastNameAr       string  `json:"last_name_ar" mapstructure:"LastNameAr"`
	ProfileImagePath *string `json:"profile_image_path" mapstructure:"ProfileImagePath"`
	Email            string  `json:"email" mapstructure:"Email"`
	Initial          string  `json:"initial" mapstructure:"Initial"`
	Color            string  `json:"color" mapstructure:"Color"`
	MobileNumber     *string `json:"mobile_number" mapstructure:"MobileNumber"`
}

type IncidentReportDetailsDTO struct {
	ID            uint                               `json:"id" mapstructure:"ID"`
	Title         string                             `json:"title" mapstructure:"Title"`
	StartDate     *time.Time                         `json:"start_date" mapstructure:"StartDate"`
	EndDate       *time.Time                         `json:"end_date" mapstructure:"EndDate"`
	CreatedAt     time.Time                          `json:"created_at" mapstructure:"CreatedAt"`
	CreatedByUser *UserInternalDTO                   `json:"created_by" mapstructure:"CreatedByUser"`
	Sections      []*IncidentReportSectionDTO        `json:"sections" mapstructure:"Sections"`
	Incidents     []*IncidentReportIncidentDetailDTO `json:"incidents" mapstructure:"Incidents"`
}

type IncidentReportSectionDTO struct {
	ID            uint   `json:"id" mapstructure:"ID"`
	Type          string `json:"type" mapstructure:"Type"`
	TitleEn       string `json:"title_en" mapstructure:"TitleEn"`
	TitleAr       string `json:"title_ar" mapstructure:"TitleAr"`
	DescriptionEn string `json:"description_en" mapstructure:"DescriptionEn"`
	DescriptionAr string `json:"description_ar" mapstructure:"DescriptionAr"`
	OrderNum      int    `json:"order_num" mapstructure:"OrderNum"`
}

type IncidentReportIncidentDetailDTO struct {
	ID                                       uint                                  `json:"id" mapstructure:"ID"`
	IncidentRegisterID                       uint                                  `json:"incident_register_id" mapstructure:"IncidentRegisterID"`
	ReferenceNumber                          *string                               `json:"reference_number" mapstructure:"ReferenceNumber"`
	IncidentTitle                            string                                `json:"incident_title" mapstructure:"IncidentTitle"`
	IncidentCategoryTitleEn                  *string                               `json:"incident_category_title_en" mapstructure:"IncidentCategoryTitleEn"`
	IncidentCategoryTitleAr                  *string                               `json:"incident_category_title_ar" mapstructure:"IncidentCategoryTitleAr"`
	IncidentSubCategoryTitleEn               *string                               `json:"incident_sub_category_title_en" mapstructure:"IncidentSubCategoryTitleEn"`
	IncidentSubCategoryTitleAr               *string                               `json:"incident_sub_category_title_ar" mapstructure:"IncidentSubCategoryTitleAr"`
	IncidentTypeTitleEn                      *string                               `json:"incident_type_title_en" mapstructure:"IncidentTypeTitleEn"`
	IncidentTypeTitleAr                      *string                               `json:"incident_type_title_ar" mapstructure:"IncidentTypeTitleAr"`
	SeverityLevelTitleEn                     *string                               `json:"severity_level_title_en" mapstructure:"SeverityLevelTitleEn"`
	SeverityLevelTitleAr                     *string                               `json:"severity_level_title_ar" mapstructure:"SeverityLevelTitleAr"`
	SeverityLevelColorCode                   *string                               `json:"severity_level_color_code" mapstructure:"SeverityLevelColorCode"`
	IncidentAt                               *time.Time                            `json:"incident_at" mapstructure:"IncidentAt"`
	ReportedAt                               *time.Time                            `json:"reported_at" mapstructure:"ReportedAt"`
	IncidentLocation                         string                                `json:"incident_location" mapstructure:"IncidentLocation"`
	IncidentStatusTitleEn                    *string                               `json:"incident_status_title_en" mapstructure:"IncidentStatusTitleEn"`
	IncidentStatusTitleAr                    *string                               `json:"incident_status_title_ar" mapstructure:"IncidentStatusTitleAr"`
	IncidentStatusColorCode                  *string                               `json:"incident_status_color_code" mapstructure:"IncidentStatusColorCode"`
	ResponsibleUserNameEn                    string                                `json:"responsible_user_name_en" mapstructure:"ResponsibleUserNameEn"`
	ResponsibleUserNameAr                    string                                `json:"responsible_user_name_ar" mapstructure:"ResponsibleUserNameAr"`
	ReportedUserNameEn                       string                                `json:"reported_user_name_en" mapstructure:"ReportedUserNameEn"`
	ReportedUserNameAr                       string                                `json:"reported_user_name_ar" mapstructure:"ReportedUserNameAr"`
	ResponsibleUnitTitleEn                   *string                               `json:"responsible_unit_title_en" mapstructure:"ResponsibleUnitTitleEn"`
	ResponsibleUnitTitleAr                   *string                               `json:"responsible_unit_title_ar" mapstructure:"ResponsibleUnitTitleAr"`
	ResponsibleUnitColorCode                 *string                               `json:"responsible_unit_color_code" mapstructure:"ResponsibleUnitColorCode"`
	BusinessOwnerNameEn                      *string                               `json:"business_owner_name_en" mapstructure:"BusinessOwnerNameEn"`
	BusinessOwnerNameAr                      *string                               `json:"business_owner_name_ar" mapstructure:"BusinessOwnerNameAr"`
	ResponseTeamLeadNameEn                   *string                               `json:"response_team_lead_name_en" mapstructure:"ResponseTeamLeadNameEn"`
	ResponseTeamLeadNameAr                   *string                               `json:"response_team_lead_name_ar" mapstructure:"ResponseTeamLeadNameAr"`
	ResponseTeamMembersEn                    *string                               `json:"response_team_members_en" mapstructure:"ResponseTeamMembersEn"`
	ResponseTeamMembersAr                    *string                               `json:"response_team_members_ar" mapstructure:"ResponseTeamMembersAr"`
	IncidentScopeTitleEn                     *string                               `json:"incident_scope_title_en" mapstructure:"IncidentScopeTitleEn"`
	IncidentScopeTitleAr                     *string                               `json:"incident_scope_title_ar" mapstructure:"IncidentScopeTitleAr"`
	IncidentSensitivityClassificationTitleEn *string                               `json:"incident_sensitivity_classification_title_en" mapstructure:"IncidentSensitivityClassificationTitleEn"`
	IncidentSensitivityClassificationTitleAr *string                               `json:"incident_sensitivity_classification_title_ar" mapstructure:"IncidentSensitivityClassificationTitleAr"`
	ConfidentialitySeverityLevelTitleEn      *string                               `json:"confidentiality_severity_level_title_en" mapstructure:"ConfidentialitySeverityLevelTitleEn"`
	ConfidentialitySeverityLevelTitleAr      *string                               `json:"confidentiality_severity_level_title_ar" mapstructure:"ConfidentialitySeverityLevelTitleAr"`
	ConfidentialitySeverityLevelColorCode    *string                               `json:"confidentiality_severity_level_color_code" mapstructure:"ConfidentialitySeverityLevelColorCode"`
	IntegritySeverityLevelTitleEn            *string                               `json:"integrity_severity_level_title_en" mapstructure:"IntegritySeverityLevelTitleEn"`
	IntegritySeverityLevelTitleAr            *string                               `json:"integrity_severity_level_title_ar" mapstructure:"IntegritySeverityLevelTitleAr"`
	IntegritySeverityLevelColorCode          *string                               `json:"integrity_severity_level_color_code" mapstructure:"IntegritySeverityLevelColorCode"`
	AvailabilitySeverityLevelTitleEn         *string                               `json:"availability_severity_level_title_en" mapstructure:"AvailabilitySeverityLevelTitleEn"`
	AvailabilitySeverityLevelTitleAr         *string                               `json:"availability_severity_level_title_ar" mapstructure:"AvailabilitySeverityLevelTitleAr"`
	AvailabilitySeverityLevelColorCode       *string                               `json:"availability_severity_level_color_code" mapstructure:"AvailabilitySeverityLevelColorCode"`
	Description                              *string                               `json:"description" mapstructure:"Description"`
	QuickAction                              *string                               `json:"quick_action" mapstructure:"QuickAction"`
	IsAcknowledge                            bool                                  `json:"is_acknowledge" mapstructure:"IsAcknowledge"`
	AcknowledgedAt                           *time.Time                            `json:"acknowledged_at" mapstructure:"AcknowledgedAt"`
	AcknowledgedByNameEn                     *string                               `json:"acknowledged_by_name_en" mapstructure:"AcknowledgedByNameEn"`
	AcknowledgedByNameAr                     *string                               `json:"acknowledged_by_name_ar" mapstructure:"AcknowledgedByNameAr"`
	IsSLATimeAchieved                        bool                                  `json:"is_sla_time_achieved" mapstructure:"IsSLATimeAchieved"`
	IsSLATimeExceeded                        bool                                  `json:"is_sla_time_exceeded" mapstructure:"IsSLATimeExceeded"`
	ClosedDate                               *time.Time                            `json:"closed_date" mapstructure:"ClosedDate"`
	ClosureResolutionSummary                 *string                               `json:"closure_resolution_summary" mapstructure:"ClosureResolutionSummary"`
	ClosureNotes                             *string                               `json:"closure_notes" mapstructure:"ClosureNotes"`
	ClosureSubmittedAt                       *time.Time                            `json:"closure_submitted_at" mapstructure:"ClosureSubmittedAt"`
	ClosureSubmittedByNameEn                 *string                               `json:"closure_submitted_by_name_en" mapstructure:"ClosureSubmittedByNameEn"`
	ClosureSubmittedByNameAr                 *string                               `json:"closure_submitted_by_name_ar" mapstructure:"ClosureSubmittedByNameAr"`
	ClosureApprovedAt                        *time.Time                            `json:"closure_approved_at" mapstructure:"ClosureApprovedAt"`
	ClosureApprovedByNameEn                  *string                               `json:"closure_approved_by_name_en" mapstructure:"ClosureApprovedByNameEn"`
	ClosureApprovedByNameAr                  *string                               `json:"closure_approved_by_name_ar" mapstructure:"ClosureApprovedByNameAr"`
	ClosureWorkflowStatusTitleEn             *string                               `json:"closure_workflow_status_title_en" mapstructure:"ClosureWorkflowStatusTitleEn"`
	ClosureWorkflowStatusTitleAr             *string                               `json:"closure_workflow_status_title_ar" mapstructure:"ClosureWorkflowStatusTitleAr"`
	ClosureWorkflowStatusColorCode           *string                               `json:"closure_workflow_status_color_code" mapstructure:"ClosureWorkflowStatusColorCode"`
	Documents                                []*IncidentReportDocumentDTO          `json:"documents" mapstructure:"Documents"`
	Evidences                                []*IncidentReportEvidenceDTO          `json:"evidences" mapstructure:"Evidences"`
	RootCauseAnalyses                        []*IncidentReportRootCauseAnalysisDTO `json:"root_cause_analyses" mapstructure:"RootCauseAnalyses"`
	ActionPlans                              []*IncidentReportActionPlanDTO        `json:"action_plans" mapstructure:"ActionPlans"`
	IncidentRegister                         *IncidentRegisterDTO                  `json:"incident_register,omitempty" mapstructure:"IncidentRegister"`
	Investigation                            *IncidentInvestigationDetailDTO       `json:"investigation,omitempty" mapstructure:"-"`
	CommunicationLogs                        []*IncidentCommunicationLogDetailDTO  `json:"communication_logs,omitempty" mapstructure:"-"`
	IncidentClosure                          *IncidentClosureDetailDTO             `json:"incident_closure,omitempty" mapstructure:"-"`
	IncidentPostReview                       *IncidentPostReviewDetailDTO          `json:"incident_post_review,omitempty" mapstructure:"-"`

	// Live associated items and timeline, fetched at retrieval time (not stored on the snapshot)
	Timelines              []*IncidentTimelineListDTO                        `json:"timelines,omitempty" mapstructure:"-"`
	KPIs                   []*commonDTOs.AssociatedKPIRegisterListDTO        `json:"kpis,omitempty" mapstructure:"-"`
	StrategicObjectives    []*commonDTOs.AssociatedStrategicObjectiveListDTO `json:"strategic_objectives,omitempty" mapstructure:"-"`
	ComplianceRequirements []*commonDTOs.AssociatedClauseListDTO             `json:"compliance_requirements,omitempty" mapstructure:"-"`
	Processes              []*commonDTOs.AssociatedProcessListDTO            `json:"processes,omitempty" mapstructure:"-"`
	Risks                  []*commonDTOs.AssociatedRiskListDTO               `json:"risks,omitempty" mapstructure:"-"`
	ControlsAssociated     []*commonDTOs.AssociatedControlListDTO            `json:"controls,omitempty" mapstructure:"-"`
	AssociatedDocuments    []*commonDTOs.AssociatedDocumentListDTO           `json:"associated_documents,omitempty" mapstructure:"-"`
	Assets                 []*commonDTOs.AssociatedAssetListDTO              `json:"assets,omitempty" mapstructure:"-"`
}

type IncidentRegisterDTO struct {
	Investigation      *IncidentInvestigationDetailDTO      `json:"investigation,omitempty" mapstructure:"Investigation"`
	CommunicationLogs  []*IncidentCommunicationLogDetailDTO `json:"communication_logs,omitempty" mapstructure:"CommunicationLogs"`
	IncidentClosure    *IncidentClosureDetailDTO            `json:"incident_closure,omitempty" mapstructure:"IncidentClosure"`
	IncidentPostReview *IncidentPostReviewDetailDTO         `json:"incident_post_review,omitempty" mapstructure:"IncidentPostReview"`
}

type IncidentReportRootCauseAnalysisDTO struct {
	ID       uint    `json:"id" mapstructure:"ID"`
	Why      uint    `json:"why" mapstructure:"Why"`
	Question string  `json:"question" mapstructure:"Question"`
	Answer   *string `json:"answer" mapstructure:"Answer"`
}

type IncidentReportActionPlanDTO struct {
	ID                        uint       `json:"id" mapstructure:"ID"`
	ActionPlanTitle           string     `json:"action_plan_title" mapstructure:"ActionPlanTitle"`
	ActionPlanDescription     *string    `json:"action_plan_description" mapstructure:"ActionPlanDescription"`
	TargetDate                *time.Time `json:"target_date" mapstructure:"TargetDate"`
	ResponsibleUserNameEn     string     `json:"responsible_user_name_en" mapstructure:"ResponsibleUserNameEn"`
	ResponsibleUserNameAr     string     `json:"responsible_user_name_ar" mapstructure:"ResponsibleUserNameAr"`
	SeverityLevelTitleEn      *string    `json:"severity_level_title_en" mapstructure:"SeverityLevelTitleEn"`
	SeverityLevelTitleAr      *string    `json:"severity_level_title_ar" mapstructure:"SeverityLevelTitleAr"`
	ActionPlanStatusTitleEn   string     `json:"action_plan_status_title_en" mapstructure:"ActionPlanStatusTitleEn"`
	ActionPlanStatusTitleAr   string     `json:"action_plan_status_title_ar" mapstructure:"ActionPlanStatusTitleAr"`
	ActionPlanStatusColorCode string     `json:"action_plan_status_color_code" mapstructure:"ActionPlanStatusColorCode"`
}

type IncidentReportIncidentPDFDTO struct {
	ID                         uint    `json:"id"`
	ReferenceNumber            *string `json:"reference_number"`
	Title                      string  `json:"title"`
	IncidentCategoryTitleEn    *string `json:"incident_category_title_en"`
	IncidentCategoryTitleAr    *string `json:"incident_category_title_ar"`
	IncidentSubCategoryTitleEn *string `json:"incident_sub_category_title_en"`
	IncidentSubCategoryTitleAr *string `json:"incident_sub_category_title_ar"`
	IncidentTypeTitleEn        *string `json:"incident_type_title_en"`
	IncidentTypeTitleAr        *string `json:"incident_type_title_ar"`
	SeverityLevelTitleEn       *string `json:"severity_level_title_en"`
	SeverityLevelTitleAr       *string `json:"severity_level_title_ar"`
	SeverityLevelColorCode     *string `json:"severity_level_color_code"`
	IncidentAtFormatted        string  `json:"incident_at_formatted"`
	ReportedAtFormatted        string  `json:"reported_at_formatted"`
	IncidentLocation           string  `json:"incident_location"`
	IncidentStatusTitleEn      *string `json:"incident_status_title_en"`
	IncidentStatusTitleAr      *string `json:"incident_status_title_ar"`
	IncidentStatusColorCode    *string `json:"incident_status_color_code"`
	ResponsibleUserNameEn      string  `json:"responsible_user_name_en"`
	ResponsibleUserNameAr      string  `json:"responsible_user_name_ar"`
	ReportedUserNameEn         string  `json:"reported_user_name_en"`
	ReportedUserNameAr         string  `json:"reported_user_name_ar"`

	ResponsibleUnitTitleEn   *string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr   *string `json:"responsible_unit_title_ar"`
	ResponsibleUnitColorCode *string `json:"responsible_unit_color_code"`

	BusinessOwnerNameEn string `json:"business_owner_name_en"`
	BusinessOwnerNameAr string `json:"business_owner_name_ar"`

	ResponseTeamLeadNameEn string `json:"response_team_lead_name_en"`
	ResponseTeamLeadNameAr string `json:"response_team_lead_name_ar"`
	ResponseTeamMembersEn  string `json:"response_team_members_en"`
	ResponseTeamMembersAr  string `json:"response_team_members_ar"`

	IncidentScopeTitleEn string `json:"incident_scope_title_en"`
	IncidentScopeTitleAr string `json:"incident_scope_title_ar"`

	IncidentSensitivityClassificationTitleEn string `json:"incident_sensitivity_classification_title_en"`
	IncidentSensitivityClassificationTitleAr string `json:"incident_sensitivity_classification_title_ar"`

	ConfidentialitySeverityLevelTitleEn   string `json:"confidentiality_severity_level_title_en"`
	ConfidentialitySeverityLevelTitleAr   string `json:"confidentiality_severity_level_title_ar"`
	ConfidentialitySeverityLevelColorCode string `json:"confidentiality_severity_level_color_code"`

	IntegritySeverityLevelTitleEn   string `json:"integrity_severity_level_title_en"`
	IntegritySeverityLevelTitleAr   string `json:"integrity_severity_level_title_ar"`
	IntegritySeverityLevelColorCode string `json:"integrity_severity_level_color_code"`

	AvailabilitySeverityLevelTitleEn   string `json:"availability_severity_level_title_en"`
	AvailabilitySeverityLevelTitleAr   string `json:"availability_severity_level_title_ar"`
	AvailabilitySeverityLevelColorCode string `json:"availability_severity_level_color_code"`

	Description string `json:"description"`
	QuickAction string `json:"quick_action"`

	IsAcknowledge        bool       `json:"is_acknowledge"`
	AcknowledgedAt       *time.Time `json:"acknowledged_at"`
	AcknowledgedAtString string     `json:"acknowledged_at_string"`
	AcknowledgedByNameEn string     `json:"acknowledged_by_name_en"`
	AcknowledgedByNameAr string     `json:"acknowledged_by_name_ar"`

	IsSLATimeAchieved bool `json:"is_sla_time_achieved"`
	IsSLATimeExceeded bool `json:"is_sla_time_exceeded"`

	ClosedDateString string `json:"closed_date_string"`

	ClosureResolutionSummary string `json:"closure_resolution_summary"`
	ClosureNotes             string `json:"closure_notes"`

	ClosureSubmittedAtString string `json:"closure_submitted_at_string"`
	ClosureSubmittedByNameEn string `json:"closure_submitted_by_name_en"`
	ClosureSubmittedByNameAr string `json:"closure_submitted_by_name_ar"`

	ClosureApprovedAtString string `json:"closure_approved_at_string"`
	ClosureApprovedByNameEn string `json:"closure_approved_by_name_en"`
	ClosureApprovedByNameAr string `json:"closure_approved_by_name_ar"`

	ClosureWorkflowStatusTitleEn   string `json:"closure_workflow_status_title_en"`
	ClosureWorkflowStatusTitleAr   string `json:"closure_workflow_status_title_ar"`
	ClosureWorkflowStatusColorCode string `json:"closure_workflow_status_color_code"`

	Documents         []*IncidentReportDocumentDTO          `json:"documents"`
	Evidences         []*IncidentReportEvidenceDTO          `json:"evidences"`
	RootCauseAnalyses []*IncidentReportRootCauseAnalysisDTO `json:"root_cause_analyses"`
	ActionPlans       []*IncidentReportActionPlanDTO        `json:"action_plans"`

	Investigation      *IncidentInvestigationDetailDTO      `json:"investigation,omitempty"`
	CommunicationLogs  []*IncidentCommunicationLogDetailDTO `json:"communication_logs,omitempty"`
	IncidentClosure    *IncidentClosureDetailDTO            `json:"incident_closure,omitempty"`
	IncidentPostReview *IncidentPostReviewDetailDTO         `json:"incident_post_review,omitempty"`

	Timelines              []*IncidentTimelineListDTO                        `json:"timelines,omitempty"`
	KPIs                   []*commonDTOs.AssociatedKPIRegisterListDTO        `json:"kpis,omitempty"`
	StrategicObjectives    []*commonDTOs.AssociatedStrategicObjectiveListDTO `json:"strategic_objectives,omitempty"`
	ComplianceRequirements []*commonDTOs.AssociatedClauseListDTO             `json:"compliance_requirements,omitempty"`
	Processes              []*commonDTOs.AssociatedProcessListDTO            `json:"processes,omitempty"`
	Risks                  []*commonDTOs.AssociatedRiskListDTO               `json:"risks,omitempty"`
	Controls               []*commonDTOs.AssociatedControlListDTO            `json:"controls,omitempty"`
	AssociatedDocuments    []*commonDTOs.AssociatedDocumentListDTO           `json:"associated_documents,omitempty"`
	Assets                 []*commonDTOs.AssociatedAssetListDTO              `json:"assets,omitempty"`
}

type IncidentReportPDFDetailDTO struct {
	ID                     uint                            `json:"id"`
	Title                  string                          `json:"title"`
	StartDate              string                          `json:"start_date"`
	EndDate                string                          `json:"end_date"`
	ReportedAt             string                          `json:"reported_at"`
	CreatedUserFirstNameEn string                          `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr string                          `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn  string                          `json:"created_user_last_name_en"`
	CreatedUserLastNameAr  string                          `json:"created_user_last_name_ar"`
	CreatedAt              string                          `json:"created_at"`
	GeneratedAt            string                          `json:"generated_at"`
	ExecutiveSummary       string                          `json:"executive_summary"`
	Scope                  string                          `json:"scope"`
	Objective              string                          `json:"objective"`
	Incidents              []*IncidentReportIncidentPDFDTO `json:"incidents"`
}

type IncidentReportSectionFormDTO struct {
	TitleEn       string  `json:"title_en" binding:"required"`
	TitleAr       string  `json:"title_ar" binding:"required"`
	DescriptionEn *string `json:"description_en"`
	DescriptionAr *string `json:"description_ar"`
}
