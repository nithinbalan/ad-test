package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"

	"time"
)

type RiskRegisterListDTO struct {
	ID                         uint       `json:"id"`
	ReferenceNumber            string     `json:"reference_number"`
	RiskCategoryTitleEn        string     `json:"risk_category_title_en"`
	RiskCategoryTitleAr        string     `json:"risk_category_title_ar"`
	RiskSourceTitleEn          *string    `json:"risk_source_title_en"`
	RiskSourceTitleAr          *string    `json:"risk_source_title_ar"`
	ResponsibleUnitID          uint       `json:"responsible_unit_id"`
	ResponsibleUnitTitleEn     string     `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr     string     `json:"responsible_unit_title_ar"`
	ResponsibleUserFirstNameEn string     `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr string     `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn  string     `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr  string     `json:"responsible_user_last_name_ar"`
	RiskStatusTitleEn          string     `json:"risk_status_title_en"`
	RiskStatusTitleAr          string     `json:"risk_status_title_ar"`
	RiskStatusColorCode        string     `json:"risk_status_color_code"`
	Title                      string     `json:"title"`
	Description                *string    `json:"description"`
	ReportedAt                 *time.Time `json:"reported_at"`
	ReviewFrequencyTitleEn     *string    `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr     *string    `json:"review_frequency_title_ar"`
	NextReviewDate             *string    `json:"next_review_date"`
	CurrentRiskRatingValueID   *uint      `json:"current_risk_rating_value_id"`
	Comments                   *string    `json:"comments"`
	CreatedAt                  string     `json:"created_at"`
	UpdatedAt                  *string    `json:"updated_at"`
	CreatedUserFirstNameEn     string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr     string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn      string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr      string     `json:"created_user_last_name_ar"`
	UpdatedUserFirstNameEn     *string    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr     *string    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn      *string    `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr      *string    `json:"updated_user_last_name_ar"`
}

type RiskRegisterListQueryParamsDTO struct {
	Sort               string
	SortBy             string
	Keyword            string
	Status             string
	Page               int
	PageSize           int
	Filters            map[string][]string
	QuickFilter        []string
	ReportedAtFrom     string
	ReportedAtTo       string
	NextReviewDateFrom string
	NextReviewDateTo   string
}

