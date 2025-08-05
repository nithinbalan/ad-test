package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"

	"time"
)

type StrategicObjectiveListDTO struct {
	ID              uint    `json:"id"`
	TitleEn         string  `json:"title_en"`
	TitleAr         string  `json:"title_ar"`
	ReferenceNumber string  `json:"reference_number"`
	DescriptionEn   *string `json:"description_en"`
	DescriptionAr   *string `json:"description_ar"`
	StartDate       *string `json:"start_date"`
	EndDate         *string `json:"end_date"`
	Status          string  `json:"status"`

	StrategicObjectiveCategoryID      uint   `json:"strategic_objective_category_id"`
	StrategicObjectiveCategoryTitleEn string `json:"strategic_objective_category_title_en"`
	StrategicObjectiveCategoryTitleAr string `json:"strategic_objective_category_title_ar"`

	ResponsibleUserID               *uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      *string `json:"responsible_user_first_name_en"`
	ResponsibleUserLastNameEn       *string `json:"responsible_user_last_name_en"`
	ResponsibleUserFirstNameAr      *string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameAr       *string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          *string `json:"responsible_user_initial"`
	ResponsibleUserEmail            *string `json:"responsible_user_email"`
	ResponsibleUserColor            *string `json:"responsible_user_color"`
	ResponsibleUserProfileImagePath *string `json:"responsible_user_profile_image_path"`

	ResponsibleUnitID      *uint   `json:"responsible_unit_id"`
	ResponsibleUnitTitleEn *string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr *string `json:"responsible_unit_title_ar"`

	ReviewUserID          *uint   `json:"review_user_id"`
	ReviewUserFirstNameEn *string `json:"review_user_first_name_en"`
	ReviewUserLastNameEn  *string `json:"review_user_last_name_en"`
	ReviewUserFirstNameAr *string `json:"review_user_first_name_ar"`
	ReviewUserLastNameAr  *string `json:"review_user_last_name_ar"`

	ReviewFrequencyID      *uint   `json:"review_frequency_id"`
	ReviewFrequencyTitleEn *string `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr *string `json:"review_frequency_title_ar"`

	CreatedByID                 uint   `json:"created_by_id"`
	CreatedUserFirstNameEn      string `json:"created_user_first_name_en"`
	CreatedUserLastNameEn       string `json:"created_user_last_name_en"`
	CreatedUserFirstNameAr      string `json:"created_user_first_name_ar"`
	CreatedUserLastNameAr       string `json:"created_user_last_name_ar"`
	CreatedUserInitial          string `json:"created_user_initial"`
	CreatedUserColor            string `json:"created_user_color"`
	CreatedUserEmail            string `json:"created_user_email"`
	CreatedUserProfileImagePath string `json:"created_user_profile_image_path"`

	UpdatedByID                 *uint   `json:"updated_by_id"`
	UpdatedUserFirstNameEn      *string `json:"updated_user_first_name_en"`
	UpdatedUserLastNameEn       *string `json:"updated_user_last_name_en"`
	UpdatedUserFirstNameAr      *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameAr       *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          *string `json:"updated_user_initial"`
	UpdatedUserColor            *string `json:"updated_user_color"`
	UpdatedUserEmail            *string `json:"updated_user_email"`
	UpdatedUserProfileImagePath *string `json:"updated_user_profile_image_path"`
}

type StrategicObjectiveFormDTO struct {
	TitleEn                      string  `json:"title_en" binding:"required"`
	TitleAr                      string  `json:"title_ar" binding:"required"`
	StrategicObjectiveCategoryID uint    `json:"strategic_objective_category_id" binding:"required"`
	DescriptionEn                *string `json:"description_en"`
	DescriptionAr                *string `json:"description_ar"`
	StartDate                    *string `json:"start_date"`
	EndDate                      *string `json:"end_date"`
	ResponsibleUserID            *uint   `json:"responsible_user_id"`
	ReviewUserID                 *uint   `json:"review_user_id"`
	ReviewFrequencyID            *uint   `json:"review_frequency_id"`
	Status                       string  `json:"status"`
}

type StrategicObjectiveListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	StartDate   string
	EndDate     string
}

type StrategicObjectiveDetailDTO struct {
	ID                           uint                                   `json:"id"`
	TitleEn                      string                                 `json:"title_en"`
	TitleAr                      string                                 `json:"title_ar"`
	ReferenceNumber              string                                 `json:"reference_number"`
	DescriptionEn                *string                                `json:"description_en"`
	DescriptionAr                *string                                `json:"description_ar"`
	StartDateString              *string                                `json:"start_date"`
	EndDateString                *string                                `json:"end_date"`
	Status                       string                                 `json:"status"`
	StrategicObjectiveCategoryID uint                                   `json:"strategic_objective_category_id"`
	StrategicObjectiveCategory   *StrategicObjectiveCategoryDTO         `json:"strategic_objective_category"`
	ResponsibleUserID            *uint                                  `json:"responsible_user_id"`
	ResponsibleUser              *organizationDTOs.UserMinimalDTO       `json:"responsible_user"`
	ResponsibleUnitID            *uint                                  `json:"responsible_unit_id"`
	ResponsibleUnit              *organizationDTOs.OrganizationUnit1DTO `json:"responsible_unit"`
	ReviewUserID                 *uint                                  `json:"review_user_id"`
	ReviewUser                   *organizationDTOs.UserMinimalDTO       `json:"review_user"`
	ReviewFrequencyID            *uint                                  `json:"review_frequency_id"`
	ReviewFrequency              *masterDTOs.ReviewFrequencyDTO         `json:"review_frequency"`
	CreatedAt                    time.Time                              `json:"created_at"`
	UpdatedAt                    *time.Time                             `json:"updated_at"`
	CreatedByUser                organizationDTOs.UserMinimalDTO        `json:"created_by_user"`
	UpdatedByUser                *organizationDTOs.UserMinimalDTO       `json:"updated_by_user"`
	StrategicObjectiveKPIs       []StrategicObjectiveKPIDetailDTO       `json:"strategic_objective_kpis"`
}
type StrategicObjectiveKPIDetailDTO struct {
	ID            uint                            `json:"id"`
	KPIRegisterID uint                            `json:"kpi_register_id"`
	KPIRegister   KPIRegisterDetailDTO            `json:"kpi_register"`
	CreatedAt     time.Time                       `json:"created_at"`
	CreatedBy     organizationDTOs.UserMinimalDTO `json:"created_by"`
}

type StrategicObjectiveListResponseDTO struct {
	Items      []*StrategicObjectiveListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type StrategicObjectiveParsedDates struct {
	StartDate *time.Time
	EndDate   *time.Time
}

// DTO for focus area mapped to a strategic objective
type StrategicObjectiveFocusAreaDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Status  string `json:"status"`
}

type StrategicObjectiveRiskSummaryDTO struct {
	TotalRiskCount   *int     `json:"total_risk_count,omitempty"`
	HighRiskCount    *int     `json:"high_risk_count,omitempty"`
	ControlCoverage  *float64 `json:"control_coverage,omitempty"`
	ControlRiskCount *int     `json:"control_risk_count,omitempty"`

	LastReviewed *string `json:"last_reviewed"`

	TopRisk struct {
		ID                *uint   `json:"id,omitempty"`
		Title             *string `json:"title,omitempty"`
		Rating            *string `json:"rating,omitempty"`
		Description       *string `json:"description,omitempty"`
		ImpactTitleEn     *string `json:"impact_title_en,omitempty"`
		ImpactTitleAr     *string `json:"impact_title_ar,omitempty"`
		LikelihoodTitleEn *string `json:"likelihood_title_en,omitempty"`
		LikelihoodTitleAr *string `json:"likelihood_title_ar,omitempty"`
	} `json:"top_risk"`
}
