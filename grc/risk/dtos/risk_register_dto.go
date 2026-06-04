package dtos

import (
	"time"

	assetDTOs "jethur/grc/asset/dtos"
	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
)

type RiskRegisterListDTO struct {
	ID                              uint    `json:"id"`
	ReferenceNumber                 string  `json:"reference_number"`
	RiskCategoryTitleEn             *string `json:"risk_category_title_en"`
	RiskCategoryTitleAr             *string `json:"risk_category_title_ar"`
	RiskSourceTitleEn               *string `json:"risk_source_title_en"`
	RiskSourceTitleAr               *string `json:"risk_source_title_ar"`
	ResponsibleUnitID               *uint   `json:"responsible_unit_id"`
	ResponsibleUnitTitleEn          *string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr          *string `json:"responsible_unit_title_ar"`
	ResponsibleUserID               *uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      *string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      *string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       *string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       *string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          *string `json:"responsible_user_initial"`
	ResponsibleUserColor            *string `json:"responsible_user_color"`
	ResponsibleUserEmail            *string `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath *string `json:"responsible_user_profile_image_path"`
	RiskStatusTitleEn               string  `json:"risk_status_title_en"`
	RiskStatusTitleAr               string  `json:"risk_status_title_ar"`
	RiskResponseStrategyTitleEn     string  `json:"risk_response_strategy_title_en"`
	RiskResponseStrategyTitleAr     string  `json:"risk_response_strategy_title_ar"`
	RiskStatusType                  string  `json:"risk_status_type"`
	RiskStatusColorCode             string  `json:"risk_status_color_code"`
	Title                           string  `json:"title"`
	Description                     *string `json:"description"`
	ReportedAt                      *string `json:"reported_at"`
	ReviewFrequencyTitleEn          *string `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr          *string `json:"review_frequency_title_ar"`
	NextReviewDate                  *string `json:"next_review_date"`
	CurrentRiskRatingValueID        *uint   `json:"current_risk_rating_value_id"`
	Comments                        *string `json:"comments"`
	SeverityLevelTitleEn            *string `json:"severity_level_title_en"`
	SeverityLevelTitleAr            *string `json:"severity_level_title_ar"`
	CreatedAt                       string  `json:"created_at"`
	UpdatedAt                       *string `json:"updated_at"`
	CreatedUserID                   uint    `json:"created_user_id"`
	CreatedUserFirstNameEn          string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr          string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn           string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr           string  `json:"created_user_last_name_ar"`
	CreatedUserInitial              string  `json:"created_user_initial"`
	CreatedUserColor                string  `json:"created_user_color"`
	CreatedUserEmail                string  `json:"created_user_email"`
	CreatedUserProfileImagePath     string  `json:"created_user_profile_image_path"`
	UpdatedUserID                   *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn          *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr          *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn           *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr           *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial              string  `json:"updated_user_initial"`
	UpdatedUserColor                string  `json:"updated_user_color"`
	UpdatedUserEmail                string  `json:"updated_user_email"`
	UpdatedUserProfileImagePath     string  `json:"updated_user_profile_image_path"`
	CurrentLikelihoodTitleEn        *string `json:"current_likelihood_title_en"`
	CurrentLikelihoodTitleAr        *string `json:"current_likelihood_title_ar"`
	CurrentImpactTitleEn            *string `json:"current_impact_title_en"`
	CurrentImpactTitleAr            *string `json:"current_impact_title_ar"`
	CurrentRiskRatingTitleEn        *string `json:"current_risk_rating_title_en"`
	CurrentRiskRatingTitleAr        *string `json:"current_risk_rating_title_ar"`
	CurrentRiskRatingValue          *uint   `json:"current_risk_rating_value"`
	CurrentRiskRatingColorCode      *string `json:"current_risk_rating_color_code"`
	InherentLikelihoodTitleEn       *string `json:"inherent_likelihood_title_en"`
	InherentLikelihoodTitleAr       *string `json:"inherent_likelihood_title_ar"`
	InherentImpactTitleEn           *string `json:"inherent_impact_title_en"`
	InherentImpactTitleAr           *string `json:"inherent_impact_title_ar"`
	InherentRiskRatingTitleEn       *string `json:"inherent_risk_rating_title_en"`
	InherentRiskRatingTitleAr       *string `json:"inherent_risk_rating_title_ar"`
	InherentRiskRatingValue         *uint   `json:"inherent_risk_rating_value"`
	InherentRiskRatingColorCode     *string `json:"inherent_risk_rating_color_code"`
	ResidualLikelihoodTitleEn       *string `json:"residual_likelihood_title_en"`
	ResidualLikelihoodTitleAr       *string `json:"residual_likelihood_title_ar"`
	ResidualImpactTitleEn           *string `json:"residual_impact_title_en"`
	ResidualImpactTitleAr           *string `json:"residual_impact_title_ar"`
	ResidualRiskRatingTitleEn       *string `json:"residual_risk_rating_title_en"`
	ResidualRiskRatingTitleAr       *string `json:"residual_risk_rating_title_ar"`
	ResidualRiskRatingValue         *uint   `json:"residual_risk_rating_value"`
	ResidualRiskRatingColorCode     *string `json:"residual_risk_rating_color_code"`
	RiskCode                        *string `json:"risk_code"`
	RegistrationDate                *string `json:"registration_date"`
	DivisionID                      *uint   `json:"division_id"`
	DivisionTitleEn                 *string `json:"division_title_en"`
	DivisionTitleAr                 *string `json:"division_title_ar"`
	SectionID                       *uint   `json:"section_id"`
	SectionTitleEn                  *string `json:"section_title_en"`
	SectionTitleAr                  *string `json:"section_title_ar"`
	SubSectionID                    *uint   `json:"sub_section_id"`
	SubSectionTitleEn               *string `json:"sub_section_title_en"`
	SubSectionTitleAr               *string `json:"sub_section_title_ar"`
	IsStrategic                     bool    `json:"is_strategic"`
	IsOperational                   bool    `json:"is_operational"`
	ExpectedCause                   *string `json:"expected_cause"`
	Consequence                     *string `json:"consequence"`
	IsInternal                      bool    `json:"is_internal"`
	IsExternal                      bool    `json:"is_external"`
}

