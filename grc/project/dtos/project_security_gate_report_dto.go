package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

// Carries both languages to feed either PDF template
type ProjectSecurityGateReportDTO struct {
	ProjectTitleEn     string                                   `json:"project_title_en"`
	ProjectTitleAr     string                                   `json:"project_title_ar"`
	Stage              *ProjectSecurityGateStageDTO             `json:"stage"`
	TemplateVersion    *ProjectSecurityGateTemplateVersionDTO   `json:"template_version"`
	BusinessOwner      *organizationDTOs.UserMinimalDTO         `json:"business_owner"`
	ReadinessStatus    *ProjectSecurityGateReadinessStatusDTO   `json:"readiness_status"`
	Notes              *string                                  `json:"notes"`
	Classifications    []*ProjectSecurityGateClassificationDTO  `json:"classifications"`
	CreatedByUser      *organizationDTOs.UserMinimalDTO         `json:"created_by"`
	CreatedAtFormatted string                                   `json:"created_at_formatted"`
	UpdatedByUser      *organizationDTOs.UserMinimalDTO         `json:"updated_by"`
	UpdatedAtFormatted string                                   `json:"updated_at_formatted"`
	GeneratedAt        string                                   `json:"generated_at"`
	Summary            *ProjectSecurityGateChecklistSummaryDTO  `json:"summary"`
	Checklists         []*ProjectSecurityGateReportChecklistDTO `json:"checklists"`
	Exceptions         []*ProjectSecurityGateReportExceptionDTO `json:"exceptions,omitempty"`
	Controls           []*ProjectSecurityGateReportControlDTO   `json:"controls"`
	Risks              []*ProjectSecurityGateReportRiskDTO      `json:"risks"`
}

type ProjectSecurityGateReportChecklistDTO struct {
	ChecklistID             uint   `json:"checklist_id"`
	TitleEn                 string `json:"title_en"`
	TitleAr                 string `json:"title_ar"`
	IsCritical              bool   `json:"is_critical"`
	IsException             bool   `json:"is_exception"`
	ResponseStatusTitleEn   string `json:"response_status_title_en"`
	ResponseStatusTitleAr   string `json:"response_status_title_ar"`
	ResponseStatusColorCode string `json:"response_status_color_code"`
	Comment                 string `json:"comment"`
	RespondedByEn           string `json:"responded_by_en"`
	RespondedByAr           string `json:"responded_by_ar"`
	RespondedAtFormatted    string `json:"responded_at_formatted"`
}

type ProjectSecurityGateReportExceptionDTO struct {
	ChecklistID             uint                                  `json:"checklist_id"`
	TitleEn                 string                                `json:"title_en"`
	TitleAr                 string                                `json:"title_ar"`
	IsCritical              bool                                  `json:"is_critical"`
	Justification           string                                `json:"justification"`
	WorkflowStatusTitleEn   string                                `json:"workflow_status_title_en"`
	WorkflowStatusTitleAr   string                                `json:"workflow_status_title_ar"`
	WorkflowStatusColorCode string                                `json:"workflow_status_color_code"`
	WorkflowStatusType      string                                `json:"workflow_status_type"`
	SubmittedByEn           string                                `json:"submitted_by_en"`
	SubmittedByAr           string                                `json:"submitted_by_ar"`
	SubmittedAtFormatted    string                                `json:"submitted_at_formatted"`
	ApprovedByEn            string                                `json:"approved_by_en"`
	ApprovedByAr            string                                `json:"approved_by_ar"`
	ApprovedAtFormatted     string                                `json:"approved_at_formatted"`
	Reviews                 []*ProjectSecurityGateReportReviewDTO `json:"reviews"`

	// New structured fields for exception report
	WorkflowStatus *WorkflowStatusDTO               `json:"workflow_status,omitempty"`
	SubmittedBy    *organizationDTOs.UserMinimalDTO `json:"submitted_by,omitempty"`
	SubmittedAt    *time.Time                       `json:"submitted_at,omitempty"`
	ApprovedBy     *organizationDTOs.UserMinimalDTO `json:"approved_by,omitempty"`
	ApprovedAt     *time.Time                       `json:"approved_at,omitempty"`
}

type ProjectSecurityGateReportReviewDTO struct {
	Level                 int    `json:"level"`
	StatusTitleEn         string `json:"status_title_en"`
	StatusTitleAr         string `json:"status_title_ar"`
	StatusType            string `json:"status_type"`
	StatusTextColor       string `json:"status_text_color"`
	StatusBorderColor     string `json:"status_border_color"`
	StatusBackgroundColor string `json:"status_background_color"`
	ActionDateLabelEn     string `json:"action_date_label_en"`
	ActionDateLabelAr     string `json:"action_date_label_ar"`
	ActionByLabelEn       string `json:"action_by_label_en"`
	ActionByLabelAr       string `json:"action_by_label_ar"`
	CreatedAtString       string `json:"created_at_string"`
	UserInitial           string `json:"user_initial"`
	UserInitialAr         string `json:"user_initial_ar"`
	UserFullNameEn        string `json:"user_full_name_en"`
	UserFullNameAr        string `json:"user_full_name_ar"`
}

