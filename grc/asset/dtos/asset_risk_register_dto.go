package dtos

import (
	"time"

	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
)

type AssetRiskRegisterFormDTO struct {
	VulnerabilityID      uint   `json:"vulnerability_id" binding:"required"`
	AssetID              *uint  `json:"asset_id"`
	RiskCode             string `json:"risk_code" binding:"required,max=500"`
	RegistrationDate     string `json:"registration_date" binding:"required,datetime=2006-01-02"`
	Title                string `json:"title" binding:"required,max=255"`
	Justification        string `json:"justification"  binding:"required"`
	SectionID            uint   `json:"section_id" binding:"required"`
	DivisionID           uint   `json:"division_id" binding:"required"`
	ResponsibleUserID    uint   `json:"responsible_user_id" binding:"required"`
	IsStrategic          *bool  `json:"is_strategic" binding:"required"`
	IsOperational        *bool  `json:"is_operational" binding:"required"`
	ExpectedCause        string `json:"expected_cause"  binding:"required"`
	RiskCategoryID       uint   `json:"risk_category_id" binding:"required"`
	IsInternal           *bool  `json:"is_internal" binding:"required"`
	IsExternal           *bool  `json:"is_external" binding:"required"`
	ConfidentialityScore int    `json:"confidentiality_score" binding:"required,min=1,max=5"`
	IntegrityScore       int    `json:"integrity_score" binding:"required,min=1,max=5"`
	AvailabilityScore    int    `json:"availability_score" binding:"required,min=1,max=5"`
	CurrentLikelihoodID  uint   `json:"current_likelihood_id" binding:"required"`
}

type AssetRiskRegisterParsedDates struct {
	RegistrationDate time.Time
}

type AssetRiskRegisterListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	Export      bool
}

type AssetRiskRegisterListDTO struct {
	ID                          uint    `json:"id"`
	ReferenceNumber             string  `json:"reference_number"`
	VulnerabilityID             uint    `json:"vulnerability_id"`
	AssetID                     *uint   `json:"asset_id"`
	RiskCode                    string  `json:"risk_code"`
	Title                       string  `json:"title"`
	RegistrationDate            string  `json:"registration_date"`
	Justification               *string `json:"justification"`
	ExpectedCause               *string `json:"expected_cause"`
	RiskCategoryTitleEn         string  `json:"risk_category_title_en"`
	RiskCategoryTitleAr         string  `json:"risk_category_title_ar"`
	DivisionTitleEn             string  `json:"division_title_en"`
	DivisionTitleAr             string  `json:"division_title_ar"`
	SectionTitleEn              string  `json:"section_title_en"`
	SectionTitleAr              string  `json:"section_title_ar"`
	CurrentLikelihoodTitleEn    string  `json:"current_likelihood_title_en"`
	CurrentLikelihoodTitleAr    string  `json:"current_likelihood_title_ar"`
	CurrentImpactTitleEn        string  `json:"current_impact_title_en"`
	CurrentImpactTitleAr        string  `json:"current_impact_title_ar"`
	CurrentRiskRatingTitleEn    string  `json:"current_risk_rating_title_en"`
	CurrentRiskRatingTitleAr    string  `json:"current_risk_rating_title_ar"`
	CurrentRiskRatingColorCode  string  `json:"current_risk_rating_color_code"`
	IsStrategic                 bool    `json:"is_strategic"`
	IsOperational               bool    `json:"is_operational"`
	IsInternal                  bool    `json:"is_internal"`
	IsExternal                  bool    `json:"is_external"`
	CreatedAt                   string  `json:"created_at"`
	UpdatedAt                   *string `json:"updated_at"`
	CreatedUserID               uint    `json:"created_user_id"`
	CreatedUserFirstNameEn      string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string  `json:"created_user_last_name_ar"`
	CreatedUserInitial          string  `json:"created_user_initial"`
	CreatedUserColor            string  `json:"created_user_color"`
	CreatedUserEmail            string  `json:"created_user_email"`
	CreatedUserProfileImagePath string  `json:"created_user_profile_image_path"`
	UpdatedUserID               *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn      *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          *string `json:"updated_user_initial"`
	UpdatedUserColor            *string `json:"updated_user_color"`
	UpdatedUserEmail            *string `json:"updated_user_email"`
	UpdatedUserProfileImagePath *string `json:"updated_user_profile_image_path"`
}

