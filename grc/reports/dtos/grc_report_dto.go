package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
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
	ID                     uint                                  `json:"id"`
	Title                  string                                `json:"title"`
	StartDate              string                                `json:"start_date"`
	EndDate                string                                `json:"end_date"`
	CreatedUserFirstNameEn string                                `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr string                                `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn  string                                `json:"created_user_last_name_en"`
	CreatedUserLastNameAr  string                                `json:"created_user_last_name_ar"`
	CreatedAt              string                                `json:"created_at"`
	Sections               []GRCReportSectionResponseDTO         `json:"sections"`
	GovernanceSummary      GRCGovernanceSummaryResponseDTO       `json:"governance_summary"`
	Documents              []GRCGovernanceDocumentResponseDTO    `json:"documents"`
	Tasks                  []GRCGovernanceTaskResponseDTO        `json:"tasks"`
	RiskSummary            GRCRiskSummaryResponseDTO             `json:"risk_summary"`
	TopFiveRisks           []GRCReportRiskTopFiveRiskResponseDTO `json:"top_five_risks"`
	RiskKRIs               []GRCReportRiskKriResponseDTO         `json:"risk_kris"`
	RiskControls           []GRCReportRiskControlResponseDTO     `json:"risk_controls"`
	ComplianceSummary      []*GRCReportComplianceSummaryDTO      `json:"compliance_summary"`
	SLA                    []*GRCReportSLAResponseDTO            `json:"sla"`
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
	ID                    uint     `json:"id"`
	ReferenceNumber       string   `json:"reference_number"`
	PoliciesInSystem      string   `json:"policies_in_system"`
	TotalPolicies         int      `json:"total_policies"`
	RolesDefined          int      `json:"roles_defined"`
	PolicyApprovedRate    *float64 `json:"policy_approved_rate"`
	Highlights            *string  `json:"highlights"`
	Challenges            *string  `json:"challenges"`
	DocumentOverviewChart *string  `json:"document_overview_chart"`
	DocumentTypeChart     *string  `json:"document_type_chart"`
	TopPublishedDocuments *string  `json:"top_published_documents"`
	CreatedByID           uint     `json:"created_by_id"`
	CreatedAt             string   `json:"created_at"`
}

