package dtos

import (
	organizationDTOs "jethur/grc/organization/dtos"
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
	Exceptions         []*ProjectSecurityGateReportExceptionDTO `json:"exceptions"`
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
	SubmittedByEn           string                                `json:"submitted_by_en"`
	SubmittedByAr           string                                `json:"submitted_by_ar"`
	SubmittedAtFormatted    string                                `json:"submitted_at_formatted"`
	ApprovedByEn            string                                `json:"approved_by_en"`
	ApprovedByAr            string                                `json:"approved_by_ar"`
	ApprovedAtFormatted     string                                `json:"approved_at_formatted"`
	Reviews                 []*ProjectSecurityGateReportReviewDTO `json:"reviews"`
}

type ProjectSecurityGateReportReviewDTO struct {
	StepNumber            int    `json:"step_number"`
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

// Cross-gate exceptions report
type ProjectSecurityGateExceptionReportDTO struct {
	GeneratedAt   string                                       `json:"generated_at"`
	TotalCount    int                                          `json:"total_count"`
	ApprovedCount int                                          `json:"approved_count"`
	PendingCount  int                                          `json:"pending_count"`
	CriticalCount int                                          `json:"critical_count"`
	Exceptions    []*ProjectSecurityGateExceptionReportItemDTO `json:"exceptions"`
}

type ProjectSecurityGateExceptionReportItemDTO struct {
	ProjectTitleEn         string `json:"project_title_en"`
	ProjectTitleAr         string `json:"project_title_ar"`
	StageTitleEn           string `json:"stage_title_en"`
	StageTitleAr           string `json:"stage_title_ar"`
	ReadinessStatusTitleEn string `json:"readiness_status_title_en"`
	ReadinessStatusTitleAr string `json:"readiness_status_title_ar"`
	BusinessOwnerEn        string `json:"business_owner_en"`
	BusinessOwnerAr        string `json:"business_owner_ar"`

	TitleEn       string `json:"title_en"`
	TitleAr       string `json:"title_ar"`
	IsCritical    bool   `json:"is_critical"`
	Justification string `json:"justification"`

	WorkflowStatusTitleEn         string `json:"workflow_status_title_en"`
	WorkflowStatusTitleAr         string `json:"workflow_status_title_ar"`
	WorkflowStatusTextColor       string `json:"workflow_status_text_color"`
	WorkflowStatusBorderColor     string `json:"workflow_status_border_color"`
	WorkflowStatusBackgroundColor string `json:"workflow_status_background_color"`

	SubmittedByEn        string `json:"submitted_by_en"`
	SubmittedByAr        string `json:"submitted_by_ar"`
	SubmittedAtFormatted string `json:"submitted_at_formatted"`
	ApprovedByEn         string `json:"approved_by_en"`
	ApprovedByAr         string `json:"approved_by_ar"`
	ApprovedAtFormatted  string `json:"approved_at_formatted"`

	Reviews []*ProjectSecurityGateReportReviewDTO `json:"reviews"`
}
