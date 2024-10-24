package dtos

import (
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
)

type AuditFinalReportFormDTO struct {
	ExecutiveSummary *string `json:"executive_summary" binding:"required"`
	Conclusion       *string `json:"conclusion"`
}

type AuditFinalReportDetailDTO struct {
	ID                     uint                             `json:"id"`
	AuditID                uint                             `json:"audit_id"`
	ReferenceNumber        string                           `json:"reference_number"`
	Title                  string                           `json:"title"`
	PreparedAt             string                           `json:"prepared_at"`
	PreparedBy             organizationDTOs.UserMinimalDTO  `json:"prepared_by"`
	MainSections           []*AuditFinalReportSectionDTO    `json:"main_sections"` // Main sections with subsections
	AuditFinalReportStatus *AuditFinalReportStatusDTO       `json:"audit_final_report_status"`
	UpdatedByUser          *organizationDTOs.UserMinimalDTO `json:"updated_by"`
	UpdatedAt              *string                          `json:"updated_at"`
}

type AuditFinalReportSectionDTO struct {
	ID                        uint                             `json:"id"`
	TitleEn                   string                           `json:"title_en"`
	TitleAr                   string                           `json:"title_ar"`
	AuditFinalReportSectionID *uint                            `json:"audit_final_report_section_id"`
	DescriptionEn             *string                          `json:"description_en"`
	DescriptionAr             *string                          `json:"description_ar"`
	Type                      string                           `json:"type"` // One of: overview, finding-summary, findings, conclusion, custom
	Order                     *float64                         `json:"order"`
	SubSections               []*AuditFinalReportSectionDTO    `json:"sub_sections"` // Recursive structure for subsections
	Findings                  []*AuditFinalReportFindingDTO    `json:"findings,omitempty"`
	CreatedByUser             organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser             *organizationDTOs.UserMinimalDTO `json:"updated_by"`
	CreatedAt                 string                           `json:"created_at"`
	UpdatedAt                 *string                          `json:"updated_at"`
}

type AuditFinalReportFindingDTO struct {
	ID                     uint                             `json:"id"`
	Title                  string                           `json:"title"`
	Description            *string                          `json:"description"`
	SeverityLevelID        uint                             `json:"severity_level_id"`
	AuditFindingCategoryID uint                             `json:"audit_finding_category_id"`
	Recommendations        *string                          `json:"recommendations"`
	AuditFindingCategory   AuditFindingCategoryDTO          `json:"audit_finding_category"`
	SeverityLevel          masterDTOs.SeverityLevelDTO      `json:"severity_level"`
	AuditableItem          AuditableItemDetailDTO           `json:"auditable_item"`
	ActionPlans            []*AuditFindingActionPlanDTO     `json:"action_plans"` // Nested action plans for each finding
	CreatedByUser          organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser          *organizationDTOs.UserMinimalDTO `json:"updated_by"`
	CreatedAt              string                           `json:"created_at"`
	UpdatedAt              *string                          `json:"updated_at"`
}

type AuditFindingActionPlanDTO struct {
	ID                 uint                             `json:"id"`
	Title              string                           `json:"title"`
	Description        *string                          `json:"description"`
	ResponsibleUserID  uint                             `json:"responsible_user_id"`
	TargetDate         string                           `json:"target_date"`
	ActionPlanStatusID uint                             `json:"action_plan_status_id"`
	CreatedByUser      organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser      *organizationDTOs.UserMinimalDTO `json:"updated_by"`
	CreatedAt          string                           `json:"created_at"`
	UpdatedAt          *string                          `json:"updated_at"`
}
