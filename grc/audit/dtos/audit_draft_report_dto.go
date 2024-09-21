package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)
type AuditDraftReportListDTO struct {
	// Audit Draft Report Information
	ID               uint   `json:"id"`
	ReportDate       string `json:"report_date"`
	ExecutiveSummary string `json:"executive_summary"`
	Conclusion       string `json:"conclusion"`
	// Submitted User Information
	SubmittedAt               string `json:"submitted_at"`
	SubmittedUserID           uint   `json:"submitted_user_id"`
	SubmittedUserFirstNameEn  string `json:"submitted_user_first_name_en"`
	SubmittedUserFirstNameAr  string `json:"submitted_user_first_name_ar"`
	SubmittedUserLastNameEn   string `json:"submitted_user_last_name_en"`
	SubmittedUserLastNameAr   string `json:"submitted_user_last_name_ar"`
	SubmittedUserInitial      string `json:"submitted_user_initial"`
	SubmittedUserColor        string `json:"submitted_user_color"`
	SubmittedUserEmail        string `json:"submitted_user_email"`
	SubmittedUserProfileImage string `json:"submitted_user_profile_image_path"`
	// Audit Draft Report Status
	AuditDraftReportStatusID        uint   `json:"audit_draft_report_status_id"`
	AuditDraftReportStatusTitleEn   string `json:"audit_draft_report_status_title_en"`
	AuditDraftReportStatusTitleAr   string `json:"audit_draft_report_status_title_ar"`
	AuditDraftReportStatusColorCode string `json:"audit_draft_report_status_color_code"`
	AuditDraftReportStatusType      string `json:"audit_draft_report_status_type"`
	// Audit Information
	AuditID     uint   `json:"audit_id"`
	AuditPlanID uint   `json:"audit_plan_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	// Audit Plan Information
	AuditPlanTitle       string `json:"audit_plan_title"`
	AuditPlanDescription string `json:"audit_plan_description"`
	AuditPlanPurpose     string `json:"audit_plan_purpose"`
	AuditPlanObjective   string `json:"audit_plan_objective"`
	// Audit Status
	AuditStatusTitleEn   string `json:"audit_status_title_en"`
	AuditStatusTitleAr   string `json:"audit_status_title_ar"`
	AuditStatusType      string `json:"audit_status_type"`
	AuditStatusColorCode string `json:"audit_status_color_code"`
	// Audit Category
	AuditCategoryTitleEn string `json:"audit_category_title_en"`
	AuditCategoryTitleAr string `json:"audit_category_title_ar"`
	// Team Lead Information
	TeamLeadID           uint   `json:"team_lead_id"`
	TeamLeadFirstNameEn  string `json:"team_lead_first_name_en"`
	TeamLeadFirstNameAr  string `json:"team_lead_first_name_ar"`
	TeamLeadLastNameEn   string `json:"team_lead_last_name_en"`
	TeamLeadLastNameAr   string `json:"team_lead_last_name_ar"`
	TeamLeadEmail        string `json:"team_lead_email"`
	TeamLeadInitial      string `json:"team_lead_initial"`
	TeamLeadColor        string `json:"team_lead_color"`
	TeamLeadProfileImage string `json:"team_lead_profile_image_path"`
	// Created User Information
	CreatedUserID           uint   `json:"created_user_id"`
	CreatedUserFirstNameEn  string `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr  string `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn   string `json:"created_user_last_name_en"`
	CreatedUserLastNameAr   string `json:"created_user_last_name_ar"`
	CreatedUserInitial      string `json:"created_user_initial"`
	CreatedUserColor        string `json:"created_user_color"`
	CreatedUserEmail        string `json:"created_user_email"`
	CreatedUserProfileImage string `json:"created_user_profile_image_path"`
	// Updated User Information
	UpdatedUserID           uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn  string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr  string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn   string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr   string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial      string `json:"updated_user_initial"`
	UpdatedUserColor        string `json:"updated_user_color"`
	UpdatedUserEmail        string `json:"updated_user_email"`
	UpdatedUserProfileImage string `json:"updated_user_profile_image_path"`
}
type AuditDraftReportListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	ReportDate  string
}

type AuditDraftReportListResponseDTO struct {
	Items      []*AuditDraftReportListDTO   `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}
type AuditDraftReportFormDTO struct {
	ReportDate       string  `json:"report_date" binding:"required"`
	ExecutiveSummary *string `json:"executive_summary" binding:"required"`
	Conclusion       *string `json:"conclusion"`
}

type AuditDraftReportDetailDTO struct {
	ID                     uint                                `json:"id"`
	AuditID                uint                                `json:"audit_id"`
	Title                  string                              `json:"title"`
	ReferenceNumber        string                              `json:"reference_number"`
	ReportDateString       string                              `json:"report_date"`
	PreparedAt             string                              `json:"prepared_at"`
	PreparedBy             organizationDTOs.UserMinimalDTO     `json:"prepared_by"`
	ExecutiveSummary       string                              `json:"executive_summary"`
	Auditors               []*AuditPlanAuditMemberDTO          `json:"auditors"`
	TeamLead               organizationDTOs.UserMinimalDTO     `json:"team_lead"`
	AuditDraftFindings     []*AuditDraftFindingReportDetailDTO `json:"draft_findings"`
	Conclusion             string                              `json:"conclusion"`
	SubmittedByUser        *organizationDTOs.UserMinimalDTO    `json:"submitted_by"`
	SubmittedAt            *string                             `json:"submitted_at"`
	AuditDraftReportStatus *AuditDraftReportStatusDTO          `json:"audit_draft_report_status"`
	CreatedByUser          organizationDTOs.UserMinimalDTO     `json:"created_by"`
	UpdatedByUser          *organizationDTOs.UserMinimalDTO    `json:"updated_by"`
}
type AuditDraftReportParsedDates struct {
	ReportDate *time.Time
}
type AuditDraftFindingReportDetailDTO struct {
	ID                             uint                                      `json:"id"`
	Title                          string                                    `json:"title"`
	Description                    *string                                   `json:"description"`
	AuditID                        uint                                      `json:"audit_id"`
	AuditableItemID                *uint                                     `json:"auditable_item_id"`
	AuditChecklistID               *uint                                     `json:"audit_checklist_id"`
	AuditFindingCategoryID         uint                                      `json:"audit_finding_category_id"`
	SeverityLevelID                uint                                      `json:"severity_level_id"`
	Recommendations                *string                                   `json:"recommendations,omitempty"`
	AuditDraftFindingStatusID      uint                                      `json:"audit_draft_finding_status_id"`
	AuditableItem                  AuditableItemDetailDTO                    `json:"auditable_item"`
	AuditFindingCategory           AuditFindingCategoryDTO                   `json:"audit_finding_category"`
	SeverityLevel                  masterDTOs.SeverityLevelDTO               `json:"severity_level"`
	AuditDraftFindingStatus        AuditDraftFindingStatusDTO                `json:"audit_draft_finding_status"`
	CreatedAt                      time.Time                                 `json:"created_at"`
	UpdatedAt                      *time.Time                                `json:"updated_at"`
	CreatedByUser                  organizationDTOs.UserMinimalDTO           `json:"created_by"`
	UpdatedByUser                  *organizationDTOs.UserMinimalDTO          `json:"updated_by"`
	Documents                      []*AuditDraftFindingDocumentDTO           `json:"documents"`
	AuditDraftFindingConversations []*AuditDraftFindingConversationDetailDTO `json:"conversations"`
	AuditDraftFindingActionPlans   []*AuditDraftFindingActionPlanDetailDTO   `json:"action_plans"`
}