type ProjectSecurityGateReportControlDTO struct {
	ChecklistTitleEn string `json:"checklist_title_en"`
	ChecklistTitleAr string `json:"checklist_title_ar"`
	ReferenceNumber  string `json:"reference_number"`
	Code             string `json:"code"`
	TitleEn          string `json:"title_en"`
	TitleAr          string `json:"title_ar"`
	Description      string `json:"description"`
}

type ProjectSecurityGateReportRiskDTO struct {
	ChecklistTitleEn          string `json:"checklist_title_en"`
	ChecklistTitleAr          string `json:"checklist_title_ar"`
	ReferenceNumber           string `json:"reference_number"`
	RiskCode                  string `json:"risk_code"`
	Title                     string `json:"title"`
	RiskStatusTitleEn         string `json:"risk_status_title_en"`
	RiskStatusTitleAr         string `json:"risk_status_title_ar"`
	RiskCategoryTitleEn       string `json:"risk_category_title_en"`
	RiskCategoryTitleAr       string `json:"risk_category_title_ar"`
	ResponsibleUnitTitleEn    string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr    string `json:"responsible_unit_title_ar"`
	ResponsibleUserEn         string `json:"responsible_user_en"`
	ResponsibleUserAr         string `json:"responsible_user_ar"`
	RegistrationDateFormatted string `json:"registration_date_formatted"`
}

// Exception report payload: stored on create, then served by details and rendered into the pdf
type ProjectSecurityGateExceptionReportDataDTO struct {
	Title         string                                   `json:"title"`
	GeneratedByEn string                                   `json:"generated_by_en"`
	GeneratedByAr string                                   `json:"generated_by_ar"`
	GeneratedAt   string                                   `json:"generated_at"`
	ExportedAt    string                                   `json:"exported_at"`
	Projects      []*ProjectSecurityGateReportProjectDTO   `json:"projects"`
	Sections      []*ProjectSecurityGateReportSectionDTO   `json:"sections"`
	SecurityGates []*ProjectSecurityGateExceptionReportDTO `json:"security_gates"`
	TotalCount    int                                      `json:"total_count"`
	ApprovedCount int                                      `json:"approved_count"`
	PendingCount  int                                      `json:"pending_count"`
	CriticalCount int                                      `json:"critical_count"`
}

type ProjectSecurityGateExceptionReportDTO struct {
	ProjectTitleEn     string                                  `json:"project_title_en"`
	ProjectTitleAr     string                                  `json:"project_title_ar"`
	Stage              *ProjectSecurityGateStageDTO            `json:"stage"`
	TemplateVersion    *ProjectSecurityGateTemplateVersionDTO  `json:"template_version"`
	BusinessOwner      *organizationDTOs.UserMinimalDTO        `json:"business_owner"`
	ReadinessStatus    *ProjectSecurityGateReadinessStatusDTO  `json:"readiness_status"`
	Notes              *string                                 `json:"notes"`
	Classifications    []*ProjectSecurityGateClassificationDTO `json:"classifications"`
	CreatedByUser      *organizationDTOs.UserMinimalDTO        `json:"created_by"`
	CreatedAtFormatted string                                  `json:"created_at_formatted"`
	UpdatedByUser      *organizationDTOs.UserMinimalDTO        `json:"updated_by"`
	UpdatedAtFormatted string                                  `json:"updated_at_formatted"`
	Exceptions         []*ProjectSecurityGateExceptionItemDTO  `json:"exceptions"`
}

type ProjectSecurityGateExceptionItemDTO struct {
	ChecklistID    uint                                  `json:"checklist_id"`
	TitleEn        string                                `json:"title_en"`
	TitleAr        string                                `json:"title_ar"`
	IsCritical     bool                                  `json:"is_critical"`
	Justification  string                                `json:"justification"`
	WorkflowStatus *WorkflowStatusDTO                    `json:"workflow_status"`
	SubmittedBy    *organizationDTOs.UserMinimalDTO      `json:"submitted_by"`
	SubmittedAt    *time.Time                            `json:"submitted_at"`
	ApprovedBy     *organizationDTOs.UserMinimalDTO      `json:"approved_by"`
	ApprovedAt     *time.Time                            `json:"approved_at"`
	Reviews        []*ProjectSecurityGateReportReviewDTO `json:"reviews"`
}

