package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type ProjectClosureReportDTO struct {
	ID                    uint                             `json:"id"`
	OverallProjectOutcome *string                          `json:"overall_project_outcome"`
	ClosureReason         *string                          `json:"closure_reason"`
	ClosureComments       *string                          `json:"closure_comments"`
	LessonLearned         *string                          `json:"lesson_learned"`
	ApproverID            uint                             `json:"approver_id"`
	IsApproved            bool                             `json:"is_approved"`
	ApproverUser          *organizationDTOs.UserMinimalDTO `json:"approver_user"`
	CreatedByUser         organizationDTOs.UserMinimalDTO  `json:"created_by"`
	CreatedAt             time.Time                        `json:"created_at"`
	CreatedAtFormatted    string                           `json:"created_at_formatted"`
}

type ProjectReportDetailDTO struct {
	ID                      uint                                              `json:"id"`
	TitleEn                 string                                            `json:"title_en"`
	TitleAr                 string                                            `json:"title_ar"`
	DescriptionEn           *string                                           `json:"description_en"`
	DescriptionAr           *string                                           `json:"description_ar"`
	Budget                  *float64                                          `json:"budget"`
	Scope                   *string                                           `json:"scope"`
	NumberOfResources       *int                                              `json:"number_of_resources"`
	CompletionPercentage    *float64                                          `json:"completion_percentage"`
	StartDateString         string                                            `json:"start_date"`
	EndDateString           string                                            `json:"end_date"`
	StartDateFormatted      string                                            `json:"start_date_formatted"`
	EndDateFormatted        string                                            `json:"end_date_formatted"`
	ResponsibleUser         organizationDTOs.UserMinimalDTO                   `json:"responsible_user"`
	ProjectType             ProjectTypeDTO                                    `json:"project_type"`
	ProjectStatus           ProjectStatusDTO                                  `json:"project_status"`
	ProjectPhaseID          *uint                                             `json:"project_phase_id"`
	ProjectPhase            *ProjectPhaseDTO                                  `json:"project_phase"`
	ProjectClassificationID *uint                                             `json:"project_classification_id"`
	ProjectClassification   *ProjectClassificationDTO                         `json:"project_classification"`
	IncidentActionPlanID    *uint                                             `json:"incident_action_plan_id"`
	IncidentActionPlan      *ProjectIncidentActionPlanDTO                     `json:"incident_action_plan" mapstructure:"-"`
	CreatedAt               time.Time                                         `json:"created_at"`
	UpdatedAt               *time.Time                                        `json:"updated_at"`
	CreatedByUser           organizationDTOs.UserMinimalDTO                   `json:"created_by"`
	UpdatedByUser           *organizationDTOs.UserMinimalDTO                  `json:"updated_by"`
	GeneratedAt             string                                            `json:"generated_at"`
	Milestones              []*ProjectMilestoneListDTO                        `json:"milestones"`
	Controls                []*commonDTOs.AssociatedControlListDTO            `json:"controls"`
	Risks                   []*commonDTOs.AssociatedRiskListDTO               `json:"risks"`
	Processes               []*commonDTOs.AssociatedProcessListDTO            `json:"processes"`
	StrategicObjectives     []*commonDTOs.AssociatedStrategicObjectiveListDTO `json:"strategic_objectives"`
	KPIs                    []*commonDTOs.AssociatedKPIRegisterListDTO        `json:"kpis"`
	AssociatedDocuments     []*commonDTOs.AssociatedDocumentListDTO           `json:"associated_documents"`
	Assets                  []*commonDTOs.AssociatedAssetListDTO              `json:"assets"`
	Incidents               []*commonDTOs.AssociatedIncidentListDTO           `json:"incidents"`
	ComplianceRequirements  []*commonDTOs.AssociatedClauseListDTO             `json:"compliance_requirements"`
	ActionPlans             []*ProjectActionPlanListDTO                       `json:"action_plans"`
	ChangeRequests          []*ProjectChangeRequestListDTO                    `json:"change_requests"`
	DecisionLogs            []*ProjectDecisionLogListDTO                      `json:"decision_logs"`
	SupportingDocuments     []*ProjectSupportingDocumentListDTO               `json:"supporting_documents"`
	Closure                 *ProjectClosureReportDTO                          `json:"closure"`
}

type ProjectReportRequestDTO struct {
	Type       string  `json:"type" binding:"required,oneof=individual date-range"`
	ProjectIDs *[]uint `json:"project_ids"`
	Title      string  `json:"title" binding:"required"`
	StartDate  *string `json:"start_date"`
	EndDate    *string `json:"end_date"`
}

type ProjectReportListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
	Filters  map[string][]string
}

