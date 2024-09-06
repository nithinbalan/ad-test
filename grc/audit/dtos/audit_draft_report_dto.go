package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type AuditDraftReportListDTO struct {
	ID                          uint    `json:"id"`
	AuditID                     uint    `json:"audit_id"`
	AuditTitle                  string  `json:"audit_title"`
	ReportDate                  string  `json:"report_date"`
	ExecutiveSummary            *string `json:"executive_summary"`
	Conclusion                  *string `json:"conclusion"`
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

type AuditDraftReportFormDTO struct {
	ReportDate       string  `json:"report_date" binding:"required"`
	ExecutiveSummary *string `json:"executive_summary" binding:"required"`
	Conclusion       *string `json:"conclusion"`
}

type AuditDraftReportDetailDTO struct {
	ID                 uint                            `json:"id"`
	AuditID            uint                            `json:"audit_id"`
	Title              string                          `json:"title"`
	ReferenceNumber    string                          `json:"reference_number"`
	ReportDateString   string                          `json:"report_date"`
	PreparedBy         organizationDTOs.UserMinimalDTO `json:"prepared_by"`
	PreparedAt         string                          `json:"prepared_at"`
	ExecutiveSummary   string                          `json:"executive_summary"`
	Auditors           []*AuditPlanAuditMemberDTO      `json:"auditors"`
	TeamLead           organizationDTOs.UserMinimalDTO `json:"team_lead"`
	AuditDraftFindings []*AuditDraftFindingDetailDTO   `json:"draft_findings"`
	Conclusion         string                          `json:"conclusion"`
}

type AuditDraftReportListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}

type AuditDraftReportListResponseDTO struct {
	Items      []*AuditDraftReportListDTO   `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type AuditDraftReportParsedDates struct {
	ReportDate *time.Time
}
