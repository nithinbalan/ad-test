package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"

	"time"
)

type StrategicIntiativeListDTO struct {
	ID            uint    `json:"id"`
	TitleEn       string  `json:"title_en"`
	TitleAr       string  `json:"title_ar"`
	DescriptionEn *string `json:"description_en"`
	DescriptionAr *string `json:"description_ar"`
	StartDate     *string `json:"start_date"`
	EndDate       *string `json:"end_date"`

	ObjectivesEn      *string `json:"objectives_en"`
	ObjectivesAr      *string `json:"objectives_ar"`
	RecommendationsEn *string `json:"recommendations_en"`
	RecommendationsAr *string `json:"recommendations_ar"`

	OwnerDepartmentID      uint    `json:"owner_department_id"`
	OwnerDepartmentTitleEn *string `json:"owner_department_title_en"`
	OwnerDepartmentTitleAr *string `json:"owner_department_title_ar"`

	StrategicInitiativeStatusID        uint    `json:"strategic_initiative_status_id"`
	StrategicInitiativeStatusTitleEn   *string `json:"strategic_initiative_status_title_en"`
	StrategicInitiativeStatusTitleAr   *string `json:"strategic_initiative_status_title_ar"`
	StrategicInitiativeStatusColorCode *string `json:"strategic_initiative_status_color_code"`
	StrategicInitiativeStatusType      *string `json:"strategic_initiative_status_type"`

	CompletionPercentage   *float64 `json:"completion_percentage"`
	KPIs                   *string  `json:"kpis" gorm:"column:kpis"`
	ExpectedBenefitsEn     *string  `json:"expected_benefits_en"`
	ExpectedBenefitsAr     *string  `json:"expected_benefits_ar"`
	Budget                 *float64 `json:"budget"`
	ResourceRequirementsEn *string  `json:"resource_requirements_en"`
	ResourceRequirementsAr *string  `json:"resource_requirements_ar"`
	DependenciesEn         *string  `json:"dependencies_en"`
	DependenciesAr         *string  `json:"dependencies_ar"`

	ResponsibleUserID               uint    `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      *string `json:"responsible_user_first_name_en"`
	ResponsibleUserLastNameEn       *string `json:"responsible_user_last_name_en"`
	ResponsibleUserFirstNameAr      *string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameAr       *string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          *string `json:"responsible_user_initial"`
	ResponsibleUserEmail            *string `json:"responsible_user_email"`
	ResponsibleUserColor            *string `json:"responsible_user_color"`
	ResponsibleUserProfileImagePath *string `json:"responsible_user_profile_image_path"`

	ResponsibleUnitID      uint    `json:"responsible_unit_id"`
	ResponsibleUnitTitleEn *string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr *string `json:"responsible_unit_title_ar"`

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

type StrategicIntiativeFormDTO struct {
	TitleEn           string  `json:"title_en" binding:"required"`
	TitleAr           string  `json:"title_ar" binding:"required"`
	DescriptionEn     *string `json:"description_en"`
	DescriptionAr     *string `json:"description_ar"`
	StartDate         *string `json:"start_date"`
	EndDate           *string `json:"end_date"`
	ResponsibleUserID uint    `json:"responsible_user_id"`
	ResponsibleUnitID uint    `json:"responsible_unit_id" binding:"required"`

	ObjectivesEn           *string  `json:"objectives_en"`
	ObjectivesAr           *string  `json:"objectives_ar"`
	RecommendationsEn      *string  `json:"recommendations_en"`
	RecommendationsAr      *string  `json:"recommendations_ar"`
	KPIs                   *string  `json:"kpis"`
	ExpectedBenefitsEn     *string  `json:"expected_benefits_en"`
	ExpectedBenefitsAr     *string  `json:"expected_benefits_ar"`
	Budget                 *float64 `json:"budget"`
	ResourceRequirementsEn *string  `json:"resource_requirements_en"`
	ResourceRequirementsAr *string  `json:"resource_requirements_ar"`
	DependenciesEn         *string  `json:"dependencies_en"`
	DependenciesAr         *string  `json:"dependencies_ar"`
}

type StrategicIntiativeStatusCompletionUpdateDTO struct {
	StrategicInitiativeStatusID uint    `json:"strategic_initiative_status_id" binding:"required"`
	CompletionPercentage        float64 `json:"completion_percentage" binding:"required"`
}

type StrategicIntiativeListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	StartDate   string
	EndDate     string
	Export      bool
}

type StrategicIntiativeDetailDTO struct {
	ID              uint    `json:"id"`
	TitleEn         string  `json:"title_en"`
	TitleAr         string  `json:"title_ar"`
	DescriptionEn   *string `json:"description_en"`
	DescriptionAr   *string `json:"description_ar"`
	StartDateString *string `json:"start_date"`
	EndDateString   *string `json:"end_date"`

	ObjectivesEn      *string `json:"objectives_en"`
	ObjectivesAr      *string `json:"objectives_ar"`
	RecommendationsEn *string `json:"recommendations_en"`
	RecommendationsAr *string `json:"recommendations_ar"`

	OwnerDepartmentID uint                                  `json:"owner_department_id"`
	OwnerDepartment   organizationDTOs.OrganizationUnit1DTO `json:"owner_department"`

	StrategicInitiativeStatusID uint                         `json:"strategic_initiative_status_id"`
	StrategicInitiativeStatus   StrategicInitiativeStatusDTO `json:"strategic_initiative_status"`

	CompletionPercentage   *float64 `json:"completion_percentage"`
	KPIs                   *string  `json:"kpis"`
	ExpectedBenefitsEn     *string  `json:"expected_benefits_en"`
	ExpectedBenefitsAr     *string  `json:"expected_benefits_ar"`
	Budget                 *float64 `json:"budget"`
	ResourceRequirementsEn *string  `json:"resource_requirements_en"`
	ResourceRequirementsAr *string  `json:"resource_requirements_ar"`
	DependenciesEn         *string  `json:"dependencies_en"`
	DependenciesAr         *string  `json:"dependencies_ar"`

	ResponsibleUserID uint                                  `json:"responsible_user_id"`
	ResponsibleUser   organizationDTOs.UserMinimalDTO       `json:"responsible_user"`
	ResponsibleUnitID uint                                  `json:"responsible_unit_id"`
	ResponsibleUnit   organizationDTOs.OrganizationUnit1DTO `json:"responsible_unit"`

	CreatedAt     time.Time                        `json:"created_at"`
	UpdatedAt     *time.Time                       `json:"updated_at"`
	CreatedByUser organizationDTOs.UserMinimalDTO  `json:"created_by_user"`
	UpdatedByUser *organizationDTOs.UserMinimalDTO `json:"updated_by_user"`
}

type StrategicIntiativeListResponseDTO struct {
	Items      []*StrategicIntiativeListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type StrategicIntiativeParsedDates struct {
	StartDate *time.Time
	EndDate   *time.Time
}
