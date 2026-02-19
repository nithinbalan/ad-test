package dtos

import commonDTOs "jethur/grc/common/dtos"

type ComplianceAssessmentAcknowledgeListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
}

type ComplianceAssessmentAcknowledgementListDTO struct {
	ID                      uint   `json:"user_id"`
	FirstNameEn             string `json:"first_name_en"`
	FirstNameAr             string `json:"first_name_ar"`
	LastNameEn              string `json:"last_name_en"`
	LastNameAr              string `json:"last_name_ar"`
	Email                   string `json:"email"`
	Initial                 string `json:"initial"`
	Color                   string `json:"color"`
	Comment                 string `json:"comment"`
	ProfileImagePath        string `json:"profile_image_path"`
	OrganizationUnitID      uint   `json:"organization_unit_id"`
	OrganizationUnitTitleEn string `json:"organization_unit_title_en"`
	OrganizationUnitTitleAr string `json:"organization_unit_title_ar"`
	Acknowledged            bool   `json:"acknowledged"`
}

type ComplianceAssessmentAcknowledgementListResponseDTO struct {
	Items      []ComplianceAssessmentAcknowledgementListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO                 `json:"pagination"`
}

type ComplianceAssessmentAcknowledgmentSummaryDTO struct {
	TotalDepartments        int64 `json:"total_departments"`
	AcknowledgedDepartments int64 `json:"acknowledged_departments"`
	PendingDepartments      int64 `json:"pending_departments"`
}

type ComplianceAssessmentAcknowledgeFormDTO struct {
	Comment *string `json:"comment" binding:"omitempty"`
}

type ComplianceAssessmentAcknowledgmentOrgUnitSummaryDTO struct {
	OrganizationUnitID      uint    `json:"organization_unit_id"`
	OrganizationUnitTitleEn string  `json:"organization_unit_title_en"`
	OrganizationUnitTitleAr string  `json:"organization_unit_title_ar"`
	TotalUsers              int64   `json:"total_users"`
	AcknowledgedUsers       int64   `json:"acknowledged_users"`
	NotAcknowledgedUsers    int64   `json:"not_acknowledged_users"`
	AcknowledgedPercent     float64 `json:"acknowledged_percent"`
}
