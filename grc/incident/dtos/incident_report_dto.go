package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"time"
)

type IncidentReportDetailDTO struct {
	ID                     uint                                              `json:"id"`
	ReferenceNumber        string                                            `json:"reference_number"`
	Title                  string                                            `json:"title"`
	Description            string                                            `json:"description"`
	IncidentAtFormatted    string                                            `json:"incident_at_formatted"`
	ReportedAtFormatted    string                                            `json:"reported_at_formatted"`
	Location               *string                                           `json:"location"`
	IncidentCategory       IncidentCategoryDTO                               `json:"incident_category"`
	IncidentType           IncidentTypeDTO                                   `json:"incident_type"`
	SeverityLevel          *masterDTOs.SeverityLevelDTO                      `json:"severity_level"`
	IncidentStatus         IncidentStatusDTO                                 `json:"incident_status"`
	ReportedByUser         organizationDTOs.UserMinimalDTO                   `json:"reported_by"`
	Documents              []*IncidentDocumentDTO                            `json:"documents"`
	Controls               []*commonDTOs.AssociatedControlListDTO            `json:"controls"`
	Risks                  []*commonDTOs.AssociatedRiskListDTO               `json:"risks"`
	Assets                 []*commonDTOs.AssociatedAssetListDTO              `json:"assets"`
	Processes              []*commonDTOs.AssociatedProcessListDTO            `json:"processes"`
	Timelines              []*IncidentTimelineListDTO                        `json:"timelines"`
	ResponsibleUnit        settingsDTOs.OrganizationUnitDTO                  `json:"responsible_unit"`
	ResponsibleUser        organizationDTOs.UserMinimalDTO                   `json:"responsible_user"`
	ResponsibleUserID      uint                                              `json:"responsible_user_id" binding:"required" gorm:"not null"`
	SLARegisters           []*IncidentSLARegisterDTO                         `json:"sla_registers"`
	AssociatedDocuments    []*commonDTOs.AssociatedDocumentListDTO           `json:"associated_documents"`
	CreatedAt              time.Time                                         `json:"created_at"`
	UpdatedAt              *time.Time                                        `json:"updated_at"`
	CreatedByUser          organizationDTOs.UserMinimalDTO                   `json:"created_by"`
	UpdatedByUser          organizationDTOs.UserMinimalDTO                   `json:"updated_by"`
	GeneratedAt            string                                            `json:"generated_at"`
	CommunicationLogs      []*IncidentCommunicationLogDetailDTO              `json:"communication_log"`
	ActionPlans            []*IncidentActionPlanListDTO                      `json:"action_plan"`
	IncidentClosure        *IncidentClosureDetailDTO                         `json:"incident_closure"`
	IncidentPostReview     *IncidentPostReviewDetailDTO                      `json:"incident_post_review"`
	RootCauseAnalyses      []*IncidentRootCauseAnalysisListDTO               `json:"root_cause_analysis"`
	StrategicObjectives    []*commonDTOs.AssociatedStrategicObjectiveListDTO `json:"strategic_objectives"`
	KPIs                   []*commonDTOs.AssociatedKPIRegisterListDTO        `json:"kpi"`
	ComplianceRequirements []*commonDTOs.AssociatedClauseListDTO             `json:"compliance_clause"`
	Investigation          *IncidentInvestigationDetailDTO                   `json:"investigation"`
	AppUrl                 string                                            `json:"url"`
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

	IncidentSubCategoryTitleEn *string `json:"incident_sub_category_title_en"`
	IncidentSubCategoryTitleAr *string `json:"incident_sub_category_title_ar"`

	SeverityTitleEn   string  `json:"severity_title_en"`
	SeverityTitleAr   string  `json:"severity_title_ar"`
	SeverityColorCode *string `json:"severity_color_code"`

	RealTime         time.Time `json:"real_time"`
	DetectedTime     time.Time `json:"detected_time"`
	IncidentLocation string    `json:"incident_location"`

	IncidentStatusTitleEn string `json:"incident_status_title_en"`
	IncidentStatusTitleAr string `json:"incident_status_title_ar"`
	ResponsibleUserNameEn string `json:"responsible_user_name_en"`
	ResponsibleUserNameAr string `json:"responsible_user_name_ar"`
	ReportedUserNameEn    string `json:"reported_user_name_en"`
	ReportedUserNameAr    string `json:"reported_user_name_ar"`

	ContainmentStrategy *string `json:"containment_strategy"`
	NextSteps           *string `json:"next_steps"`

	ResolvedDate     *time.Time `json:"resolved_date"`
	RecoveryMeasures *string    `json:"recovery_measures"`

	AssignedUserInitial     *string `json:"assigned_user_initial"`
	AssignedUserFirstNameEn *string `json:"assigned_user_first_name_en"`
	AssignedUserLastNameEn  *string `json:"assigned_user_last_name_en"`
	AssignedUserFirstNameAr *string `json:"assigned_user_first_name_ar"`
	AssignedUserLastNameAr  *string `json:"assigned_user_last_name_ar"`
}

