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
	ID                         uint                                  `json:"id" mapstructure:"ID"`
	ReferenceNumber            *string                               `json:"reference_number" mapstructure:"ReferenceNumber"`
	IncidentCategoryTitleEn    *string                               `json:"incident_category_title_en" mapstructure:"IncidentCategoryTitleEn"`
	IncidentCategoryTitleAr    *string                               `json:"incident_category_title_ar" mapstructure:"IncidentCategoryTitleAr"`
	IncidentSubCategoryTitleEn *string                               `json:"incident_sub_category_title_en" mapstructure:"IncidentSubCategoryTitleEn"`
	IncidentSubCategoryTitleAr *string                               `json:"incident_sub_category_title_ar" mapstructure:"IncidentSubCategoryTitleAr"`
	SeverityLevelTitleEn       *string                               `json:"severity_level_title_en" mapstructure:"SeverityLevelTitleEn"`
	SeverityLevelTitleAr       *string                               `json:"severity_level_title_ar" mapstructure:"SeverityLevelTitleAr"`
	SeverityLevelColorCode     *string                               `json:"severity_level_color_code" mapstructure:"SeverityLevelColorCode"`
	IncidentAt                 *time.Time                            `json:"incident_at" mapstructure:"IncidentAt"`
	ReportedAt                 *time.Time                            `json:"reported_at" mapstructure:"ReportedAt"`
	IncidentLocation           string                                `json:"incident_location" mapstructure:"IncidentLocation"`
	IncidentStatusTitleEn      *string                               `json:"incident_status_title_en" mapstructure:"IncidentStatusTitleEn"`
	IncidentStatusTitleAr      *string                               `json:"incident_status_title_ar" mapstructure:"IncidentStatusTitleAr"`
	ResponsibleUserNameEn      string                                `json:"responsible_user_name_en" mapstructure:"ResponsibleUserNameEn"`
	ResponsibleUserNameAr      string                                `json:"responsible_user_name_ar" mapstructure:"ResponsibleUserNameAr"`
	ReportedUserNameEn         string                                `json:"reported_user_name_en" mapstructure:"ReportedUserNameEn"`
	ReportedUserNameAr         string                                `json:"reported_user_name_ar" mapstructure:"ReportedUserNameAr"`
	ClosedDate                 *time.Time                            `json:"closed_date" mapstructure:"ClosedDate"`
	Documents                  []*IncidentReportDocumentDTO          `json:"documents" mapstructure:"Documents"`
	RootCauseAnalyses          []*IncidentReportRootCauseAnalysisDTO `json:"root_cause_analyses" mapstructure:"RootCauseAnalyses"`
	ActionPlans                []*IncidentReportActionPlanDTO        `json:"action_plans" mapstructure:"ActionPlans"`
	IncidentRegister           *IncidentRegisterDTO                  `json:"incident_register,omitempty" mapstructure:"IncidentRegister"`
	Investigation              *IncidentInvestigationDetailDTO       `json:"investigation,omitempty" mapstructure:"-"`
	CommunicationLogs          []*IncidentCommunicationLogDetailDTO  `json:"communication_logs,omitempty" mapstructure:"-"`
	IncidentClosure            *IncidentClosureDetailDTO             `json:"incident_closure,omitempty" mapstructure:"-"`
	IncidentPostReview         *IncidentPostReviewDetailDTO          `json:"incident_post_review,omitempty" mapstructure:"-"`
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

type IncidentReportSectionFormDTO struct {
	TitleEn       string  `json:"title_en" binding:"required"`
	TitleAr       string  `json:"title_ar" binding:"required"`
	DescriptionEn *string `json:"description_en"`
	DescriptionAr *string `json:"description_ar"`
}