type GRCGovernanceDocumentResponseDTO struct {
	ID                      uint    `json:"id"`
	ReferenceNumber         string  `json:"reference_number"`
	DocumentTitleEn         string  `json:"document_title_en"`
	DocumentTitleAr         string  `json:"document_title_ar"`
	ResponsibleUserNameEn   string  `json:"responsible_user_name_en"`
	ResponsibleUserNameAr   string  `json:"responsible_user_name_ar"`
	ResponsibleUnitTitleEn  string  `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr  string  `json:"responsible_unit_title_ar"`
	TotalEmployees          int     `json:"total_employees"`
	AcknowledgedCount       int     `json:"acknowledged_count"`
	PendingCount            int     `json:"pending_count"`
	DocumentStatusTitleEn   string  `json:"document_status_title_en"`
	DocumentStatusTitleAr   string  `json:"document_status_title_ar"`
	DocumentStatusColorCode string  `json:"document_status_color_code"`
	AcknowledgeDueDate      *string `json:"acknowledge_due_date"`
	CreatedByID             uint    `json:"created_by_id"`
	CreatedAt               string  `json:"created_at"`
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

type TopFiveRiskDTO struct {
	ID              uint    `json:"id"`
	ReferenceNumber string  `json:"reference_number"`
	Title           string  `json:"title"`
	Description     *string `json:"description"`

	RiskCategoryTitleEn *string `json:"risk_category_title_en"`
	RiskCategoryTitleAr *string `json:"risk_category_title_ar"`

	ResponsibleUserFirstNameEn string `json:"-"`
	ResponsibleUserLastNameEn  string `json:"-"`
	ResponsibleUserFirstNameAr string `json:"-"`
	ResponsibleUserLastNameAr  string `json:"-"`

	ResponsibleUserNameEn *string `json:"responsible_user_name_en"`
	ResponsibleUserNameAr *string `json:"responsible_user_name_ar"`

	ResponsibleUnitTitleEn *string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr *string `json:"responsible_unit_title_ar"`

	RiskStatusType *string `json:"risk_status_type"`
	IsOpen         bool    `json:"is_open"`

	InherentRiskRatingTitleEn   *string `json:"inherent_risk_rating_title_en"`
	InherentRiskRatingTitleAr   *string `json:"inherent_risk_rating_title_ar"`
	InherentRiskRatingColorCode *string `json:"inherent_risk_rating_color_code"`
	InherentRiskRatingScore     *uint   `json:"inherent_risk_rating_score"`

	InherentLikelihoodTitleEn *string `json:"inherent_likelihood_title_en"`
	InherentLikelihoodTitleAr *string `json:"inherent_likelihood_title_ar"`
	InherentImpactTitleEn     *string `json:"inherent_impact_title_en"`
	InherentImpactTitleAr     *string `json:"inherent_impact_title_ar"`

	ResidualRiskRatingTitleEn   *string `json:"residual_risk_rating_title_en"`
	ResidualRiskRatingTitleAr   *string `json:"residual_risk_rating_title_ar"`
	ResidualRiskRatingColorCode *string `json:"residual_risk_rating_color_code"`
	ResidualRiskRatingScore     *uint   `json:"residual_risk_rating_score"`

	ResidualLikelihoodTitleEn *string `json:"residual_likelihood_title_en"`
	ResidualLikelihoodTitleAr *string `json:"residual_likelihood_title_ar"`
	ResidualImpactTitleEn     *string `json:"residual_impact_title_en"`
	ResidualImpactTitleAr     *string `json:"residual_impact_title_ar"`

	RiskResponseStrategyTitleEn *string `json:"risk_response_strategy_title_en"`
	RiskResponseStrategyTitleAr *string `json:"risk_response_strategy_title_ar"`

	ExistingControls *string `json:"existing_controls"`
}

type GRCReportRiskKRIDTO struct {
	RiskRegisterTitle        string  `json:"risk_register_title"`
	KRIRegisterTitle         string  `json:"kri_register_title"`
	KRIRegisterDescription   *string `json:"kri_register_description"`
	LowThreshold             *int    `json:"low_threshold"`
	MediumThreshold          *int    `json:"medium_threshold"`
	CurrentValue             *int    `json:"current_value"`
	RiskStatusTitleEn        string  `json:"risk_status_title_en"`
	RiskStatusTitleAr        string  `json:"risk_status_title_ar"`
	RiskStatusTitleColorCode string  `json:"risk_status_title_color_code"`
}

type RiskControlReportDTO struct {
	ControlCode                          string     `json:"control_code"`
	ControlCategoryTitleEn               string     `json:"control_category_title_en"`
	ControlCategoryTitleAr               string     `json:"control_category_title_ar"`
	Title                                string     `json:"title"`
	Description                          *string    `json:"description"`
	ResponsibleUserNameEn                string     `json:"responsible_user_name_en"`
	ResponsibleUserNameAr                string     `json:"responsible_user_name_ar"`
	ResponsibleUnitTitleEn               string     `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr               string     `json:"responsible_unit_title_ar"`
	ControlImplementationStatusTitleEn   string     `json:"control_implementation_status_title_en"`
	ControlImplementationStatusTitleAr   string     `json:"control_implementation_status_title_ar"`
	ControlImplementationStatusColorCode string     `json:"control_implementation_status_color_code"`
	ControlEfficiencyStatusTitleEn       string     `json:"control_efficiency_status_title_en"`
	ControlEfficiencyStatusTitleAr       string     `json:"control_efficiency_status_title_ar"`
	ControlEfficiencyStatusColorCode     string     `json:"control_efficiency_status_color_code"`
	TestFrequencyTitleEn                 *string    `json:"test_frequency_title_en"`
	TestFrequencyTitleAr                 *string    `json:"test_frequency_title_ar"`
	NextTestDate                         *time.Time `json:"next_test_date"`
	TestResponsibleUserNameEn            *string    `json:"test_responsible_user_name_en"`
	TestResponsibleUserNameAr            *string    `json:"test_responsible_user_name_ar"`
}

