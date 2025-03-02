package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type AuditFinalReportListDTO struct {
	// Audit Draft Report Information
	ID                        uint   `json:"id"`
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
	AuditFinalReportStatusID        uint   `json:"audit_final_report_status_id"`
	AuditFinalReportStatusTitleEn   string `json:"audit_final_report_status_title_en"`
	AuditFinalReportStatusTitleAr   string `json:"audit_final_report_status_title_ar"`
	AuditFinalReportStatusColorCode string `json:"audit_final_report_status_color_code"`
	AuditFinalReportStatusType      string `json:"audit_final_report_status_type"`
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
type AuditFinalReportFormDTO struct {
	ExecutiveSummary *string `json:"executive_summary" binding:"required"`
	Conclusion       *string `json:"conclusion"`
}

type AuditFinalReportDetailDTO struct {
	ID                         uint                             `json:"id"`
	AuditID                    uint                             `json:"audit_id"`
	ReferenceNumber            string                           `json:"reference_number"`
	Title                      string                           `json:"title"`
	Description                *string                          `json:"description"`
	Comments                   *string                          `json:"comments"`
	NextReviewLevel            int                              `json:"next_review_level"`
	PreparedAt                 string                           `json:"prepared_at"`
	PreparedBy                 organizationDTOs.UserMinimalDTO  `json:"prepared_by"`
	MainSections               []*AuditFinalReportSectionDTO    `json:"main_sections"` // Main sections with subsections
	AuditFinalReportStatus     *AuditFinalReportStatusDTO       `json:"audit_final_report_status"`
	AuditFinalReportRecipients []*AuditFinalReportRecipientDTO  `json:"audit_final_report_recipients"`
	PublishedByUser            *organizationDTOs.UserMinimalDTO `json:"published_by"`
	PublishedAt                *string                          `json:"published_at"`
	UpdatedByUser              *organizationDTOs.UserMinimalDTO `json:"updated_by"`
	UpdatedAt                  *string                          `json:"updated_at"`
}

type AuditFinalReportSectionDTO struct {
	ID                        uint                                 `json:"id"`
	TitleEn                   string                               `json:"title_en"`
	TitleAr                   string                               `json:"title_ar"`
	AuditFinalReportSectionID *uint                                `json:"audit_final_report_section_id"`
	DescriptionEn             *string                              `json:"description_en"`
	DescriptionAr             *string                              `json:"description_ar"`
	Type                      string                               `json:"type"` // One of: overview, finding-summary, findings, conclusion, custom
	Order                     *float64                             `json:"order"`
	SubSections               []*AuditFinalReportSectionDTO        `json:"sub_sections"` // Recursive structure for subsections
	Findings                  []*AuditFinalReportFindingDTO        `json:"findings,omitempty"`
	FindingCategorySummary    []*FindingCategorySummaryDTO         `json:"finding_category_summary,omitempty"`
	Appendices                []*AuditFinalReportAppendixDetailDTO `json:"appendices"`
	CreatedByUser             organizationDTOs.UserMinimalDTO      `json:"created_by"`
	UpdatedByUser             *organizationDTOs.UserMinimalDTO     `json:"updated_by"`
	CreatedAt                 string                               `json:"created_at"`
	UpdatedAt                 *string                              `json:"updated_at"`
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
	Comments  *string `json:"comments"`
}
type AuditFinalReportRecipientDTO struct {
	ID                 *uint                           `json:"id"`
	AuditFinalReportID uint                            `json:"audit_final_report_id"`
	UserID             uint                            `json:"user_id"`
	Type               string                          `json:"type"`
	User               organizationDTOs.UserMinimalDTO `json:"user"`
}

type AuditFinalReportListQueryParamsDTO struct {
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

type AuditFinalReportListResponseDTO struct {
	Items      []*AuditFinalReportListDTO   `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type FindingCategorySummaryDTO struct {
	Category *AuditFindingCategoryDTO `json:"category"`
	Count    int                      `json:"count"`
}