type AssetRiskRegisterListResponseDTO struct {
	Items      []*AssetRiskRegisterListDTO  `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type AssetRiskRegisterDetailDTO struct {
	ID                   uint       `json:"id"`
	ReferenceNumber      string     `json:"reference_number"`
	RiskCode             string     `json:"risk_code"`
	RegistrationDate     time.Time  `json:"registration_date"`
	Title                string     `json:"title"`
	Justification        *string    `json:"justification"`
	IsStrategic          bool       `json:"is_strategic"`
	IsOperational        bool       `json:"is_operational"`
	IsInternal           bool       `json:"is_internal"`
	IsExternal           bool       `json:"is_external"`
	ExpectedCause        *string    `json:"expected_cause"`
	ConfidentialityScore int        `json:"confidentiality_score"`
	IntegrityScore       int        `json:"integrity_score"`
	AvailabilityScore    int        `json:"availability_score"`
	IsCyberRisk          bool       `json:"is_cyber_risk"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            *time.Time `json:"updated_at"`

	RiskStatus             RiskStatusDTO                    `json:"risk_status"`
	Section                settingsDTOs.OrganizationUnitDTO `json:"section"`
	Division               settingsDTOs.OrganizationUnitDTO `json:"division"`
	ResponsibleUser        organizationDTOs.UserMinimalDTO  `json:"responsible_user"`
	RiskCategory           RiskCategoryDTO                  `json:"risk_category"`
	CurrentLikelihood      RiskMatrixLikelihoodDTO          `json:"current_likelihood"`
	CurrentImpact          RiskMatrixImpactDTO              `json:"current_impact"`
	CurrentRiskRating      RiskMatrixRiskRatingDTO          `json:"current_risk_rating"`
	CurrentRiskRatingValue RiskMatrixRiskRatingValueDTO     `json:"current_risk_rating_value"`
	CreatedByUser          organizationDTOs.UserMinimalDTO  `json:"created_by"`

	UpdatedByUser   *organizationDTOs.UserMinimalDTO `json:"updated_by"`
	Vulnerabilities []*RiskVulnerabilityDTO          `json:"vulnerabilities"`
	Assets          []*RiskAssetDTO                  `json:"assets"`
}

type RiskCategoryDTO struct {
	ID      int    `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Status  string `json:"status"`
}

type RiskStatusDTO struct {
	ID            int    `json:"id"`
	TitleEn       string `json:"title_en"`
	TitleAr       string `json:"title_ar"`
	DescriptionEn string `json:"description_en"`
	DescriptionAr string `json:"description_ar"`
	ColorCode     string `json:"color_code"`
	Type          string `json:"type"`
	Status        string `json:"status"`
}

type RiskAssetDTO struct {
	ID              *uint                  `json:"id"`
	RiskRegisterID  uint                   `json:"risk_register_id"`
	AssetRegisterID uint                   `json:"asset_register_id"`
	AssetRegister   AssetRegisterDetailDTO `json:"asset"`
}

type RiskVulnerabilityDTO struct {
	ID              *uint                   `json:"id"`
	RiskRegisterID  uint                    `json:"risk_register_id"`
	VulnerabilityID uint                    `json:"vulnerability_id"`
	Vulnerability   VulnerabilityDetailsDTO `json:"vulnerability"`
}

type RiskMatrixLikelihoodDTO struct {
	ID            uint   `json:"id"`
	RiskMatrixID  uint   `json:"risk_matrix_id"`
	TitleEn       string `json:"title_en"`
	TitleAr       string `json:"title_ar"`
	Score         uint   `json:"score"`
	DescriptionEn string `json:"description_en"`
	DescriptionAr string `json:"description_ar"`
	Probability   string `json:"probability"`
	Status        string `json:"status"`
}

type RiskMatrixImpactDTO struct {
	ID            uint   `json:"id"`
	RiskMatrixID  uint   `json:"risk_matrix_id"`
	TitleEn       string `json:"title_en"`
	TitleAr       string `json:"title_ar"`
	Score         uint   `json:"score"`
	DescriptionEn string `json:"description_en"`
	DescriptionAr string `json:"description_ar"`
	Status        string `json:"status"`
}

type RiskMatrixRiskRatingDTO struct {
	ID           uint   `json:"id"`
	RiskMatrixID uint   `json:"risk_matrix_id"`
	TitleEn      string `json:"title_en"`
	TitleAr      string `json:"title_ar"`
	ColorCode    string `json:"color_code"`
	Type         string `json:"type"`
	Status       string `json:"status"`
}

type RiskMatrixRiskRatingValueDTO struct {
	ID                     uint `json:"id"`
	RiskMatrixRiskRatingID uint `json:"risk_matrix_risk_rating_id"`
	Value                  uint `json:"value"`
}