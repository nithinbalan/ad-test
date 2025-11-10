package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

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

type OrganizationUnitReportListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
	Filters  map[string][]string
}

type OrganizationUnitReportListDTO struct {
	ID                          uint   `json:"id"`
	Title                       string `json:"title"`
	StartDate                   string `json:"start_date"`
	EndDate                     string `json:"end_date"`
	OrganizationUnitID          uint   `json:"organization_unit_id"`
	OrganizationUnitTitleEn     string `json:"organization_unit_title_en"`
	OrganizationUnitTitleAr     string `json:"organization_unit_title_ar"`
	OrganizationUnitColorCode   string `json:"organization_unit_color_code"`
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

type OrganizationUnitReportListReponseDTO struct {
	Items      []*OrganizationUnitReportListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO     `json:"pagination"`
}

type OrganizationUnitReportFormDTO struct {
	Title              string `json:"title" binding:"required"`
	StartDate          string `json:"start_date" binding:"required"`
	EndDate            string `json:"end_date" binding:"required"`
	OrganizationUnitID uint   `json:"organization_unit_id" binding:"required"`
}

type OrganizationUnitReportParsedDates struct {
	StartDate *time.Time
	EndDate   *time.Time
}

type OrganizationUnitKeyMatrixDTO struct {
	TotalRisks       int64 `json:"total_risks"`
	TotalPolicies    int64 `json:"total_policies"`
	TotalProcesses   int64 `json:"total_processes"`
	TotalCompliances int64 `json:"total_compliances"`
}

type OrganizationUnitUserDetailsDTO struct {
	ID                 uint    `json:"id"`
	FirstNameEn        string  `json:"first_name_en"`
	LastNameEn         string  `json:"last_name_en"`
	FirstNameAr        string  `json:"first_name_ar"`
	LastNameAr         string  `json:"last_name_ar"`
	DesignationTitleEn string  `json:"designation_title_en"`
	DesignationTitleAr string  `json:"designation_title_ar"`
	RoleTitleEn        string  `json:"role_title_en"`
	RoleTitleAr        string  `json:"role_title_ar"`
	IsHead             bool    `json:"is_head"`
	Initial            string  `json:"initial"`
	ProfileImagePath   *string `json:"profile_image_path"`
}

type OrganizationUnitUserRoleDistributionDTO struct {
	ID          uint   `json:"id"`
	RoleTitleEn string `json:"role_title_en"`
	RoleTitleAr string `json:"role_title_ar"`
	Count       int64  `json:"count"`
}

type OrganizationUnitOverviewDTO struct {
	OrganizationUnitID         uint    `json:"oraganization_unit_id"`
	OrganizationUnitTitleEn    string  `json:"organization_unit_title_en"`
	OrganizationUnitTitleAr    string  `json:"organization_unit_title_ar"`
	OrganizationUnitHeadNameEn *string `json:"organization_unit_head_name_en"`
	OrganizationUnitHeadNameAr *string `json:"organization_unit_head_name_ar"`
	OrganizationUnitHeadEmail  *string `json:"organization_unit_head_email"`
	TotalEmployees             int64   `json:"total_employees"`
}

type OrganizationUnitRiskListDTO struct {
	ID                         uint    `json:"id"`
	RiskRegisterID             uint    `json:"risk_register_id"`
	ReferenceNumber            string  `json:"reference_number"`
	RiskCategoryTitleEn        string  `json:"risk_category_title_en"`
	RiskCategoryTitleAr        string  `json:"risk_category_title_ar"`
	RiskStatusTitleEn          string  `json:"risk_status_title_en"`
	RiskStatusTitleAr          string  `json:"risk_status_title_ar"`
	RiskStatusColorCode        string  `json:"risk_status_color_code"`
	Title                      string  `json:"title"`
	Description                *string `json:"description"`
	ReportedAt                 *string `json:"reported_at"`
	CreatedAt                  string  `json:"created_at"`
	UpdatedAt                  *string `json:"updated_at"`
	CurrentLikelihoodTitleEn   *string `json:"current_likelihood_title_en"`
	CurrentLikelihoodTitleAr   *string `json:"current_likelihood_title_ar"`
	CurrentImpactTitleEn       *string `json:"current_impact_title_en"`
	CurrentImpactTitleAr       *string `json:"current_impact_title_ar"`
	CurrentRiskRatingTitleEn   *string `json:"current_risk_rating_title_en"`
	CurrentRiskRatingTitleAr   *string `json:"current_risk_rating_title_ar"`
	CurrentRiskRatingColorCode *string `json:"current_risk_rating_color_code"`
}

type OrganizationUnitPolicyListDTO struct {
	ID                              uint       `json:"id"`
	TitleEn                         string     `json:"title_en"`
	TitleAr                         string     `json:"title_ar"`
	DocumentCategoryTitleEn         string     `json:"document_category_title_en"`
	DocumentCategoryTitleAr         string     `json:"document_category_title_ar"`
	DocumentTypeTitleEn             string     `json:"document_type_title_en"`
	DocumentTypeTitleAr             string     `json:"document_type_title_ar"`
	ResponsibleUserID               uint       `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string     `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string     `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string     `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string     `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string     `json:"responsible_user_initial"`
	ResponsibleUserColor            string     `json:"responsible_user_color"`
	ResponsibleUserEmail            string     `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string     `json:"responsible_user_profile_image_path"`
	DocumentStatusTitleEn           string     `json:"document_status_title_en"`
	DocumentStatusTitleAr           string     `json:"document_status_title_ar"`
	DocumentStatusColorCode         string     `json:"document_status_color_code"`
	DocumentStatusType              string     `json:"document_status_type"`
	DocumentVersion                 float64    `json:"document_version"`
	DocumentURL                     string     `json:"document_url"`
	DocumentName                    string     `json:"document_name"`
	DocumentSize                    int        `json:"document_size"`
	DocumentToken                   string     `json:"document_token"`
	DocumentOriginalName            *string    `json:"document_original_name"`
	DocumentType                    string     `json:"document_type"`
	CreatedAt                       time.Time  `json:"created_at"`
	UpdatedAt                       *time.Time `json:"updated_at"`
	Acknowledged                    bool       `json:"acknowledged"`
}

type OrganizationUnitProcessListDTO struct {
	ID                             uint       `json:"id"`
	ProcessGroupTitleEn            string     `json:"process_group_title_en"`
	ProcessGroupTitleAr            string     `json:"process_group_title_ar"`
	OwnerFirstNameEn               string     `json:"owner_first_name_en"`
	OwnerFirstNameAr               string     `json:"owner_first_name_ar"`
	OwnerLastNameEn                string     `json:"owner_last_name_en"`
	OwnerLastNameAr                string     `json:"owner_last_name_ar"`
	OwnerInitial                   string     `json:"owner_initial"`
	OwnerColor                     string     `json:"owner_color"`
	OwnerEmail                     string     `json:"owner_email"`
	OwnerProfileImagePath          *string    `json:"owner_profile_image_path"`
	ResponsibleUnitTitleEn         string     `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr         string     `json:"responsible_unit_title_ar"`
	ResponsibleUnitColorCode       string     `json:"responsible_unit_color_code"`
	LifecycleStatusTitleEn         string     `json:"lifecycle_status_title_en"`
	LifecycleStatusTitleAr         string     `json:"lifecycle_status_title_ar"`
	LifecycleStatusColorCode       string     `json:"lifecycle_status_color_code"`
	ReferenceNumber                string     `json:"reference_number"`
	TitleEn                        string     `json:"title_en"`
	TitleAr                        string     `json:"title_ar"`
	Description                    *string    `json:"description"`
	MarkedAsAuditableItem          bool       `json:"marked_as_auditable_item"`
	RiskLevelId                    *uint      `json:"risk_level_id"`
	RiskLevelTitleEn               *string    `json:"risk_level_title_en"`
	RiskLevelTitleAr               *string    `json:"risk_level_title_ar"`
	RiskLevelColorCode             *string    `json:"risk_level_color_code"`
	RiskLevelType                  *string    `json:"risk_level_type"`
	ProcessCriticalityLevelID      uint       `json:"process_criticality_level_id"`
	ProcessCriticalityLevelTitleEn string     `json:"process_criticality_level_title_en"`
	ProcessCriticalityLevelTitleAr string     `json:"process_criticality_level_title_ar"`
	CreatedAt                      time.Time  `json:"created_at"`
	UpdatedAt                      *time.Time `json:"updated_at"`
	RiskRatingTitleEn              *string    `json:"risk_rating_title_en"`
	RiskRatingTitleAr              *string    `json:"risk_rating_title_ar"`
	RiskRatingColorCode            *string    `json:"risk_rating_color_code"`
	RiskRatingValue                *uint      `json:"risk_rating_value"`
}

type OrganizationUnitComplianceListDTO struct {
	ID                              uint   `json:"id"`
	ComplianceCategoryTitleEn       string `json:"compliance_category_title_en"`
	ComplianceCategoryTitleAr       string `json:"compliance_category_title_ar"`
	ComplianceStatusTitleEn         string `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr         string `json:"compliance_status_title_ar"`
	ComplianceStatusColorCode       string `json:"compliance_status_color_code"`
	ResponsibleUserFirstNameEn      string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string `json:"responsible_user_initial"`
	ResponsibleUserColor            string `json:"responsible_user_color"`
	ResponsibleUserEmail            string `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string `json:"responsible_user_profile_image_path"`
	SeverityLevelTitleEn            string `json:"severity_level_title_en"`
	SeverityLevelTitleAr            string `json:"severity_level_title_ar"`
	SeverityLevelColorCode          string `json:"severity_level_color_code"`
	ReferenceNumber                 string `json:"reference_number"`
	StandardCode                    string `json:"standard_code"`
	TitleEn                         string `json:"title_en"`
	TitleAr                         string `json:"title_ar"`
	DescriptionEn                   string `json:"description_en"`
	DescriptionAr                   string `json:"description_ar"`
	Version                         string `json:"version"`
	Deadline                        string `json:"deadline"`
	IssueDate                       string `json:"issue_date"`
	ExpiryDate                      string `json:"expiry_date"`
	NextReviewDate                  string `json:"next_review_date"`
	LastReviewedDate                string `json:"last_reviewed_date"`
	CreatedAt                       string `json:"created_at"`
	UpdatedAt                       string `json:"updated_at"`
}

type OrganizationUnitComplianceRequirementListDTO struct {
	ID                                   uint     `json:"id"`
	TitleEn                              string   `json:"title_en"`
	TitleAr                              string   `json:"title_ar"`
	Score                                *float64 `json:"score"`
	Deadline                             *string  `json:"deadline"`
	CreatedAt                            string   `json:"created_at"`
	UpdatedAt                            *string  `json:"updated_at"`
	ComplianceRequirementCategoryTitleEn string   `json:"compliance_requirement_category_title_en"`
	ComplianceRequirementCategoryTitleAr string   `json:"compliance_requirement_category_title_ar"`
	ComplianceStatusTitleEn              *string  `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr              *string  `json:"compliance_status_title_ar"`
	ComplianceStatusColorCode            *string  `json:"compliance_status_color_code"`
	SeverityLevelTitleEn                 string   `json:"severity_level_title_en"`
	SeverityLevelTitleAr                 string   `json:"severity_level_title_ar"`
	SeverityLevelColorCode               string   `json:"severity_level_color_code"`
}

type OrganizationUnitDistributionItemDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Count     int64  `json:"count"`
}