type RiskRegisterListQueryParamsDTO struct {
	Sort                         string
	SortBy                       string
	Keyword                      string
	Status                       string
	Page                         int
	PageSize                     int
	Value                        int
	Filters                      map[string][]string
	QuickFilter                  []string
	ReportedAtFrom               string
	ReportedAtTo                 string
	NextReviewDateFrom           string
	NextReviewDateTo             string
	Export                       bool
	DateFrom                     string
	DateTo                       string
	RiskAppetiteToleranceLevelID uint
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
	ID                      int                                `json:"id"`
	ReferenceNumber         string                             `json:"reference_number"`
	RiskCategory            *RiskCategoryDTO                   `json:"risk_category"`
	RiskSource              *RiskSourceDTO                     `json:"risk_source"`
	ResponsibleUnit         *settingsDTOs.OrganizationUnitDTO  `json:"responsible_unit"`
	ResponsibleUser         *organizationDTOs.UserMinimalDTO   `json:"responsible_user"`
	ResponsibleUserID       uint                               `json:"responsible_user_id" binding:"required" gorm:"not null"`
	RiskStatus              RiskStatusDTO                      `json:"risk_status"`
	Title                   string                             `json:"title"`
	Description             *string                            `json:"description"`
	ReportedAtString        *string                            `json:"reported_at"`
	CurrentLikelihood       *RiskMatrixLikelihoodDTO           `json:"current_likelihood"`
	CurrentImpact           *RiskMatrixImpactDTO               `json:"current_impact"`
	CurrentRiskRating       *RiskMatrixRiskRatingDTO           `json:"current_risk_rating"`
	CurrentRiskRatingValue  *RiskMatrixRiskRatingValueDTO      `json:"current_risk_rating_value"`
	InherentLikelihood      *RiskMatrixLikelihoodDTO           `json:"inherent_likelihood"`
	InherentImpact          *RiskMatrixImpactDTO               `json:"inherent_impact"`
	InherentRiskRating      *RiskMatrixRiskRatingDTO           `json:"inherent_risk_rating"`
	InherentRiskRatingValue *RiskMatrixRiskRatingValueDTO      `json:"inherent_risk_rating_value"`
	ResidualLikelihood      *RiskMatrixLikelihoodDTO           `json:"residual_likelihood"`
	ResidualImpact          *RiskMatrixImpactDTO               `json:"residual_impact"`
	ResidualRiskRating      *RiskMatrixRiskRatingDTO           `json:"residual_risk_rating"`
	ResidualRiskRatingValue *RiskMatrixRiskRatingValueDTO      `json:"residual_risk_rating_value"`
	RiskResponseStrategy    *RiskResponseStrategyDTO           `json:"risk_response_strategy"`
	ReviewFrequency         *masterDTOs.ReviewFrequencyDTO     `json:"review_frequency"`
	NextReviewDateString    *string                            `json:"next_review_date"`
	LastReviewNote          *string                            `json:"last_review_note"`
	NextReviewLevel         *int                               `json:"next_review_level"`
	Documents               []*RiskDocumentDTO                 `json:"documents"`
	SeverityLevel           *masterDTOs.SeverityLevelDTO       `json:"severity_level"`
	CreatedAt               time.Time                          `json:"created_at"`
	UpdatedAt               *time.Time                         `json:"updated_at"`
	CreatedByUser           organizationDTOs.UserMinimalDTO    `json:"created_by"`
	UpdatedByUser           organizationDTOs.UserMinimalDTO    `json:"updated_by"`
	SubmittedByUser         *organizationDTOs.UserMinimalDTO   `json:"submitted_by"`
	SubmittedAtString       *string                            `json:"submitted_at"`
	ReviewedByUser          *organizationDTOs.UserMinimalDTO   `json:"reviewed_by"`
	Controls                []*RiskControlDTO                  `json:"controls"`
	ComplianceRequirements  []*RiskComplianceRequirementDTO    `json:"compliance_requirements"`
	Processes               []*RiskProcessDTO                  `json:"processes"`
	StrategicObjectives     []*RiskStrategicObjectiveDTO       `json:"strategic_objectives"`
	SLARegisters            []*RiskSLARegisterDTO              `json:"sla_registers"`
	Kris                    []*RiskKRIDTO                      `json:"kris"`
	AssociatedDocuments     []*RiskAssociatedDocumentDTO       `json:"associated_documents"`
	Assets                  []*RiskAssetDTO                    `json:"assets"`
	Threats                 []*RiskThreatDTO                   `json:"threats"`
	ActionPlans             []*RiskActionPlanDetailDTO         `json:"action_plans"`
	MarkedAsAuditableItem   bool                               `json:"marked_as_auditable_item"`
	AuditableItem           *AuditableItemMinimalRiskDetailDTO `json:"auditable_item"`
	RiskCode                *string                            `json:"risk_code"`
	RegistrationDateString  *string                            `json:"registration_date"`
	Division                *settingsDTOs.OrganizationUnitDTO  `json:"division"`
	Section                 *settingsDTOs.OrganizationUnitDTO  `json:"section"`
	SubSection              *settingsDTOs.OrganizationUnitDTO  `json:"sub_section"`
	IsStrategic             bool                               `json:"is_strategic"`
	IsOperational           bool                               `json:"is_operational"`
	ExpectedCause           *string                            `json:"expected_cause"`
	Consequence             *string                            `json:"consequence"`
	IsInternal              bool                               `json:"is_internal"`
	IsExternal              bool                               `json:"is_external"`
}