type IncidentReportDocumentDTO struct {
	ID                       uint    `json:"id"`
	IncidentDocumentID       uint    `json:"incident_document_id"`
	IncidentReportIncidentID uint    `json:"incident_report_incident_id"`
	URL                      string  `json:"url"`
	Name                     *string `json:"name"`
	Size                     int     `json:"size"`
	Token                    string  `json:"token"`
	OriginalName             string  `json:"original_name"`
	Type                     string  `json:"type"`
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

type IncidentReportDetailsDTO struct {
	ID        uint                               `json:"id"`
	Title     string                             `json:"title"`
	StartDate *string                            `json:"start_date"`
	EndDate   *string                            `json:"end_date"`
	CreatedAt time.Time                          `json:"created_at"`
	CreatedBy organizationDTOs.UserMinimalDTO    `json:"created_by"`
	Sections  []*IncidentReportSectionDTO        `json:"sections"`
	Incidents []*IncidentReportIncidentDetailDTO `json:"incidents"`
}

type IncidentReportSectionDTO struct {
	ID            uint   `json:"id"`
	Type          string `json:"type"`
	TitleEn       string `json:"title_en"`
	TitleAr       string `json:"title_ar"`
	DescriptionEn string `json:"description_en"`
	DescriptionAr string `json:"description_ar"`
	OrderNum      int    `json:"order_num"`
}

type IncidentReportIncidentDetailDTO struct {
	ID                         uint                                  `json:"id"`
	ReferenceNumber            *string                               `json:"reference_number"`
	IncidentCategoryTitleEn    string                                `json:"incident_category_title_en"`
	IncidentCategoryTitleAr    string                                `json:"incident_category_title_ar"`
	IncidentSubCategoryTitleEn *string                               `json:"incident_sub_category_title_en"`
	IncidentSubCategoryTitleAr *string                               `json:"incident_sub_category_title_ar"`
	SeverityTitleEn            string                                `json:"severity_title_en"`
	SeverityTitleAr            string                                `json:"severity_title_ar"`
	IncidentAt                 *time.Time                            `json:"incident_at"`
	ReportedAt                 *time.Time                            `json:"reported_at"`
	IndicatorsOfCompromise     *string                               `json:"indicators_of_compromise"`
	IncidentStatusTitleEn      string                                `json:"incident_status_title_en"`
	IncidentStatusTitleAr      string                                `json:"incident_status_title_ar"`
	ContainmentStrategy        *string                               `json:"containment_strategy"`
	NextSteps                  *string                               `json:"next_steps"`
	ResolvedDate               *time.Time                            `json:"resolved_date"`
	RecoveryMeasures           *string                               `json:"recovery_measures"`
	AssignedUserInitial        *string                               `json:"assigned_user_initial"`
	AssignedUserFirstNameEn    *string                               `json:"assigned_user_first_name_en"`
	AssignedUserLastNameEn     *string                               `json:"assigned_user_last_name_en"`
	AssignedUserFirstNameAr    *string                               `json:"assigned_user_first_name_ar"`
	AssignedUserLastNameAr     *string                               `json:"assigned_user_last_name_ar"`
	Documents                  []*IncidentReportDocumentDTO          `json:"documents"`
	RootCauseAnalyses          []*IncidentReportRootCauseAnalysisDTO `json:"root_cause_analyses"`
	ActionPlans                []*IncidentReportActionPlanDTO        `json:"action_plans"`
	Investigation              *IncidentInvestigationDetailDTO       `json:"investigation"`
	CommunicationLogs          []*IncidentCommunicationLogDetailDTO  `json:"communication_logs"`
	IncidentClosure            *IncidentClosureDetailDTO             `json:"incident_closure"`
	IncidentPostReview         *IncidentPostReviewDetailDTO          `json:"incident_post_review"`
}

type IncidentReportRootCauseAnalysisDTO struct {
	ID        uint                             `json:"id"`
	Why       uint                             `json:"why"`
	Question  string                           `json:"question"`
	Answer    string                           `json:"answer"`
	CreatedBy organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedBy *organizationDTOs.UserMinimalDTO `json:"updated_by"`
	CreatedAt time.Time                        `json:"created_at"`
	UpdatedAt *time.Time                       `json:"updated_at"`
}

type IncidentReportActionPlanDTO struct {
	ID          uint                             `json:"id"`
	Title       string                           `json:"title"`
	Description *string                          `json:"description"`
	TargetDate  *time.Time                       `json:"target_date"`
	StatusTitle string                           `json:"status_title"`
	CreatedBy   *organizationDTOs.UserMinimalDTO `json:"created_by"`
	CreatedAt   time.Time                        `json:"created_at"`
}

type IncidentReportSectionFormDTO struct {
	TitleEn       string  `json:"title_en" binding:"required"`
	TitleAr       string  `json:"title_ar" binding:"required"`
	DescriptionEn *string `json:"description_en"`
	DescriptionAr *string `json:"description_ar"`
}