type OrganizationUnitRiskSummaryDTO struct {
	TotalRisks        int64   `json:"total_risks"`
	OpenRisks         int64   `json:"open_risks"`
	ClosedRisks       int64   `json:"closed_risks"`
	ClosurePercentage float64 `json:"closure_percentage"`
}

type OrganizationUnitProcessHealthSummaryDTO struct {
	TotalProcesses   int64   `json:"total_processes"`
	ActiveProcesses  int64   `json:"active_processes"`
	ActivePercentage float64 `json:"active_percentage"`
}

type OrganizationUnitComplianceSummaryDTO struct {
	TotalRequirements  int64                                 `json:"total_requirements"`
	ComplianceStatuses []OrganizationUnitDistributionItemDTO `json:"compliance_statuses"`
}

type OrganizationUnitPolicySummaryDTO struct {
	TotalPolicies             int64   `json:"total_policies"`
	ActivePolicies            int64   `json:"active_policies"`
	PoliciesUnderReview       int64   `json:"policies_under_review"`
	PolicyHealthPercentage    float64 `json:"policy_health_percentage"`
	PolicyAwarenessPercentage float64 `json:"policy_awareness_percentage"`
}

type OrganizationUnitSummaryReportDTO struct {
	ID                               uint                                                      `json:"id"`
	Title                            string                                                    `json:"title"`
	OrganizationUnitID               uint                                                      `json:"organization_unit_id"`
	OrganizationUnitTitleEn          string                                                    `json:"organization_unit_title_en"`
	OrganizationUnitTitleAr          string                                                    `json:"organization_unit_title_ar"`
	OrganizationUnitColorCode        string                                                    `json:"organization_unit_color_code"`
	CreatedAt                        string                                                    `json:"created_at"`
	CreatedUserFirstNameEn           string                                                    `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr           string                                                    `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn            string                                                    `json:"created_user_last_name_en"`
	CreatedUserLastNameAr            string                                                    `json:"created_user_last_name_ar"`
	OrganizationUnitOverview         *OrganizationUnitReportOverviewResponseDTO                `json:"organization_unit_overview"`
	KeyMetrics                       *OrganizationUnitReportKeyMetricsResponseDTO              `json:"key_metrics"`
	UserRoleDistribution             []*OrganizationUnitReportUserRoleResponseDTO              `json:"user_role_distribution"`
	OrganizationUnitUserDetails      []*OrganizationUnitReportUserResponseDTO                  `json:"organization_unit_user_details"`
	RiskSummary                      *OrganizationUnitReportRiskSummaryResponseDTO             `json:"risk_summary"`
	RiskRatingDistribution           []*OrganizationUnitReportRiskRatingResponseDTO            `json:"risk_rating_distribution"`
	RiskMitigationStatusDistribution []*OrganizationUnitReportRiskMitigationResponsDTO         `json:"risk_mitigation_status_distribution"`
	RiskDetails                      []*OrganizationUnitReportRiskResponseDTO                  `json:"risk_details"`
	PolicySummary                    *OrganizationUnitReportPolicySummaryResponseDTO           `json:"policy_summary"`
	PolicyDetails                    []*OrganizationUnitReportPolicyResponseDTO                `json:"policy_details"`
	ProcessHealthSummary             *OrganizationUnitReportProcessSummaryResponseDTO          `json:"process_health_summary"`
	ProcessDetails                   []*OrganizationUnitReportProcessResponseDTO               `json:"process_details"`
	ComplianceSummary                *OrganizationUnitReportComplianceSummaryResponseDTO       `json:"compliance_summary"`
	ComplianceRequirementDetails     []*OrganizationUnitReportComplianceRequirementResponseDTO `json:"compliance_requirement_details"`
	AppUrl                           string                                                    `json:"url"`
}

