package dtos

import (
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
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
	Description            *string                          `json:"description"`
	NextReviewLevel        int                              `json:"next_review_level"`
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
	ID                       uint                             `json:"id"`
	Title                    string                           `json:"title"`
	Description              *string                          `json:"description"`
	SeverityLevelID          uint                             `json:"severity_level_id"`
	AuditFindingCategoryID   uint                             `json:"audit_finding_category_id"`
	Recommendations          *string                          `json:"recommendations"`
	AuditFindingCategory     AuditFindingCategoryDTO          `json:"audit_finding_category"`
	SeverityLevel            masterDTOs.SeverityLevelDTO      `json:"severity_level"`
	AuditableItem            AuditableItemDetailDTO           `json:"auditable_item"`
	ActionPlans              []*AuditFindingActionPlanDTO     `json:"action_plans"`
	DepartmentResponse       *string                          `json:"department_response"`
	DepartmentResponseByID   *uint                            `json:"department_response_by_id"`
	DepartmentResponseByUser *organizationDTOs.UserMinimalDTO `json:"department_response_by"`
	DepartmentResponseAt     *string                          `json:"department_response_at"`
	CreatedByUser            organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser            *organizationDTOs.UserMinimalDTO `json:"updated_by"`
	CreatedAt                string                           `json:"created_at"`
	UpdatedAt                *string                          `json:"updated_at"`
}

type AuditFindingActionPlanDTO struct {
	ID                 uint                             `json:"id"`
	Title              string                           `json:"title"`
	Description        *string                          `json:"description"`
	ResponsibleUserID  uint                             `json:"responsible_user_id"`
	ResponsibleUser    organizationDTOs.UserMinimalDTO  `json:"responsible_user"`
	TargetDate         string                           `json:"target_date"`
	ActionPlanStatusID uint                             `json:"action_plan_status_id"`
	ActionPlanStatus   AuditFindingActionPlanStatusDTO  `json:"action_plan_status"`
	CreatedByUser      organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser      *organizationDTOs.UserMinimalDTO `json:"updated_by"`
	CreatedAt          string                           `json:"created_at"`
	UpdatedAt          *string                          `json:"updated_at"`
}

type AuditReportEditHistoryDTO struct {
	ID                 *uint   `json:"id"`
	AuditID            uint    `json:"audit_id"`
	AuditFinalReportID uint    `json:"audit_final_report_id"`
	ActionType         string  `json:"action_type"`
	Comment            *string `json:"comment"`
}

type AuditFinalReportLastEditSummaryDTO struct {
	LatestActivity ActivityDTO                        `json:"latest_activity"`
	CurrentLevel   *int                               `json:"current_level"`
	Participants   []*organizationDTOs.UserMinimalDTO `json:"participants"`
	TotalComments  *int                               `json:"total_comments"`
}

type ActivityDTO struct {
	Activity  string                           `json:"activity"` // "checkout" or "checkin"
	By        *organizationDTOs.UserMinimalDTO `json:"by"`
	Comment   *string                          `json:"comment"`
	Timestamp time.Time                        `json:"timestamp"`
}

type AuditFinalReportPublishFormDTO struct {
	CCUserIDs []uint  `json:"cc_user_ids" binding:"required"`
	ToUserID  uint    `json:"to_user_id" binding:"required"`
	Comment   *string `json:"comment"`
}
