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
	Type 		   string `json:"type" binding:"required,oneof=individual date-range"`
	IncidentIDs    *[]uint `json:"incident_ids"`
	Title  	 	   string `json:"title" binding:"required"`
	StartDate      *string  `json:"start_date"`
	EndDate        *string  `json:"end_date"`
}

type IncidentReportParsedDates struct {
	StartDate *time.Time
	EndDate   *time.Time
}

type IncidentReportIncidentDTO struct {
	ID              uint   `json:"id"`
	ReferenceNumber string `json:"reference_number"`

	IncidentCategoryTitleEn string `json:"incident_category_title_en"`
	IncidentCategoryTitleAr string `json:"incident_category_title_ar"`

	IncidentSubCategoryTitleEn *string `json:"incident_sub_category_title_en"`
	IncidentSubCategoryTitleAr *string `json:"incident_sub_category_title_ar"`

	SeverityTitleEn string `json:"severity_title_en"`
	SeverityTitleAr string `json:"severity_title_ar"`

	RealTime     time.Time `json:"real_time"`
	DetectedTime time.Time `json:"detected_time"`

	IncidentStatusTitleEn string `json:"incident_status_title_en"`
	IncidentStatusTitleAr string `json:"incident_status_title_ar"`

	ContainmentStrategy *string `json:"containment_strategy"`
	NextSteps           *string `json:"next_steps"`

	ResolvedDate     *time.Time `json:"resolved_date"`
	RecoveryMeasures *string    `json:"recovery_measures"`

	AssignedUserInitial        *string `json:"assigned_user_initial"`
	AssignedUserFirstNameEn    *string `json:"assigned_user_first_name_en"`
	AssignedUserLastNameEn     *string `json:"assigned_user_last_name_en"`
	AssignedUserFirstNameAr    *string `json:"assigned_user_first_name_ar"`
	AssignedUserLastNameAr     *string `json:"assigned_user_last_name_ar"`
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