type OrganizationUnitReportOverviewResponseDTO struct {
	ID                         uint   `json:"id"`
	OrganizationUnitReportID   uint   `json:"organization_unit_report_id"`
	OrganizationUnitTitleEn    string `json:"organization_unit_title_en"`
	OrganizationUnitTitleAr    string `json:"organization_unit_title_ar"`
	OrganizationUnitHeadNameEn string `json:"organization_unit_head_name_en"`
	OrganizationUnitHeadNameAr string `json:"organization_unit_head_name_ar"`
	OrganizationUnitHeadEmail  string `json:"organization_unit_head_email"`
	TotalEmployees             int    `json:"total_employees"`
}

type OrganizationUnitReportKeyMetricsResponseDTO struct {
	ID                       uint `json:"id"`
	OrganizationUnitReportID uint `json:"organization_unit_report_id"`
	TotalUsers               int  `json:"total_users"`
	TotalRisks               int  `json:"total_risks"`
	TotalPolicies            int  `json:"total_policies"`
	TotalProcesses           int  `json:"total_processes"`
	TotalCompliances         int  `json:"total_compliances"`
}

type OrganizationUnitReportUserRoleResponseDTO struct {
	ID                       uint   `json:"id"`
	OrganizationUnitReportID uint   `json:"organization_unit_report_id"`
	RoleTitleEn              string `json:"role_title_en"`
	RoleTitleAr              string `json:"role_title_ar"`
	Count                    int    `json:"count"`
}

