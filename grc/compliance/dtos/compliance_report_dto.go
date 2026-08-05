package dtos

import (
	"html/template"
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"time"
)

type ComplianceReportListDTO struct {
	ID                          uint   `json:"id"`
	Title                       string `json:"title"`
	AssessmentType              string `json:"assessment_type"`
	StartDate                   string `json:"start_date"`
	EndDate                     string `json:"end_date"`
	CreatedAt                   string `json:"created_at"`
	IsCyber                     bool   `json:"is_cyber"`
	ComplianceRegisterTitleEn   string `json:"compliance_register_title_en"`
	ComplianceRegisterTitleAr   string `json:"compliance_register_title_ar"`
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

type ComplianceReportListResponseDTO struct {
	Items      []*ComplianceReportListDTO   `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type ComplianceReportDetailsDTO struct {
	ID                       uint                                  `json:"id"`
	Title                    string                                `json:"title"`
	StartDate                string                                `json:"start_date"`
	EndDate                  string                                `json:"end_date"`
	AssessmentType           string                                `json:"assessment_type"`
	ComplianceRegisterID     *uint                                 `json:"compliance_register_id"`
	ResponsibleUnitID        *uint                                 `json:"responsible_unit_id"`
	CreatedAt                string                                `json:"created_at"`
	UpdatedAt                *string                               `json:"updated_at"`
	CreatedByUser            organizationDTOs.UserMinimalDTO       `json:"created_by"`
	ComplianceRegister       *ComplianceRegister                   `json:"compliance_register"`
	ResponsibleUnit          *organizationDTOs.OrganizationUnitDTO `json:"responsible_unit"`
	ComplianceReportSections []*ComplianceReportSection            `json:"compliance_report_sections"`
}

type ComplianceRegister struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
}

type ComplianceReportSection struct {
	ID                          uint                                 `json:"id"`
	TitleEn                     string                               `json:"title_en"`
	TitleAr                     string                               `json:"title_ar"`
	DescriptionEn               template.HTML                        `json:"description_en"`
	DescriptionAr               template.HTML                        `json:"description_ar"`
	Type                        string                               `json:"type"`
	Order                       int                                  `json:"order"`
	ComplianceReportDepartments []*ComplianceReportDepartmentSection `json:"compliance_report_departments"`
	ComplianceReportStandards   []*ComplianceReportClauseSection     `json:"compliance_report_standards"`
	ComplianceReportActionPlans []*ComplianceReportActionPlan        `json:"compliance_report_action_plans"`
}

type ComplianceReportDepartmentSection struct {
	ComplianceRegister         *ComplianceRegister           `json:"compliance_register"`
	ComplianceReportDepartment []*ComplianceReportDepartment `json:"compliance_report_department"`
}

type ComplianceReportDepartment struct {
	OrganizationUnit                   settingsDTOs.OrganizationUnitDTO `json:"organization_unit"`
	ComplianceRegister                 *ComplianceRegister              `json:"compliance_register"`
	ComplianceStatus                   *ComplianceStatusDTO             `json:"compliance_status"`
	CompliancePercentage               *float64                         `json:"compliance_percentage"`
	IsRiskBasedAssessment              *bool                            `json:"is_risk_based_assessment"`
	IsRegulationNotificationAssessment *bool                            `json:"is_regulation_notification_assessment"`
	NonCompliantAreas                  []*ComplianceReportRequirement   `json:"non_compliant_areas"`
}

// type ComplianceReportStandard struct {
// 	ComplianceReportClauses []*ComplianceReportClause `json:"compliance_report_clauses"`
// }

type ComplianceReportRequirement struct {
	ID                                 uint                             `json:"id"`
	TitleEn                            string                           `json:"title_en"`
	TitleAr                            string                           `json:"title_ar"`
	ComplianceStatus                   *ComplianceStatusDTO             `json:"compliance_status"`
	AverageComplianceStatus            *ComplianceStatusDTO             `json:"average_compliance_status"`
	Score                              *int                             `json:"score"`
	AverageScore                       *float64                         `json:"average_score"`
	InherentScore                      *float64                         `json:"inherent_score"`
	InherentRiskRatingID               *uint                            `json:"inherent_risk_rating_id"`
	InherentRiskRating                 *RiskMatrixRiskRatingDTO         `json:"inherent_risk_rating"`
	ResidualScore                      *float64                         `json:"residual_score"`
	ResidualRiskRatingID               *uint                            `json:"residual_risk_rating_id"`
	ResidualRiskRating                 *RiskMatrixRiskRatingDTO         `json:"residual_risk_rating"`
	IsAcknowledged                     bool                             `json:"is_acknowledged"`
	AcknowledgeUser                    *organizationDTOs.UserMinimalDTO `json:"acknowledge_user"`
	AcknowledgeComment                 *string                          `json:"acknowledge_comment"`
	IsRiskBasedAssessment              *bool                            `json:"is_risk_based_assessment"`
	IsRegulationNotificationAssessment *bool                            `json:"is_regulation_notification_assessment"`
	ActionPlans                        []*ComplianceReportActionPlan    `json:"action_plans"`
}

type RiskMatrixRiskRatingDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Type      string `json:"type"`
}

type ComplianceReportActionPlan struct {
	OrganizationUnit settingsDTOs.OrganizationUnitDTO `json:"organization_unit"`
	Title            string                           `json:"title"`
	Description      *string                          `json:"description"`
	ActionPlanStatus masterDTOs.ActionPlanStatusDTO   `json:"action_plan_status"`
	TargetDate       string                           `json:"target_date"`
	SeverityLevel    *masterDTOs.SeverityLevelDTO     `json:"severity_level"`
}

type ComplianceReportClauseSection struct {
	ComplianceRegister     *ComplianceRegister       `json:"compliance_register"`
	ComplianceReportClause []*ComplianceReportClause `json:"compliance_report_clause"`
}

type ComplianceReportClause struct {
	ClauseNo               string                         `json:"clause_no"`
	ClauseTitleEn          string                         `json:"clause_title_en"`
	ClauseTitleAr          *string                        `json:"clause_title_ar"`
	DescriptionEn          *string                        `json:"description_en"`
	DescriptionAr          *string                        `json:"description_ar"`
	ComplianceRequirements []*ComplianceReportRequirement `json:"compliance_requirements"`
}

type ComplianceReportFormDTO struct {
	Title                string `json:"title" binding:"required"`
	StartDate            string `json:"start_date" binding:"required"`
	EndDate              string `json:"end_date" binding:"required"`
	AssessmentType       string `json:"assessment_type" binding:"required"`
	ComplianceRegisterID *uint  `json:"compliance_register_id"`
	ResponsibleUnitID    *uint  `json:"responsible_unit_id"`
}

type ComplianceReportParsedDates struct {
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
}

type ComplianceReportListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
	Filters  map[string][]string
	IsCyber  bool
}

type ComplianceReportSectionFormDTO struct {
	TitleEn       string  `json:"title_en" binding:"required"`
	TitleAr       string  `json:"title_ar" binding:"required"`
	DescriptionEn *string `json:"description_en"`
	DescriptionAr *string `json:"description_ar"`
}