type ProjectReportListDTO struct {
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

type ProjectReportListResponseDTO struct {
	Items      []*ProjectReportListDTO      `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type ProjectReportSectionDTO struct {
	ID            uint   `json:"id"`
	Type          string `json:"type"`
	TitleEn       string `json:"title_en"`
	TitleAr       string `json:"title_ar"`
	DescriptionEn string `json:"description_en"`
	DescriptionAr string `json:"description_ar"`
	OrderNum      int    `json:"order_num"`
}

type ProjectReportSectionFormDTO struct {
	TitleEn       string  `json:"title_en" binding:"required"`
	TitleAr       string  `json:"title_ar" binding:"required"`
	DescriptionEn *string `json:"description_en"`
	DescriptionAr *string `json:"description_ar"`
}

type ProjectReportDetailsDTO struct {
	ID            uint                             `json:"id"`
	Title         string                           `json:"title"`
	StartDate     *time.Time                       `json:"start_date"`
	EndDate       *time.Time                       `json:"end_date"`
	CreatedAt     time.Time                        `json:"created_at"`
	CreatedByUser *organizationDTOs.UserMinimalDTO `json:"created_by"`
	Sections      []*ProjectReportSectionDTO       `json:"sections"`
	Projects      []*ProjectReportProjectDetailDTO `json:"projects"`
}

type ProjectReportProjectDetailDTO struct {
	ID                           uint                                              `json:"id"`
	ProjectID                    uint                                              `json:"project_id"`
	TitleEn                      string                                            `json:"title_en"`
	TitleAr                      string                                            `json:"title_ar"`
	DescriptionEn                *string                                           `json:"description_en"`
	DescriptionAr                *string                                           `json:"description_ar"`
	Budget                       *float64                                          `json:"budget"`
	Scope                        *string                                           `json:"scope"`
	StartDate                    string                                            `json:"start_date"`
	EndDate                      string                                            `json:"end_date"`
	ProjectTypeTitleEn           *string                                           `json:"project_type_title_en"`
	ProjectTypeTitleAr           *string                                           `json:"project_type_title_ar"`
	ProjectStatusTitleEn         *string                                           `json:"project_status_title_en"`
	ProjectStatusTitleAr         *string                                           `json:"project_status_title_ar"`
	ProjectStatusColorCode       *string                                           `json:"project_status_color_code,omitempty"`
	ProjectPhaseTitleEn          *string                                           `json:"project_phase_title_en"`
	ProjectPhaseTitleAr          *string                                           `json:"project_phase_title_ar"`
	ProjectPhaseColorCode        *string                                           `json:"project_phase_color_code,omitempty"`
	ProjectClassificationTitleEn *string                                           `json:"project_classification_title_en"`
	ProjectClassificationTitleAr *string                                           `json:"project_classification_title_ar"`
	ResponsibleUserNameEn        string                                            `json:"responsible_user_name_en"`
	ResponsibleUserNameAr        string                                            `json:"responsible_user_name_ar"`
	IncidentActionPlanID         *uint                                             `json:"incident_action_plan_id,omitempty"`
	IncidentActionPlanTitle      *string                                           `json:"incident_action_plan_title,omitempty"`
	IncidentActionPlan           *ProjectIncidentActionPlanDTO                     `json:"incident_action_plan,omitempty"`
	CreatedAt                    string                                            `json:"created_at"`
	Documents                    []*ProjectReportDocumentDTO                       `json:"documents"`
	Milestones                   []*ProjectMilestoneListDTO                        `json:"milestones,omitempty"`
	Controls                     []*commonDTOs.AssociatedControlListDTO            `json:"controls,omitempty"`
	Risks                        []*commonDTOs.AssociatedRiskListDTO               `json:"risks,omitempty"`
	Processes                    []*commonDTOs.AssociatedProcessListDTO            `json:"processes,omitempty"`
	StrategicObjectives          []*commonDTOs.AssociatedStrategicObjectiveListDTO `json:"strategic_objectives,omitempty"`
	KPIs                         []*commonDTOs.AssociatedKPIRegisterListDTO        `json:"kpis,omitempty"`
	AssociatedDocuments          []*commonDTOs.AssociatedDocumentListDTO           `json:"associated_documents,omitempty"`
	Assets                       []*commonDTOs.AssociatedAssetListDTO              `json:"assets,omitempty"`
	Incidents                    []*commonDTOs.AssociatedIncidentListDTO           `json:"incidents,omitempty"`
	ComplianceRequirements       []*commonDTOs.AssociatedClauseListDTO             `json:"compliance_requirements,omitempty"`
	ActionPlans                  []*ProjectActionPlanListDTO                       `json:"action_plans,omitempty"`
	ChangeRequests               []*ProjectChangeRequestListDTO                    `json:"change_requests,omitempty"`
	DecisionLogs                 []*ProjectDecisionLogListDTO                      `json:"decision_logs,omitempty"`
	SupportingDocuments          []*ProjectSupportingDocumentListDTO               `json:"supporting_documents,omitempty"`
	Closure                      *ProjectClosureReportDTO                          `json:"closure,omitempty"`
}

type ProjectReportDocumentDTO struct {
	ID           uint    `json:"id"`
	Name         *string `json:"name"`
	Size         *int    `json:"size"`
	Token        *string `json:"token"`
	OriginalName *string `json:"original_name"`
	Type         *string `json:"type"`
}

type ProjectReportPDFDetailDTO struct {
	ID                     uint                             `json:"id"`
	Title                  string                           `json:"title"`
	StartDate              string                           `json:"start_date"`
	EndDate                string                           `json:"end_date"`
	ReportedAt             string                           `json:"reported_at"`
	CreatedUserFirstNameEn string                           `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr string                           `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn  string                           `json:"created_user_last_name_en"`
	CreatedUserLastNameAr  string                           `json:"created_user_last_name_ar"`
	CreatedAt              string                           `json:"created_at"`
	GeneratedAt            string                           `json:"generated_at"`
	ExecutiveSummary       string                           `json:"executive_summary"`
	Scope                  string                           `json:"scope"`
	Objective              string                           `json:"objective"`
	Projects               []*ProjectReportProjectDetailDTO `json:"projects"`
}
