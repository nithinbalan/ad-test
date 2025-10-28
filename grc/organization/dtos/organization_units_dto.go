package dtos

import commonDTOs "jethur/grc/common/dtos"

type OrganizationUnitsDTO struct {
	ID           uint                   `json:"id"`
	TitleEn      string                 `json:"title_en"`
	TitleAr      string                 `json:"title_ar"`
	Level        uint8                  `json:"level"`
	LevelTitleEn string                 `json:"level_title_en,omitempty"`
	LevelTitleAr string                 `json:"level_title_ar,omitempty"`
	IsPrimary    bool                   `json:"is_primary"`
	Children     []OrganizationUnitsDTO `json:"children,omitempty"`
}

type OrganizationUnitUserPrimaryDetailsDTO struct {
	ID               uint    `json:"id"`
	FirstNameEn      string  `json:"first_name_en"`
	LastNameEn       string  `json:"last_name_en"`
	FirstNameAr      string  `json:"first_name_ar"`
	LastNameAr       string  `json:"last_name_ar"`
	DesignationEn    string  `json:"designation_en"`
	DesignationAr    string  `json:"designation_ar"`
	IsHead           bool    `json:"is_head"`
	Initial          string  `json:"initial"`
	ProfileImagePath *string `json:"profile_image_path"`
}

type OrganizationUnitUserPrimaryDetailsQueryParamsDTO struct {
	Sort     string `form:"sort"`
	SortBy   string `form:"sort_by"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

type OrganizationUnitUserPrimaryDetailsResponseDTO struct {
	Items      []OrganizationUnitUserPrimaryDetailsDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO            `json:"pagination"`
}

type OverallComplianceStatusDTO struct {
	TitleEn       string  `json:"title_en"`
	TitleAr       string  `json:"title_ar"`
	ColorCode     string  `json:"color_code"`
	CompliancePct float64 `json:"compliance_pct"`
}

type ComplianceStatusCountDTO struct {
	ComplianceStatusID uint   `json:"compliance_status_id"`
	TitleEn            string `json:"title_en"`
	TitleAr            string `json:"title_ar"`
	Type               string `json:"type"`
	ColorCode          string `json:"color_code"`
	Count              int64  `json:"count"`
	ScoreFrom          *int   `json:"score_from"`
	ScoreTo            *int   `json:"score_to"`
}

type OrganizationUnitSummaryDTO struct {
	ResponsibleUnitID        uint                        `json:"responsible_unit_id"`
	ComplianceStatus         *OverallComplianceStatusDTO `json:"compliance_status"`
	OpenRiskCount            int64                       `json:"open_risk"`
	OverallControlEfficiency *ControlEfficiencyStatusDTO `json:"overall_control_efficiency"`
}

type DepartmentHeadUpdateDTO struct {
	IsHead bool `json:"is_head"`
}

type ParentOrganizationUnitWithHeadDTO struct {
	OrganizationUnit OrganizationUnitDTO      `json:"organization_unit"`
	HeadUser         *OrganizationUnitUserDTO `json:"head_user,omitempty"`
}

type OrganizationUnitDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
}

type OrganizationUnitUserDTO struct {
	ID               uint    `json:"id"`
	FirstNameEn      string  `json:"first_name_en"`
	LastNameEn       string  `json:"last_name_en"`
	FirstNameAr      string  `json:"first_name_ar"`
	LastNameAr       string  `json:"last_name_ar"`
	ProfileImagePath *string `json:"profile_image_path,omitempty"`
	IsHead           bool    `json:"is_head"`
	Initial          string  `json:"initial"`
}

type OrganizationUnitUserRoleDistributionDTO struct {
	ID    uint   `json:"id"`
	Role  string `json:"Role"`
	Count int64  `json:"count"`
}

type OrganizationUnitOverviewDTO struct {
	OrganizationUnitID        uint   `json:"oraganization_unit_id"`
	OrganizationUnitTitle     string `json:"organization_unit_title"`
	OrganizationUnitHead      string `json:"organization_unit_head"`
	OrganizationUnitHeadEmail string `json:"organization_unit_head_email"`
	TotalEmployees            int64  `json:"total_employees"`
}

type OrganizationUnitKeyMatrixDTO struct {
	TotalRisks       int64 `json:"total_risks"`
	TotalPolicies    int64 `json:"total_policies"`
	TotalProcesses   int64 `json:"total_processes"`
	TotalCompliances int64 `json:"total_compliances"`
}

type OrganizationUnitDistributionItemDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Count     int64  `json:"count"`
}

type OrganizationUnitSummaryReportDTO struct {
	OrganizationUnitOverview  OrganizationUnitOverviewDTO `json:"organization_unit_overview"`
	OrganizationUnitKeyMatrix OrganizationUnitKeyMatrixDTO
}