type OrganizationUnitReportUserResponseDTO struct {
	ID                       uint   `json:"id"`
	OrganizationUnitReportID uint   `json:"organization_unit_report_id"`
	UserNameEn               string `json:"user_name_en"`
	UserNameAr               string `json:"user_name_ar"`
	DesignationTitleEn       string `json:"designation_title_en"`
	DesignationTitleAr       string `json:"designation_title_ar"`
	RoleTitleEn              string `json:"role_title_en"`
	RoleTitleAr              string `json:"role_title_ar"`
}

type OrganizationUnitReportRiskResponseDTO struct {
	ID                       uint    `json:"id"`
	OrganizationUnitReportID uint    `json:"organization_unit_report_id"`
	RiskTitle                string  `json:"risk_title"`
	Description              *string `json:"description"`
	RiskCategoryTitleEn      string  `json:"risk_category_title_en"`
	RiskCategoryTitleAr      string  `json:"risk_category_title_ar"`
	RiskStatusTitleEn        string  `json:"risk_status_title_en"`
	RiskStatusTitleAr        string  `json:"risk_status_title_ar"`
	RiskStatusColorCode      string  `json:"risk_status_color_code"`
	RiskLikelihoodTitleEn    *string `json:"risk_likelihood_title_en"`
	RiskLikelihoodTitleAr    *string `json:"risk_likelihood_title_ar"`
	RiskImpactTitleEn        *string `json:"risk_impact_title_en"`
	RiskImpactTitleAr        *string `json:"risk_impact_title_ar"`
	RiskRatingTitleEn        *string `json:"risk_rating_title_en"`
	RiskRatingTitleAr        *string `json:"risk_rating_title_ar"`
	RiskRatingColorCode      *string `json:"risk_rating_color_code"`
}