type AuditableItemMinimalRiskDetailDTO struct {
	ID                  uint                             `json:"id"`
	Type                string                           `json:"type"`
	TitleEn             *string                          `json:"title_en"`
	TitleAr             *string                          `json:"title_ar"`
	Comment             *string                          `json:"comment"`
	ResponsibleUnitID   uint                             `json:"responsible_unit_id"`
	ResponsibleUserID   uint                             `json:"responsible_user_id"`
	CreatedByID         uint                             `json:"created_by_id"`
	UpdatedByID         *uint                            `json:"updated_by_id"`
	AuditorRiskRatingID uint                             `json:"auditor_risk_rating_id"`
	AuditorRiskRating   *AuditorRiskRatingMinimalRiskDTO `json:"auditor_risk_rating"`
	ResponsibleUnit     settingsDTOs.OrganizationUnitDTO `json:"responsible_unit"`
	CreatedAt           time.Time                        `json:"created_at"`
	UpdatedAt           *time.Time                       `json:"updated_at"`
	ResponsibleUser     *UserMinimalDTO                  `json:"responsible_user"`
	CreatedByUser       *UserMinimalDTO                  `json:"created_by"`
	UpdatedByUser       *UserMinimalDTO                  `json:"updated_by"`
}
type AuditorRiskRatingMinimalRiskDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Type      string `json:"type"`
	Status    string `json:"status"`
}

