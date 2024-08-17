package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	settingsDTOs "jethur/grc/settings/dtos"

	"time"
)

type ControlListDTO struct {
	ID                                   uint       `json:"id"`
	ControlTypeTitleEn                   string     `json:"control_type_title_en"`
	ControlTypeTitleAr                   string     `json:"control_type_title_ar"`
	ControlCategoryTitleEn               string     `json:"control_category_title_en"`
	ControlCategoryTitleAr               string     `json:"control_category_title_ar"`
	OwnerUnitTitleEn                     *string    `json:"owner_unit_title_en"`
	OwnerUnitTitleAr                     *string    `json:"owner_unit_title_ar"`
	OwnerID                              *uint      `json:"owner_id"`
	OwnerFirstNameEn                     *string    `json:"owner_first_name_en"`
	OwnerFirstNameAr                     *string    `json:"owner_first_name_ar"`
	OwnerLastNameEn                      *string    `json:"owner_last_name_en"`
	OwnerLastNameAr                      *string    `json:"owner_last_name_ar"`
	OwnerInitial                         string     `json:"owner_initial"`
	OwnerColor                           string     `json:"owner_color"`
	OwnerEmail                           string     `json:"owner_email"`
	OwnerProfileImagePath                string     `json:"owner_profile_image_path"`
	ControlImplementationStatusTitleEn   *string    `json:"control_implementation_status_title_en"`
	ControlImplementationStatusTitleAr   *string    `json:"control_implementation_status_title_ar"`
	ControlImplementationStatusColorCode *string    `json:"control_implementation_status_color_code"`
	ControlEfficiencyStatusTitleEn       *string    `json:"control_efficiency_status_title_en"`
	ControlEfficiencyStatusTitleAr       *string    `json:"control_efficiency_status_title_ar"`
	ControlEfficiencyStatusColorCode     *string    `json:"control_efficiency_status_color_code"`
	ReferenceNumber                      string     `json:"reference_number"`
	Code                                 *string    `json:"code"`
	TitleEn                              string     `json:"title_en"`
	TitleAr                              string     `json:"title_ar"`
	Description                          *string    `json:"description"`
	Objective                            *string    `json:"objective"`
	ReviewFrequencyTitleEn               string     `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr               string     `json:"review_frequency_title_ar"`
	LastReviewDate                       *time.Time `json:"last_review_date"`
	NextReviewDate                       *time.Time `json:"next_review_date"`
	LastReviewNote                       *string    `json:"last_review_note"`
	TestResponsibleUserID                uint       `json:"test_responsible_user_id"`
	TestResponsibleUserFirstNameEn       string     `json:"test_responsible_user_first_name_en"`
	TestResponsibleUserFirstNameAr       string     `json:"test_responsible_user_first_name_ar"`
	TestResponsibleUserLastNameEn        string     `json:"test_responsible_user_last_name_en"`
	TestResponsibleUserLastNameAr        string     `json:"test_responsible_user_last_name_ar"`
	TestResponsibleUserInitial           string     `json:"test_responsible_user_initial"`
	TestResponsibleUserColor             string     `json:"test_responsible_user_color"`
	TestResponsibleUserEmail             string     `json:"test_responsible_user_email"`
	TestResponsibleUserProfileImagePath  string     `json:"test_responsible_user_profile_image_path"`
	TestFrequencyTitleEn                 *string    `json:"test_frequency_title_en"`
	TestFrequencyTitleAr                 *string    `json:"test_frequency_title_ar"`
	LastTestDate                         *time.Time `json:"last_test_date"`
	NextTestDate                         *time.Time `json:"next_test_date"`
	ImplementedAt                        *time.Time `json:"implemented_at"`
	LifecycleStatusTitleEn               string     `json:"lifecycle_status_title_en"`
	LifecycleStatusTitleAr               string     `json:"lifecycle_status_title_ar"`
	LifecycleStatusColorCode             string     `json:"lifecycle_status_color_code"`
	CreatedUserID                        uint       `json:"created_user_id"`
	CreatedUserFirstNameEn               string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr               string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn                string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr                string     `json:"created_user_last_name_ar"`
	CreatedUserInitial                   string     `json:"created_user_initial"`
	CreatedUserColor                     string     `json:"created_user_color"`
	CreatedUserEmail                     string     `json:"created_user_email"`
	CreatedUserProfileImagePath          string     `json:"created_user_profile_image_path"`
	UpdatedUserID                        *uint      `json:"updated_user_id"`
	UpdatedUserFirstNameEn               *string    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr               *string    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn                *string    `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr                *string    `json:"updated_user_last_name_ar"`
	UpdatedUserInitial                   string     `json:"updated_user_initial"`
	UpdatedUserColor                     string     `json:"updated_user_color"`
	UpdatedUserEmail                     string     `json:"updated_user_email"`
	UpdatedUserProfileImagePath          string     `json:"updated_user_profile_image_path"`
	CreatedAt                            time.Time  `json:"created_at"`
	UpdatedAt                            *time.Time `json:"updated_at"`
}

type ControlMinimalListDTO struct {
	ID              uint    `json:"id"`
	ReferenceNumber string  `json:"reference_number"`
	Code            *string `json:"code"`
	TitleEn         string  `json:"title_en"`
	TitleAr         string  `json:"title_ar"`
	Description     *string `json:"description"`
}

