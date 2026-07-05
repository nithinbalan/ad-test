package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type AwarenessReportListDTO struct {
	ID                             uint       `json:"id"`
	ReferenceNumber                string     `json:"reference_number"`
	ReportSourceChannel            *string    `json:"report_source_channel"`
	ReportDate                     *time.Time `json:"report_date"`
	ResultedFromCampaign           bool       `json:"resulted_from_campaign"`
	ClosureDate                    *time.Time `json:"closure_date"`
	Result                         *string    `json:"result"`
	AwarenessReportTypeID          uint       `json:"awareness_report_type_id"`
	AwarenessReportTypeTitleEn     string     `json:"awareness_report_type_title_en"`
	AwarenessReportTypeTitleAr     string     `json:"awareness_report_type_title_ar"`
	SubCategoryID                  *uint      `json:"sub_category_id"`
	SubCategoryTitleEn             string     `json:"sub_category_title_en"`
	SubCategoryTitleAr             string     `json:"sub_category_title_ar"`
	SeverityID                     *uint      `json:"severity_id"`
	SeverityTitleEn                string     `json:"severity_title_en"`
	SeverityTitleAr                string     `json:"severity_title_ar"`
	SeverityColorCode              string     `json:"severity_color_code"`
	AwarenessHandlingStatusID      *uint      `json:"awareness_handling_status_id"`
	AwarenessHandlingStatusTitleEn string     `json:"awareness_handling_status_title_en"`
	AwarenessHandlingStatusTitleAr string     `json:"awareness_handling_status_title_ar"`
	HandlingOwnerID                *uint      `json:"handling_owner_id"`
	HandlingOwnerFirstNameEn       string     `json:"handling_owner_first_name_en"`
	HandlingOwnerFirstNameAr       string     `json:"handling_owner_first_name_ar"`
	HandlingOwnerLastNameEn        string     `json:"handling_owner_last_name_en"`
	HandlingOwnerLastNameAr        string     `json:"handling_owner_last_name_ar"`
	HandlingOwnerInitial           string     `json:"handling_owner_initial"`
	CreatedAt                      string     `json:"created_at"`
	UpdatedAt                      string     `json:"updated_at"`
	CreatedUserID                  uint       `json:"created_user_id"`
	CreatedUserFirstNameEn         string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr         string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn          string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr          string     `json:"created_user_last_name_ar"`
	CreatedUserInitial             string     `json:"created_user_initial"`
	CreatedUserColor               string     `json:"created_user_color"`
	UpdatedUserID                  *uint      `json:"updated_user_id"`
	UpdatedUserFirstNameEn         string     `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr         string     `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn          string     `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr          string     `json:"updated_user_last_name_ar"`
	UpdatedUserInitial             string     `json:"updated_user_initial"`
	UpdatedUserColor               string     `json:"updated_user_color"`
}

type AwarenessReportListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
	Filters  map[string][]string
	Export   bool
	DateFrom string
	DateTo   string
}

type AwarenessReportListResponseDTO struct {
	Items      []*AwarenessReportListDTO    `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type AwarenessReportFormDTO struct {
	AwarenessReportTypeID     uint       `json:"awareness_report_type_id" binding:"required"`
	SubCategoryID             *uint      `json:"sub_category_id"`
	ReportSourceChannel       *string    `json:"report_source_channel"`
	SeverityID                *uint      `json:"severity_id"`
	ReportDate                *time.Time `json:"report_date"`
	ResultedFromCampaign      bool       `json:"resulted_from_campaign"`
	AwarenessHandlingStatusID *uint      `json:"awareness_handling_status_id"`
	HandlingOwnerID           *uint      `json:"handling_owner_id"`
	ClosureDate               *time.Time `json:"closure_date"`
	Result                    *string    `json:"result"`
}

type AwarenessReportTypeDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
}

type AwarenessReportDetailResponseDTO struct {
	ID                      uint                             `json:"id"`
	ReferenceNumber         string                           `json:"reference_number"`
	ReportSourceChannel     *string                          `json:"report_source_channel"`
	ReportDate              *time.Time                       `json:"report_date"`
	ResultedFromCampaign    bool                             `json:"resulted_from_campaign"`
	ClosureDate             *time.Time                       `json:"closure_date"`
	Result                  *string                          `json:"result"`
	AwarenessReportType     AwarenessReportTypeDTO           `json:"awareness_report_type"`
	SubCategory             *AwarenessReportSubCategoryDTO   `json:"sub_category"`
	SeverityLevel           *masterDTOs.SeverityLevelDTO     `json:"severity_level"`
	AwarenessHandlingStatus *AwarenessHandlingStatusDTO      `json:"awareness_handling_status"`
	HandlingOwner           *organizationDTOs.UserMinimalDTO `json:"handling_owner"`
	CreatedAt               time.Time                        `json:"created_at"`
	UpdatedAt               *time.Time                       `json:"updated_at"`
	CreatedBYUser           organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedBYUser           organizationDTOs.UserMinimalDTO  `json:"updated_by"`
}


type AwarenessReportParsedDates struct {
	ReportDate  *time.Time
	ClosureDate *time.Time
}
