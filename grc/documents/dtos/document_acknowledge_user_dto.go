package dtos

import commonDTOs "jethur/grc/common/dtos"

type DocumentAcknowledgeListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
}

type DocumentAcknowledgementListDTO struct {
	ID                      uint   `json:"user_id"`
	FirstNameEn             string `json:"first_name_en"`
	FirstNameAr             string `json:"first_name_ar"`
	LastNameEn              string `json:"last_name_en"`
	LastNameAr              string `json:"last_name_ar"`
	Email                   string `json:"email"`
	OrganizationUnitTitleEn string `json:"organization_unit_title_en"`
	OrganizationUnitTitleAr string `json:"organization_unit_title_ar"`
	Acknowledged            bool   `json:"acknowledged"`
}

type DocumentAcknowledgementListResponseDTO struct {
	Items      []DocumentAcknowledgementListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO     `json:"pagination"`
}

type DocumentAcknowledgmentSummaryDTO struct {
	TotalUsers             int64   `json:"total_users"`
	AcknowledgedUsers      int64   `json:"acknowledged_users"`
	NotAcknowledgedUsers   int64   `json:"not_acknowledged_users"`
	AcknowledgedPercent    float64 `json:"acknowledged_percent"`
	NotAcknowledgedPercent float64 `json:"not_acknowledged_percent"`
}