type WorkflowStatusDTO struct {
	ID              uint   `json:"id"`
	TitleEn         string `json:"title_en"`
	TitleAr         string `json:"title_ar"`
	Type            string `json:"type"`
	ColorCode       string `json:"color_code"`
	Status          string `json:"status"`
	TextColor       string `json:"text_color"`
	BorderColor     string `json:"border_color"`
	BackgroundColor string `json:"background_color"`
}

// Checklist report payload: stored on create, then served by details and rendered into the pdf
type ProjectSecurityGateReportDataDTO struct {
	Title         string                                 `json:"title"`
	GeneratedByEn string                                 `json:"generated_by_en"`
	GeneratedByAr string                                 `json:"generated_by_ar"`
	GeneratedAt   string                                 `json:"generated_at"`
	ExportedAt    string                                 `json:"exported_at"`
	Projects      []*ProjectSecurityGateReportProjectDTO `json:"projects"`
	Sections      []*ProjectSecurityGateReportSectionDTO `json:"sections"`
	SecurityGates []*ProjectSecurityGateReportDTO        `json:"security_gates"`
}

type ProjectSecurityGateReportSectionDTO struct {
	ID            uint   `json:"id"`
	Type          string `json:"type"`
	TitleEn       string `json:"title_en"`
	TitleAr       string `json:"title_ar"`
	DescriptionEn string `json:"description_en"`
	DescriptionAr string `json:"description_ar"`
	OrderNum      int    `json:"order_num"`
}

type ProjectSecurityGateReportSectionFormDTO struct {
	TitleEn       string  `json:"title_en" binding:"required,max=255"`
	TitleAr       string  `json:"title_ar" binding:"required,max=255"`
	DescriptionEn *string `json:"description_en"`
	DescriptionAr *string `json:"description_ar"`
}

type ProjectSecurityGateReportFormDTO struct {
	Title             string `json:"title" binding:"required,max=255"`
	SecurityGateIDs   []uint `json:"security_gate_ids" binding:"required,min=1"`
	ResponseStatusIDs []uint `json:"response_status_ids" binding:"required,min=1"`
}

type ProjectSecurityGateExceptionReportFormDTO struct {
	Title           string `json:"title" binding:"required,max=255"`
	SecurityGateIDs []uint `json:"security_gate_ids" binding:"required,min=1"`
}

type ProjectSecurityGateReportListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Filters  map[string][]string
	Page     int
	PageSize int
}

type ProjectSecurityGateReportListDTO struct {
	ID                          uint                      `json:"id"`
	Title                       string                    `json:"title"`
	CreatedAt                   string                    `json:"created_at"`
	CreatedUserID               uint                      `json:"created_user_id"`
	CreatedUserFirstNameEn      string                    `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string                    `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string                    `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string                    `json:"created_user_last_name_ar"`
	CreatedUserInitial          string                    `json:"created_user_initial"`
	CreatedUserColor            string                    `json:"created_user_color"`
	CreatedUserEmail            string                    `json:"created_user_email"`
	CreatedUserProfileImagePath string                    `json:"created_user_profile_image_path"`
	SecurityGates               []*SecurityGateTitleDTO   `json:"security_gates"`
	ResponseStatuses            []*ResponseStatusTitleDTO `json:"response_statuses"`
}

type SecurityGateTitleDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
}

type ResponseStatusTitleDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
}

type ProjectSecurityGateReportListResponseDTO struct {
	Items      []*ProjectSecurityGateReportListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO        `json:"pagination"`
}

// Report record details, the row metadata plus the frozen payload
type ProjectSecurityGateReportDetailsDTO struct {
	ID               uint                                             `json:"id"`
	CreatedAt        time.Time                                        `json:"created_at"`
	CreatedByUser    *organizationDTOs.UserMinimalDTO                 `json:"created_by"`
	ResponseStatuses []*ProjectSecurityGateChecklistResponseStatusDTO `json:"response_statuses"`
	*ProjectSecurityGateReportDataDTO
}

type ProjectSecurityGateExceptionReportDetailsDTO struct {
	ID            uint                             `json:"id"`
	CreatedAt     time.Time                        `json:"created_at"`
	CreatedByUser *organizationDTOs.UserMinimalDTO `json:"created_by"`
	*ProjectSecurityGateExceptionReportDataDTO
}

type ProjectSecurityGateReportProjectDTO struct {
	ID             uint   `json:"id"`
	ProjectTitleEn string `json:"project_title_en"`
	ProjectTitleAr string `json:"project_title_ar"`
	StageTitleEn   string `json:"stage_title_en"`
	StageTitleAr   string `json:"stage_title_ar"`
}