type RiskAssociatedDocumentDTO struct {
	ID             *uint                               `json:"id"`
	RiskRegisterID uint                                `json:"risk_register_id"`
	DocumentID     uint                                `json:"document_id"`
	Document       organizationDTOs.DocumentDetailsDTO `json:"document"`
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

type RiskAssetDTO struct {
	ID              *uint                            `json:"id"`
	RiskRegisterID  uint                             `json:"risk_register_id"`
	AssetRegisterID uint                             `json:"asset_register_id"`
	AssetRegister   assetDTOs.AssetRegisterDetailDTO `json:"asset"`
}

type RiskThreatDTO struct {
	ID             *uint                     `json:"id"`
	RiskRegisterID uint                      `json:"risk_register_id"`
	ThreatID       uint                      `json:"threat_id"`
	Threat         assetDTOs.ThreatDetailDTO `json:"threat"`
}

type RiskRegisterFormDTO struct {
	RiskCategoryID    *uint                         `json:"risk_category_id"`
	RiskSourceID      *uint                         `json:"risk_source_id"`
	ResponsibleUserID *uint                         `json:"responsible_user_id"`
	ResponsibleUnitID *uint                         `json:"responsible_unit_id" binding:"required"`
	Title             string                        `json:"title" binding:"required,max=255"`
	Description       *string                       `json:"description"`
	ReportedAt        *string                       `json:"reported_at"`
	ReviewFrequencyID *uint                         `json:"review_frequency_id"`
	NextReviewDate    *string                       `json:"next_review_date" binding:"omitempty,datetime=2006-01-02"`
	LastReviewNote    *string                       `json:"last_review_note" binding:"omitempty"`
	Documents         []*commonDTOs.DocumentFormDTO `json:"documents"`
	SLARegisterIDs    []*uint                       `json:"sla_register_ids"`
	AssetIds          []*uint                       `json:"asset_ids"`
	SeverityLevelID   *uint                         `json:"severity_level_id"`

	InherentRiskMatrixLikelihoodID *uint   `json:"inherent_risk_matrix_likelihood_id"`
	InherentRiskMatrixImpactID     *uint   `json:"inherent_risk_matrix_impact_id"`
	RiskResponseStrategyID         *uint   `json:"risk_response_strategy_id"`
	ResidualRiskMatrixLikelihoodID *uint   `json:"residual_risk_matrix_likelihood_id"`
	ResidualRiskMatrixImpactID     *uint   `json:"residual_risk_matrix_impact_id"`
	RiskCode                       *string `json:"risk_code"`
	DivisionID                     uint    `json:"division_id" binding:"required"`
	SectionID                      *uint   `json:"section_id"`
	SubSectionID                   *uint   `json:"sub_section_id"`
	IsStrategic                    bool    `json:"is_strategic"`
	IsOperational                  bool    `json:"is_operational"`
	ExpectedCause                  *string `json:"expected_cause"`
	RegistrationDate               string  `json:"registration_date" binding:"required,datetime=2006-01-02"`
	Consequence                    *string `json:"consequence"`
	IsInternal                     bool    `json:"is_internal"`
	IsExternal                     bool    `json:"is_external"`
}

type RiskRegisterParsedDates struct {
	ReportedAt       *time.Time
	NextReviewDate   *time.Time
	RegistrationDate *time.Time
}

type RiskResponseStrategyFormDTO struct {
	RiskResponseStrategyID uint `json:"risk_response_strategy_id" binding:"required"`
}