type OrganizationUnitReportPolicyResponseDTO struct {
	ID                       uint   `json:"id"`
	OrganizationUnitReportID uint   `json:"organization_unit_report_id"`
	DocumentTitleEn          string `json:"document_title_en"`
	DocumentTitleAr          string `json:"document_title_ar"`
	DocumentCategoryTitleEn  string `json:"document_category_title_en"`
	DocumentCategoryTitleAr  string `json:"document_category_title_ar"`
	DocumentTypeTitleEn      string `json:"document_type_title_en"`
	DocumentTypeTitleAr      string `json:"document_type_title_ar"`
	ResponsibleUserNameEn    string `json:"responsible_user_name_en"`
	ResponsibleUserNameAr    string `json:"responsible_user_name_ar"`
	Acknowledged             bool   `json:"acknowledged"`
}

type OrganizationUnitReportProcessResponseDTO struct {
	ID                       uint    `json:"id"`
	OrganizationUnitReportID uint    `json:"organization_unit_report_id"`
	ProcessTitleEn           string  `json:"process_title_en"`
	ProcessTitleAr           string  `json:"process_title_ar"`
	ProcessGroupTitleEn      string  `json:"process_group_title_en"`
	ProcessGroupTitleAr      string  `json:"process_group_title_ar"`
	ProcessOwnerNameEn       string  `json:"process_owner_name_en"`
	ProcessOwnerNameAR       string  `json:"process_owner_name_ar"`
	ResponsibleUnitTitleEn   string  `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr   string  `json:"responsible_unit_title_ar"`
	RiskRatingTitleEn        *string `json:"risk_rating_title_en"`
	RiskRatingTitleAr        *string `json:"risk_rating_title_ar"`
	RiskRatingColorCode      *string `json:"risk_rating_color_code"`
	LifecycleStatusTitleEn   string  `json:"lifecycle_status_title_en"`
	LifecycleStatusTitleAr   string  `json:"lifecycle_status_title_ar"`
	LifecycleStatusColorCode string  `json:"lifecycle_status_color_code"`
}

type OrganizationUnitReportComplianceRequirementResponseDTO struct {
	ID                                   uint   `json:"id"`
	OrganizationUnitReportID             uint   `json:"organization_unit_report_id"`
	ComplianceRequirementTitleEn         string `json:"compliance_requirement_title_en"`
	ComplianceRequirementTitleAr         string `json:"compliance_requirement_title_ar"`
	ComplianceRequirementCategoryTitleEn string `json:"compliance_requirement_category_title_en"`
	ComplianceRequirementCategoryTitleAr string `json:"compliance_requirement_category_title_ar"`
	ComplianceStatusTitleEn              string `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr              string `json:"compliance_status_title_ar"`
	ComplianceStatusColorCode            string `json:"compliance_status_color_code"`
	SeverityLevelTitleEn                 string `json:"severity_level_title_en"`
	SeverityLevelTitleAr                 string `json:"severity_level_title_ar"`
	SeverityLevelColorCode               string `json:"severity_level_color_code"`
}

