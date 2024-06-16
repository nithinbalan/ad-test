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
	OwnerFirstNameEn                     *string    `json:"owner_first_name_en"`
	OwnerFirstNameAr                     *string    `json:"owner_first_name_ar"`
	OwnerLastNameEn                      *string    `json:"owner_last_name_en"`
	OwnerLastNameAr                      *string    `json:"owner_last_name_ar"`
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
	TestResponsibleUserFirstNameEn       string     `json:"test_responsible_user_first_name_en"`
	TestResponsibleUserFirstNameAr       string     `json:"test_responsible_user_first_name_ar"`
	TestResponsibleUserLastNameEn        string     `json:"test_responsible_user_last_name_en"`
	TestResponsibleUserLastNameAr        string     `json:"test_responsible_user_last_name_ar"`
	TestFrequencyTitleEn                 *string    `json:"test_frequency_title_en"`
	TestFrequencyTitleAr                 *string    `json:"test_frequency_title_ar"`
	LastTestDate                         *time.Time `json:"last_test_date"`
	NextTestDate                         *time.Time `json:"next_test_date"`
	ImplementedAt                        *time.Time `json:"implemented_at"`
	LifecycleStatusTitleEn               string     `json:"lifecycle_status_title_en"`
	LifecycleStatusTitleAr               string     `json:"lifecycle_status_title_ar"`
	LifecycleStatusColorCode             string     `json:"lifecycle_status_color_code"`
	CreatedUserFirstNameEn               string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr               string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn                string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr                string     `json:"created_user_last_name_ar"`
	UpdatedUserFirstNameEn               *string    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr               *string    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn                *string    `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr                *string    `json:"updated_user_last_name_ar"`
	CreatedAt                            time.Time  `json:"created_at"`
	UpdatedAt                            *time.Time `json:"updated_at"`
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
	OwnerUnitID                   *uint   `json:"owner_unit_id" binding:"omitempty"`
	OwnerID                       *uint   `json:"owner_id" binding:"omitempty"`
	ControlImplementationStatusID *uint   `json:"control_implementation_status_id" binding:"omitempty"`
	ControlEfficiencyStatusID     *uint   `json:"control_efficiency_status_id" binding:"omitempty"`
	Code                          *string `json:"code" binding:"omitempty"`
	TitleEn                       string  `json:"title_en" binding:"required"`
	TitleAr                       string  `json:"title_ar" binding:"required"`
	Description                   *string `json:"description" binding:"omitempty"`
	Objective                     *string `json:"objective" binding:"omitempty"`
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
	ID                          uint                              `json:"id"`
	ReferenceNumber             string                            `json:"reference_number"`
	Code                        string                            `json:"code"`
	TitleEn                     string                            `json:"title_en"`
	TitleAr                     string                            `json:"title_ar"`
	Description                 *string                           `json:"description"`
	Objective                   *string                           `json:"objective"`
	ControlType                 ControlTypeDTO                    `json:"control_type"`
	ControlCategory             ControlCategoryDTO                `json:"control_category"`
	OwnerUnit                   *settingsDTOs.OrganizationUnitDTO `json:"owner_unit"`
	Owner                       *UserMinimalDTO                   `json:"owner"`
	ControlImplementationStatus *ControlImplementationStatusDTO   `json:"control_implementation_status"`
	ControlEfficiencyStatus     *ControlEfficiencyStatusDTO       `json:"control_efficiency_status"`
	ReviewFrequency             *masterDTOs.ReviewFrequencyDTO    `json:"review_frequency"`
	LastReviewDateString        *string                           `json:"last_review_date"`
	NextReviewDateString        *string                           `json:"next_review_date"`
	LastReviewNote              *string                           `json:"last_review_note"`
	TestResponsibleUser         UserMinimalDTO                    `json:"test_responsible_user"`
	TestFrequency               *masterDTOs.ReviewFrequencyDTO    `json:"test_frequency"`
	LastTestDateString          *string                           `json:"last_test_date"`
	NextTestDateString          *string                           `json:"next_test_date"`
	ImplementedAtString         *string                           `json:"implemented_at"`
	CreatedAt                   time.Time                         `json:"created_at"`
	UpdatedAt                   *time.Time                        `json:"updated_at"`
	CreatedByUser               UserMinimalDTO                    `json:"created_by"`
	UpdatedByUser               *UserMinimalDTO                   `json:"updated_by"`
	LifecycleStatus             masterDTOs.LifecycleStatusDTO     `json:"status"`
}
