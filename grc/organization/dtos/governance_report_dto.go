package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type GovernanceReportListDTO struct {
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

type GovernanceReportListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
	Filters  map[string][]string
}

type GovernanceReportListResponseDTO struct {
	Items      []*GovernanceReportListDTO   `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type GovernanceReportDetailDTO struct {
	ID                     uint                                 `json:"id"`
	Title                  string                               `json:"title"`
	StartDate              string                               `json:"start_date"`
	EndDate                string                               `json:"end_date"`
	CreatedUserFirstNameEn string                               `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr string                               `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn  string                               `json:"created_user_last_name_en"`
	CreatedUserLastNameAr  string                               `json:"created_user_last_name_ar"`
	CreatedAt              string                               `json:"created_at"`
	Sections               []GovernanceReportSectionResponseDTO `json:"sections"`
	Documents              []GovernanceDocumentResponseDTO      `json:"documents"`
	Meetings               []GovernanceReportMeetingResponseDTO `json:"meetings"`
	Tasks                  []GovernanceReportTaskResponseDTO    `json:"tasks"`
}

type GovernanceReportSectionResponseDTO struct {
	ID            uint   `json:"id"`
	Type          string `json:"type"`
	TitleEn       string `json:"title_en"`
	TitleAr       string `json:"title_ar"`
	DescriptionEn string `json:"description_en"`
	DescriptionAr string `json:"description_ar"`
}

type GovernanceDocumentResponseDTO struct {
	ID                     uint   `json:"id"`
	ReferenceNumber        string `json:"reference_number"`
	DocumentTitleEn        string `json:"document_title_en"`
	DocumentTitleAr        string `json:"document_title_ar"`
	ResponsibleUserNameEn  string `json:"responsible_user_name_en"`
	ResponsibleUserNameAr  string `json:"responsible_user_name_ar"`
	ResponsibleUnitTitleEn string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr string `json:"responsible_unit_title_ar"`
	TotalEmployees         int    `json:"total_employees"`
	AcknowledgedCount      int    `json:"acknowledged_count"`
	PendingCount           int    `json:"pending_count"`
	CreatedByID            uint   `json:"created_by_id"`
	CreatedAt              string `json:"created_at"`
}

type GovernanceReportMeetingResponseDTO struct {
	ID                uint      `json:"id"`
	MeetingTitle      string    `json:"meeting_title"`
	MeetingDate       string    `json:"meeting_date"`
	ParticipantsCount int       `json:"participants_count"`
	KeyDecisions      string    `json:"key_decisions"`
	ActionsAssigned   string    `json:"actions_assigned"`
	CreatedByID       uint      `json:"created_by_id"`
	CreatedAt         time.Time `json:"created_at"`
}

type GovernanceReportTaskResponseDTO struct {
	ID                    uint    `json:"id"`
	GovernanceReportID    uint    `json:"governance_report_id"`
	TaskTitle             string  `json:"task_title"`
	TaskDescription       string  `json:"task_description,omitempty"`
	ResponsibleUserNameEn string  `json:"responsible_user_name_en"`
	ResponsibleUserNameAr string  `json:"responsible_user_name_ar"`
	TaskStatusTitleEn     string  `json:"task_status_title_en"`
	TaskStatusTitleAr     string  `json:"task_status_title_ar"`
	TaskStatusColorCode   string  `json:"task_status_color_code"`
	TargetDate            *string `json:"target_date"`
	CreatedByID           uint    `json:"created_by_id"`
	CreatedAt             string  `json:"created_at"`

	CreatedUserFirstNameEn  string `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr  string `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn   string `json:"created_user_last_name_en"`
	CreatedUserLastNameAr   string `json:"created_user_last_name_ar"`
	CreatedUserEmail        string `json:"created_user_email"`
	CreatedUserInitial      string `json:"created_user_initial"`
	CreatedUserColor        string `json:"created_user_color"`
	CreatedUserProfileImage string `json:"created_user_profile_image_path"`
}

type GovernanceReportFormDTO struct {
	Title               string `json:"title" binding:"required"`
	StartDate           string `json:"start_date" binding:"required"`
	EndDate             string `json:"end_date" binding:"required"`
	OrganizationUnitIDs []uint `json:"organization_unit_ids"`
}
type GovernanceReportSectionFormDTO struct {
	TitleEn       string  `json:"title_en" binding:"required"`
	TitleAr       string  `json:"title_ar" binding:"required"`
	DescriptionEn *string `json:"description_en"`
	DescriptionAr *string `json:"description_ar"`
}

type GovernanceReportParsedDates struct {
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

type GovernanceReportPDFDetailDTO struct {
	ID                     uint                                 `json:"id"`
	Title                  string                               `json:"title"`
	StartDate              string                               `json:"start_date"`
	EndDate                string                               `json:"end_date"`
	CreatedUserFirstNameEn string                               `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr string                               `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn  string                               `json:"created_user_last_name_en"`
	CreatedUserLastNameAr  string                               `json:"created_user_last_name_ar"`
	CreatedAt              string                               `json:"created_at"`
	Sections               []GovernanceReportSectionWithDataDTO `json:"sections"`
}

type GovernanceReportSectionWithDataDTO struct {
	ID            uint                                 `json:"id"`
	Type          string                               `json:"type"`
	TitleEn       string                               `json:"title_en"`
	TitleAr       string                               `json:"title_ar"`
	DescriptionEn string                               `json:"description_en"`
	DescriptionAr string                               `json:"description_ar"`
	Documents     []GovernanceDocumentResponseDTO      `json:"documents,omitempty"`
	Meetings      []GovernanceReportMeetingResponseDTO `json:"meetings,omitempty"`
	Tasks         []GovernanceReportTaskResponseDTO    `json:"tasks,omitempty"`
}