type RiskRegisterListResponseDTO struct {
	Items      []*RiskRegisterListDTO       `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type RiskDocumentDTO struct {
	ID           *uint   `json:"id"`
	URL          *string `json:"url"`
	Name         string  `json:"name"`
	OriginalName string  `json:"original_name"`
	Token        string  `json:"token"`
	Type         string  `json:"type"`
	Size         int     `json:"size"`
}

type RiskRegisterDetailDTO struct {
	ID                        int                              `json:"id"`
	ReferenceNumber           string                           `json:"reference_number"`
	RiskCategory              RiskCategoryDTO                  `json:"risk_category"`
	RiskSource                *RiskSourceDTO                   `json:"risk_source"`
	ResponsibleUnit           settingsDTOs.OrganizationUnitDTO `json:"responsible_unit"`
	ResponsibleUser           organizationDTOs.UserMinimalDTO  `json:"responsible_user"`
	RiskStatus                RiskStatusDTO                    `json:"risk_status"`
	Title                     string                           `json:"title"`
	Description               *string                          `json:"description"`
	ReportedAt                *time.Time                       `json:"reported_at"`
	RiskMatrixLikelihood      *RiskMatrixLikelihoodDTO         `json:"current_risk_matrix_likelihood"`
	RiskMatrixImpact          *RiskMatrixImpactDTO             `json:"current_risk_matrix_impact"`
	RiskMatrixRiskRating      *RiskMatrixRiskRatingDTO         `json:"current_risk_rating"`
	RiskMatrixRiskRatingValue *RiskMatrixRiskRatingValueDTO    `json:"current_risk_rating_value_id"`
	ReviewFrequency           *masterDTOs.ReviewFrequencyDTO   `json:"review_frequency"`
	NextReviewDate            *time.Time                       `json:"next_review_date"`
	Comments                  *string                          `json:"comments"`
	Documents                 []*RiskDocumentDTO               `json:"documents"`
	CreatedAt                 time.Time                        `json:"created_at"`
	UpdatedAt                 *time.Time                       `json:"updated_at"`
	CreatedByUser             organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser             organizationDTOs.UserMinimalDTO  `json:"updated_by"`
	Controls                  []*RiskControlDTO                `json:"controls"`
	ComplianceRequirements    []*RiskComplianceRequirementDTO  `json:"compliance_requirements"`
	Processes                 []*RiskProcessDTO                `json:"processes"`
	StrategicObjectives       []*RiskStrategicObjectiveDTO     `json:"strategic_objectives"`
	SLARegisters              []*RiskSLARegisterDTO            `json:"sla_registers"`
	Kris                      []*RiskKRIDTO                    `json:"kris"`
}

type RiskControlDTO struct {
	ID             *uint                             `json:"id"`
	RiskRegisterID uint                              `json:"risk_register_id"`
	ControlID      uint                              `json:"control_id"`
	Control        organizationDTOs.ControlDetailDTO `json:"control"`
}

type RiskComplianceRequirementDTO struct {
	ID                      *uint                                         `json:"id"`
	RiskRegisterID          uint                                          `json:"risk_register_id"`
	ComplianceRequirementID uint                                          `json:"compliance_requirement_id"`
	ComplianceRequirement   complianceDTOs.ComplianceRequirementDetailDTO `json:"compliance_requirement"`
}

type RiskProcessDTO struct {
	ID             *uint                             `json:"id"`
	RiskRegisterID uint                              `json:"risk_register_id"`
	ProcessID      uint                              `json:"process_id"`
	Process        organizationDTOs.ProcessDetailDTO `json:"process"`
}

type RiskStrategicObjectiveDTO struct {
	ID                   *uint                                  `json:"id"`
	RiskRegisterID       uint                                   `json:"risk_register_id"`
	StrategicObjectiveID uint                                   `json:"strategic_objective_id"`
	StrategicObjective   organizationDTOs.StrategicObjectiveDTO `json:"strategic_objective"`
}

type RiskKRIDTO struct {
	ID             *uint `json:"id"`
	RiskRegisterID uint  `json:"risk_register_id"`
	KriID          uint  `json:"kri_id"`
}

type RiskSLARegisterDTO struct {
	ID             *uint                               `json:"id"`
	RiskRegisterID uint                                `json:"risk_register_id"`
	SLARegisterID  uint                                `json:"sla_register_id"`
	SLARegister    complianceDTOs.SLARegisterDetailDTO `json:"sla_register"`
}
type RiskRegisterFormDTO struct {
	RiskCategoryID           uint                          `json:"risk_category_id" binding:"required"`
	RiskSourceID             *uint                         `json:"risk_source_id"`
	ResponsibleUserID        uint                          `json:"responsible_user_id" binding:"required"`
	Title                    string                        `json:"title" binding:"required,max=255"`
	Description              *string                       `json:"description"`
	ReportedAt               *string                       `json:"reported_at" binding:"required,datetime=2006-01-02"`
	ReviewFrequencyID        *uint                         `json:"review_frequency_id"`
	NextReviewDate           *string                       `json:"next_review_date" binding:"omitempty,datetime=2006-01-02"`
	Comments                 *string                       `json:"comments"`
	Documents                []*commonDTOs.DocumentFormDTO `json:"documents"`
	ComplianceRequirementIDs []*uint                       `json:"compliance_requirement_ids"`
	ControlIDs               []*uint                       `json:"control_ids"`
	ProcessIDs               []*uint                       `json:"process_ids"`
	StrategicObjectiveIDs    []*uint                       `json:"strategic_objective_ids"`
	SLARegisterIDs           []*uint                       `json:"sla_register_ids"`
}

type RiskRegisterParsedDates struct {
	ReportedAt     *time.Time
	NextReviewDate *time.Time
}

type RiskResponseStrategyFormDTO struct {
	RiskResponseStrategyID uint `json:"risk_response_strategy_id" binding:"required"`
}