type ControlListQueryParamsDTO struct {
	Sort               string
	SortBy             string
	Keyword            string
	Status             string
	Page               int
	PageSize           int
	Filters            map[string][]string
	QuickFilter        []string
	LastReviewDateFrom string
	LastReviewDateTo   string
	NextReviewDateFrom string
	NextReviewDateTo   string
	LastTestDateFrom   string
	LastTestDateTo     string
	NextTestDateFrom   string
	NextTestDateTo     string
	ImplementedAtFrom  string
	ImplementedAtTo    string
}

type ControlListResponseDTO struct {
	Items      []*ControlListDTO            `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type ControlFormDTO struct {
	ControlTypeID                 uint    `json:"control_type_id" binding:"required"`
	ControlCategoryID             uint    `json:"control_category_id" binding:"required"`
	OwnerID                       uint    `json:"owner_id" binding:"required"`
	ControlImplementationStatusID *uint   `json:"control_implementation_status_id" binding:"omitempty"`
	ControlEfficiencyStatusID     *uint   `json:"control_efficiency_status_id" binding:"omitempty"`
	Code                          *string `json:"code" binding:"omitempty"`
	TitleEn                       string  `json:"title_en" binding:"required"`
	TitleAr                       string  `json:"title_ar" binding:"required"`
	Description                   *string `json:"description" binding:"omitempty"`
	Objective                     *string `json:"objective" binding:"omitempty"`
	EfficiencyScore               *uint   `json:"efficiency_score" binding:"omitempty"`
	ReviewFrequencyID             *uint   `json:"review_frequency_id" binding:"omitempty"`
	LastReviewDate                *string `json:"last_review_date" binding:"omitempty,datetime=2006-01-02"`
	NextReviewDate                *string `json:"next_review_date" binding:"omitempty,datetime=2006-01-02"`
	LastReviewNote                *string `json:"last_review_note" binding:"omitempty"`
	TestResponsibleUserID         uint    `json:"test_responsible_user_id" binding:"required"`
	TestFrequencyID               *uint   `json:"test_frequency_id,omitempty"`
	LastTestDate                  *string `json:"last_test_date" binding:"omitempty,datetime=2006-01-02"`
	NextTestDate                  *string `json:"next_test_date" binding:"omitempty,datetime=2006-01-02"`
	ImplementedAt                 *string `json:"implemented_at" binding:"omitempty,datetime=2006-01-02"`
}

type ControlParsedDates struct {
	LastReviewDate *time.Time
	NextReviewDate *time.Time
	LastTestDate   *time.Time
	NextTestDate   *time.Time
	ImplementedAt  *time.Time
}

type ControlDetailDTO struct {
	ID                          uint                             `json:"id"`
	ReferenceNumber             string                           `json:"reference_number"`
	Code                        string                           `json:"code"`
	TitleEn                     string                           `json:"title_en"`
	TitleAr                     string                           `json:"title_ar"`
	Description                 *string                          `json:"description"`
	Objective                   *string                          `json:"objective"`
	ControlType                 ControlTypeDTO                   `json:"control_type"`
	ControlCategory             ControlCategoryDTO               `json:"control_category"`
	OwnerUnit                   settingsDTOs.OrganizationUnitDTO `json:"owner_unit"`
	Owner                       *UserMinimalDTO                  `json:"owner"`
	ControlImplementationStatus *ControlImplementationStatusDTO  `json:"control_implementation_status"`
	ControlEfficiencyStatus     *ControlEfficiencyStatusDTO      `json:"control_efficiency_status"`
	EfficiencyScore             *uint                            `json:"efficiency_score"`
	ReviewFrequency             *masterDTOs.ReviewFrequencyDTO   `json:"review_frequency"`
	LastReviewDateString        *string                          `json:"last_review_date"`
	NextReviewDateString        *string                          `json:"next_review_date"`
	LastReviewNote              *string                          `json:"last_review_note"`
	TestResponsibleUser         UserMinimalDTO                   `json:"test_responsible_user"`
	TestFrequency               *masterDTOs.ReviewFrequencyDTO   `json:"test_frequency"`
	LastTestDateString          *string                          `json:"last_test_date"`
	NextTestDateString          *string                          `json:"next_test_date"`
	ImplementedAtString         *string                          `json:"implemented_at"`
	CreatedAt                   time.Time                        `json:"created_at"`
	UpdatedAt                   *time.Time                       `json:"updated_at"`
	CreatedByUser               UserMinimalDTO                   `json:"created_by"`
	UpdatedByUser               *UserMinimalDTO                  `json:"updated_by"`
	LifecycleStatus             masterDTOs.LifecycleStatusDTO    `json:"status"`
	ControlReviewStatus         ControlReviewStatusDTO           `json:"review_status"`
	LatestControlTest           ControlTestDetailsDTO            `json:"latest_control_test"`
	LatestControlReview         *ControlReviewDTO                `json:"control_review"`
}
type ControlReviewDTO struct {
	ID             uint                          `json:"id"`
	Comment        string                        `json:"comment"`
	ReviewedAt     time.Time                     `json:"reviewed_at"`
	UpdatedAt      *time.Time                    `json:"updated_at"`
	ReviewedByUser UserMinimalDTO                `json:"reviewed_by"`
	UpdatedByUser  *UserMinimalDTO               `json:"updated_by"`
	Documents      []*commonDTOs.DocumentFormDTO `json:"documents"`
}