type OrganizationUnitReportRiskSummaryResponseDTO struct {
	ID                       uint    `json:"id"`
	OrganizationUnitReportID uint    `json:"organization_unit_report_id"`
	TotalRiskCount           int     `json:"total_risk_count"`
	OpenRiskCount            int     `json:"open_risk_count"`
	ClosedRiskCount          int     `json:"closed_risk_count"`
	ClosurePercentage        float32 `json:"closure_percentage"`
}

type OrganizationUnitReportRiskRatingResponseDTO struct {
	ID                       uint   `json:"id"`
	OrganizationUnitReportID uint   `json:"organization_unit_report_id"`
	RiskRatingTitleEn        string `json:"risk_rating_title_en"`
	RiskRatingTitleAr        string `json:"risk_rating_title_ar"`
	RiskRatingColorCode      string `json:"risk_rating_color_code"`
	Count                    int    `json:"count"`
}

type OrganizationUnitReportRiskMitigationResponsDTO struct {
	ID                            uint   `json:"id"`
	OrganizationUnitReportID      uint   `json:"organization_unit_report_id"`
	RiskActionPlanStatusTitleEn   string `json:"risk_action_plan_status_title_en"`
	RiskActionPlanStatusTitleAr   string `json:"risk_action_plan_status_title_ar"`
	RiskActionPlanStatusColorCode string `json:"risk_action_plan_status_color_code"`
	Count                         int    `json:"count"`
}

type OrganizationUnitReportPolicySummaryResponseDTO struct {
	ID                       uint    `json:"id"`
	OrganizationUnitReportID uint    `json:"organization_unit_report_id"`
	TotalPolicyCount         int     `json:"total_policy_count"`
	ActivePolicyCount        int     `json:"active_policy_count"`
	UnderReviewPolicyCount   int     `json:"under_review_policy_count"`
	PolicyHealthPercent      float32 `json:"policy_health_percent"`
	PolicyAwarenessPercent   float32 `json:"policy_awareness_percent"`
}

type OrganizationUnitReportProcessSummaryResponseDTO struct {
	ID                       uint    `json:"id"`
	OrganizationUnitReportID uint    `json:"organization_unit_report_id"`
	TotalProcessCount        int     `json:"total_process_count"`
	ActiveProcessCount       int     `json:"active_process_count"`
	ActivePercentage         float32 `json:"active_percentage"`
}

type OrganizationUnitReportComplianceSummaryResponseDTO struct {
	ID                       uint                                   `json:"id"`
	OrganizationUnitReportID uint                                   `json:"organization_unit_report_id"`
	TotalRequirementCount    int                                    `json:"total_requirement_count"`
	ComplianceStatuses       *string                                `json:"compliance_statuses"`
	ComplianceStatusesParsed []*OrganizationUnitDistributionItemDTO `json:"compliance_statuses_parsed" gorm:"-"`
}