type GRCRiskSummaryResponseDTO struct {
	RiskIdentifiedCount int     `json:"risk_identified_count"`
	CriticalRisks       *string `json:"critical_risks"`
	TreatmentStrategies *string `json:"treatment_strategies"`
	RiskScoreSummary    *string `json:"risk_score_summary"`
	Kris                *string `json:"kris"`
	RiskOverviewChart   *string `json:"risk_overview_chart"`
	RiskHeatmapChart    *string `json:"risk_heatmap_chart"`
	TopRisks            *string `json:"top_risks"`
	RiskTrend           *string `json:"risk_trend"`
}

type GRCReportRiskKriResponseDTO struct {
	ID                       uint    `json:"id"`
	GRCReportID              uint    `json:"grc_report_id"`
	RiskRegisterTitle        string  `json:"risk_register_title"`
	KRIRegisterTitle         string  `json:"kri_register_title"`
	KRIRegisterDescription   *string `json:"kri_register_description"`
	LowThreshold             *int    `json:"low_threshold"`
	MediumThreshold          *int    `json:"medium_threshold"`
	CurrentValue             *int    `json:"current_value"`
	RiskStatusTitleEn        string  `json:"risk_status_title_en"`
	RiskStatusTitleAr        string  `json:"risk_status_title_ar"`
	RiskStatusTitleColorCode string  `json:"risk_status_color_code"`
}

type GRCReportRiskControlResponseDTO struct {
	ID                                   uint       `json:"id"`
	GRCReportID                          uint       `json:"grc_report_id"`
	ControlCode                          string     `json:"control_code"`
	ControlCategoryTitleEn               string     `json:"control_category_title_en"`
	ControlCategoryTitleAr               string     `json:"control_category_title_ar"`
	Title                                string     `json:"title"`
	Description                          *string    `json:"description"`
	ResponsibleUserNameEn                string     `json:"responsible_user_name_en"`
	ResponsibleUserNameAr                string     `json:"responsible_user_name_ar"`
	ResponsibleUnitTitleEn               string     `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr               string     `json:"responsible_unit_title_ar"`
	ControlImplementationStatusTitleEn   string     `json:"control_implementation_status_title_en"`
	ControlImplementationStatusTitleAr   string     `json:"control_implementation_status_title_ar"`
	ControlImplementationStatusColorCode string     `json:"control_implementation_status_color_code"`
	ControlEfficiencyStatusTitleEn       string     `json:"control_efficiency_status_title_en"`
	ControlEfficiencyStatusTitleAr       string     `json:"control_efficiency_status_title_ar"`
	ControlEfficiencyStatusColorCode     string     `json:"control_efficiency_status_color_code"`
	TestFrequencyTitleEn                 *string    `json:"test_frequency_title_en"`
	TestFrequencyTitleAr                 *string    `json:"test_frequency_title_ar"`
	NextTestDate                         *time.Time `json:"next_test_date"`
	TestResponsibleUserNameEn            *string    `json:"test_responsible_user_name_en"`
	TestResponsibleUserNameAr            *string    `json:"test_responsible_user_name_ar"`
}

type GRCReportRiskTopFiveRiskResponseDTO struct {
	ID                          uint    `json:"id"`
	GRCReportID                 uint    `json:"grc_report_id"`
	RiskReferenceNumber         string  `json:"risk_reference_number"`
	RiskCategoryTitleEn         string  `json:"risk_category_title_en"`
	RiskCategoryTitleAr         string  `json:"risk_category_title_ar"`
	Title                       string  `json:"title"`
	Description                 *string `json:"description"`
	ResponsibleUserNameEn       string  `json:"responsible_user_name_en"`
	ResponsibleUserNameAr       string  `json:"responsible_user_name_ar"`
	ResponsibleUnitTitleEn      string  `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr      string  `json:"responsible_unit_title_ar"`
	IsOpen                      bool    `json:"is_open"`
	InherentRiskRatingTitleEn   *string `json:"inherent_risk_rating_title_en"`
	InherentRiskRatingTitleAr   *string `json:"inherent_risk_rating_title_ar"`
	InherentRiskRatingColorCode *string `json:"inherent_risk_rating_color_code"`
	InherentRiskRatingScore     *uint   `json:"inherent_risk_rating_score"`
	InherentLikelihoodTitleEn   *string `json:"inherent_likelihood_title_en"`
	InherentLikelihoodTitleAr   *string `json:"inherent_likelihood_title_ar"`
	InherentImpactTitleEn       *string `json:"inherent_impact_title_en"`
	InherentImpactTitleAr       *string `json:"inherent_impact_title_ar"`
	ResidualRiskRatingTitleEn   *string `json:"residual_risk_rating_title_en"`
	ResidualRiskRatingTitleAr   *string `json:"residual_risk_rating_title_ar"`
	ResidualRiskRatingColorCode *string `json:"residual_risk_rating_color_code"`
	ResidualRiskRatingScore     *uint   `json:"residual_risk_rating_score"`
	ResidualLikelihoodTitleEn   *string `json:"residual_likelihood_title_en"`
	ResidualLikelihoodTitleAr   *string `json:"residual_likelihood_title_ar"`
	ResidualImpactTitleEn       *string `json:"residual_impact_title_en"`
	ResidualImpactTitleAr       *string `json:"residual_impact_title_ar"`
	ExistingControls            *string `json:"existing_controls"`
	RiskResponseStrategyTitleEn *string `json:"risk_response_strategy_title_en"`
	RiskResponseStrategyTitleAr *string `json:"risk_response_strategy_title_ar"`
	ApprovalNote                *string `json:"approval_note"`
}

