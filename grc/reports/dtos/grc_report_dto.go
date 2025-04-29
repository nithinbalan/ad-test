package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type GRCReportListDTO struct {
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

type GRCReportListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
	Filters  map[string][]string
}

type GRCReportListResponseDTO struct {
	Items      []*GRCReportListDTO          `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type GRCReportDetailDTO struct {
	ID                     uint                          `json:"id"`
	Title                  string                        `json:"title"`
	StartDate              string                        `json:"start_date"`
	EndDate                string                        `json:"end_date"`
	CreatedUserFirstNameEn string                        `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr string                        `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn  string                        `json:"created_user_last_name_en"`
	CreatedUserLastNameAr  string                        `json:"created_user_last_name_ar"`
	CreatedAt              string                        `json:"created_at"`
	Sections               []GRCReportSectionResponseDTO `json:"sections"`
	// Summary                []GRCGovernanceSummaryResponseDTO  `json:"documents"`
	Documents []GRCGovernanceDocumentResponseDTO `json:"documents"`
	Tasks     []GRCGovernanceTaskResponseDTO     `json:"tasks"`
}

type GRCReportSectionResponseDTO struct {
	ID            uint   `json:"id"`
	Type          string `json:"type"`
	TitleEn       string `json:"title_en"`
	TitleAr       string `json:"title_ar"`
	DescriptionEn string `json:"description_en"`
	DescriptionAr string `json:"description_ar"`
}

type GRCGovernanceSummaryResponseDTO struct {
	ID                 uint     `json:"id"`
	ReferenceNumber    string   `json:"reference_number"`
	PoliciesInSystem   string   `json:"policies_in_system"`
	TotalPolicies      int      `json:"total_policies"`
	RolesDefined       int      `json:"roles_defined"`
	PolicyApprovedRate *float64 `json:"policy_approved_rate"`
	CreatedByID        uint     `json:"created_by_id"`
	CreatedAt          string   `json:"created_at"`
}

type GRCGovernanceDocumentResponseDTO struct {
	ID                      uint   `json:"id"`
	ReferenceNumber         string `json:"reference_number"`
	DocumentTitleEn         string `json:"document_title_en"`
	DocumentTitleAr         string `json:"document_title_ar"`
	ResponsibleUserNameEn   string `json:"responsible_user_name_en"`
	ResponsibleUserNameAr   string `json:"responsible_user_name_ar"`
	ResponsibleUnitTitleEn  string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr  string `json:"responsible_unit_title_ar"`
	TotalEmployees          int    `json:"total_employees"`
	AcknowledgedCount       int    `json:"acknowledged_count"`
	PendingCount            int    `json:"pending_count"`
	DocumentStatusTitleEn   string `json:"document_status_title_en"`
	DocumentStatusTitleAr   string `json:"document_status_title_ar"`
	DocumentStatusColorCode string `json:"document_status_color_code"`
	CreatedByID             uint   `json:"created_by_id"`
	CreatedAt               string `json:"created_at"`
}

type GRCGovernanceTaskResponseDTO struct {
	ID                      uint    `json:"id"`
	GRCReportID             uint    `json:"grc_report_id"`
	TaskTitle               string  `json:"task_title"`
	TaskDescription         string  `json:"task_description,omitempty"`
	ResponsibleUserNameEn   string  `json:"responsible_user_name_en"`
	ResponsibleUserNameAr   string  `json:"responsible_user_name_ar"`
	TaskStatusTitleEn       string  `json:"task_status_title_en"`
	TaskStatusTitleAr       string  `json:"task_status_title_ar"`
	TaskStatusColorCode     string  `json:"task_status_color_code"`
	SeverityLevelTitleEn    string  `json:"severity_level_title_en"`
	SeverityLevelTitleAr    string  `json:"severity_level_title_ar"`
	SeverityLevelColorCode  string  `json:"severity_level_color_code"`
	TargetDate              *string `json:"target_date"`
	CreatedByID             uint    `json:"created_by_id"`
	CreatedAt               string  `json:"created_at"`
	CreatedUserFirstNameEn  string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr  string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn   string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr   string  `json:"created_user_last_name_ar"`
	CreatedUserEmail        string  `json:"created_user_email"`
	CreatedUserInitial      string  `json:"created_user_initial"`
	CreatedUserColor        string  `json:"created_user_color"`
	CreatedUserProfileImage string  `json:"created_user_profile_image_path"`
}

type GRCReportFormDTO struct {
	Title         string `json:"title" binding:"required"`
	StartDate     string `json:"start_date" binding:"required"`
	EndDate       string `json:"end_date" binding:"required"`
	ReviewUserID  uint   `json:"review_user_id"`
	ApproveUserID uint   `json:"approve_user_id"`
}

type GRCReportSectionFormDTO struct {
	TitleEn       string  `json:"title_en" binding:"required"`
	TitleAr       string  `json:"title_ar" binding:"required"`
	DescriptionEn *string `json:"description_en"`
	DescriptionAr *string `json:"description_ar"`
}

type GRCReportParsedDates struct {
	StartDate *time.Time
	EndDate   *time.Time
}

type DocumentAcknowledgmentSummaryDTO struct {
	TotalUsers             int64   `json:"total_users"`
	AcknowledgedUsers      int64   `json:"acknowledged_users"`
	NotAcknowledgedUsers   int64   `json:"not_acknowledged_users"`
	AcknowledgedPercent    float64 `json:"acknowledged_percent"`
	NotAcknowledgedPercent float64 `json:"not_acknowledged_percent"`
}
