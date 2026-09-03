package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
)

type AuditDashboardQueryParamsDTO struct {
	DateFrom string
	DateTo   string
	Year     int
}

type AuditDashboardListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
	Export   bool
	DateFrom string
	DateTo   string
	Year     int
}

type BreakdownReferenceDTO struct {
	ID        uint   `json:"id"`
	Type      string `json:"type"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
}

type AuditStatusBreakdownItemDTO struct {
	Status     BreakdownReferenceDTO `json:"status"`
	Count      int64                 `json:"count"`
	Percentage float64               `json:"percentage"`
}

type AuditStatusBreakdownDTO struct {
	Total    int64                          `json:"total"`
	Statuses []AuditStatusBreakdownItemDTO `json:"statuses"`
}

type AuditSummaryDTO struct {
	TotalAudits      int64 `json:"total_audits"`
	CompletedAudits  int64 `json:"completed_audits"`
	InProgressAudits int64 `json:"in_progress_audits"`
	NotStartedAudits int64 `json:"not_started_audits"`
}

type AuditFindingsSummaryDTO struct {
	FindingsTotal        int64   `json:"findings_total"`
	FindingsOpen         int64   `json:"findings_open"`
	FindingsCriticalHigh int64   `json:"findings_critical_or_high"`
	ReassessmentPending  int64   `json:"reassessment_pending"`
	FindingClosureRate   float64 `json:"finding_closure_rate"`
}

type AuditFindingsSeverityBreakdownItemDTO struct {
	SeverityLevel BreakdownReferenceDTO `json:"severity_level"`
	Count         int64                 `json:"count"`
	Percentage    float64               `json:"percentage"`
}

type AuditFindingsSeverityBreakdownDTO struct {
	Total          int64                                    `json:"total"`
	SeverityLevels []AuditFindingsSeverityBreakdownItemDTO `json:"severity_levels"`
}

type AuditFindingsBySeverityLevelListDTO struct {
	ID                          uint    `json:"id"`
	AuditID                     uint    `json:"audit_id"`
	AuditTitle                  string  `json:"audit_title"`
	ReferenceNumber             string  `json:"reference_number"`
	Title                       string  `json:"title"`
	AuditableItemTitleEn        string  `json:"auditable_item_title_en"`
	AuditableItemTitleAr        string  `json:"auditable_item_title_ar"`
	AuditFindingCategoryTitleEn string  `json:"audit_finding_category_title_en"`
	AuditFindingCategoryTitleAr string  `json:"audit_finding_category_title_ar"`
	SeverityLevelTitleEn        string  `json:"severity_level_title_en"`
	SeverityLevelTitleAr        string  `json:"severity_level_title_ar"`
	SeverityLevelColorCode      string  `json:"severity_level_color_code"`
	AuditFindingStatusTitleEn   string  `json:"audit_finding_status_title_en"`
	AuditFindingStatusTitleAr   string  `json:"audit_finding_status_title_ar"`
	AuditFindingStatusColorCode string  `json:"audit_finding_status_color_code"`
	CreatedAt                   string  `json:"created_at"`
	UpdatedAt                   *string `json:"updated_at"`
	CreatedUserID               uint    `json:"created_user_id"`
	CreatedUserFirstNameEn      string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string  `json:"created_user_last_name_ar"`
	CreatedUserInitial          string  `json:"created_user_initial"`
	CreatedUserColor            string  `json:"created_user_color"`
	CreatedUserEmail            string  `json:"created_user_email"`
	CreatedUserProfileImagePath string  `json:"created_user_profile_image_path"`
	UpdatedUserID               *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn      *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          *string `json:"updated_user_initial"`
	UpdatedUserColor            *string `json:"updated_user_color"`
	UpdatedUserEmail            *string `json:"updated_user_email"`
	UpdatedUserProfileImagePath *string `json:"updated_user_profile_image_path"`
}

type AuditFindingsBySeverityLevelListResponseDTO struct {
	Items      []*AuditFindingsBySeverityLevelListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO           `json:"pagination"`
}

type AuditFindingsTrendMonthDTO struct {
	Month   string `json:"month"`
	MonthEn string `json:"month_en"`
	MonthAr string `json:"month_ar"`
	New     int64  `json:"new"`
	Closed  int64  `json:"closed"`
}

type AuditOverdueActionPlanListDTO struct {
	ID                               uint   `json:"id"`
	AuditFindingID                   uint   `json:"audit_finding_id"`
	AuditFindingReferenceNumber      string `json:"audit_finding_reference_number"`
	AuditFindingTitle                string `json:"audit_finding_title"`
	AuditFindingDescription          string `json:"audit_finding_description"`
	Title                            string `json:"title"`
	SeverityLevelTitleEn             string `json:"severity_level_title_en"`
	SeverityLevelTitleAr             string `json:"severity_level_title_ar"`
	SeverityLevelColorCode           string `json:"severity_level_color_code"`
	ResponsibleUserID                uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn       string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr       string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn        string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr        string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial           string `json:"responsible_user_initial"`
	ResponsibleUserColor             string `json:"responsible_user_color"`
	ResponsibleUserEmail             string `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath  string `json:"responsible_user_profile_image_path"`
	ResponsibleUserDepartmentTitleEn string `json:"responsible_user_department_title_en"`
	ResponsibleUserDepartmentTitleAr string `json:"responsible_user_department_title_ar"`
	TargetDate                       string `json:"target_date"`
	OverdueDays                      int    `json:"overdue_days"`
	ActionPlanStatusTitleEn          string `json:"action_plan_status_title_en"`
	ActionPlanStatusTitleAr          string `json:"action_plan_status_title_ar"`
	ActionPlanStatusColorCode        string `json:"action_plan_status_color_code"`
	CreatedUserFirstNameEn           string `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr           string `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn            string `json:"created_user_last_name_en"`
	CreatedUserLastNameAr            string `json:"created_user_last_name_ar"`
	CreatedUserInitial               string `json:"created_user_initial"`
	CreatedUserColor                 string `json:"created_user_color"`
	CreatedUserEmail                 string `json:"created_user_email"`
	CreatedUserProfileImagePath      string `json:"created_user_profile_image_path"`
	CreatedAt                        string `json:"created_at"`
}

type AuditOverdueActionPlanListResponseDTO struct {
	Items      []*AuditOverdueActionPlanListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO     `json:"pagination"`
}