//compliance and sla

type GRCReportComplianceSummaryDTO struct {
	ID                     uint                                `json:"id"`
	GRCReportID            uint                                `json:"grc_report_id"`
	OrganizationUnit       *settingsDTOs.OrganizationUnitDTO   `json:"organization_unit"`
	OverallComplianceScore float64                             `json:"overall_compliance_score"`
	ReviewUser             *organizationDTOs.UserMinimalDTO    `json:"review_user"`
	ApproveUser            *organizationDTOs.UserMinimalDTO    `json:"approve_user"`
	Departments            []*GRCReportComplianceDepartmentDTO `json:"departments"`
}

type GRCReportComplianceDepartmentDTO struct {
	ID               uint                                `json:"id"`
	OrganizationUnit *settingsDTOs.OrganizationUnitDTO   `json:"organization_unit"`
	ComplianceStatus *complianceDTOs.ComplianceStatusDTO `json:"compliance_status"`
	ComplianceScore  float64                             `json:"compliance_score"`
}

type GRCReportSLAResponseDTO struct {
	ID                            uint    `json:"id"`
	GRCReportID                   uint    `json:"grc_report_id"`
	TitleEn                       string  `json:"title_en"`
	TitleAr                       *string `json:"title_ar"`
	ParentOrganizationUnitTitleEn string  `json:"parent_organization_unit_title_en"`
	ParentOrganizationUnitTitleAr string  `json:"parent_organization_unit_title_ar"`
	OrganizationUnitTitleEn       string  `json:"organization_unit_title_en"`
	OrganizationUnitTitleAr       string  `json:"organization_unit_title_ar"`
	ComplianceStatusTitleEn       *string `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr       *string `json:"compliance_status_title_ar"`
	ComplianceStatusColorCode     string  `json:"compliance_status_color_code"`
	LifeCycleStatusTitleEn        string  `json:"life_cycle_status_title_en"`
	LifeCycleStatusTitleAr        string  `json:"life_cycle_status_title_ar"`
	LifeCycleStatusColorCode      string  `json:"life_cycle_status_color_code"`
	ExperyDate                    *string `json:"expery_date"`
	IssueDate                     *string `json:"issue_date"`
}

type RiskCountDTO struct {
	Year        int    `json:"year"`
	Month       int    `json:"month"`
	MonthName   string `json:"month_name"`
	NewRisks    int    `json:"new_risks"`
	ClosedRisks int    `json:"closed_risks"`
}

type GRCReportHighlightChallengeFormDTO struct {
	Description *string `json:"description" binding:"required"`
}
