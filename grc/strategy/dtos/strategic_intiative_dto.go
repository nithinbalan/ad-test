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
}

type StrategicIntiativeDetailDTO struct {
	ID                uint                                  `json:"id"`
	TitleEn           string                                `json:"title_en"`
	TitleAr           string                                `json:"title_ar"`
	DescriptionEn     *string                               `json:"description_en"`
	DescriptionAr     *string                               `json:"description_ar"`
	StartDateString   *string                               `json:"start_date"`
	EndDateString     *string                               `json:"end_date"`
	ResponsibleUserID uint                                  `json:"responsible_user_id"`
	ResponsibleUser   organizationDTOs.UserMinimalDTO       `json:"responsible_user"`
	ResponsibleUnitID uint                                  `json:"responsible_unit_id"`
	ResponsibleUnit   organizationDTOs.OrganizationUnit1DTO `json:"responsible_unit"`
	CreatedAt         time.Time                             `json:"created_at"`
	UpdatedAt         *time.Time                            `json:"updated_at"`
	CreatedByUser     organizationDTOs.UserMinimalDTO       `json:"created_by_user"`
	UpdatedByUser     *organizationDTOs.UserMinimalDTO      `json:"updated_by_user"`
}

type StrategicIntiativeListResponseDTO struct {
	Items      []*StrategicIntiativeListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type StrategicIntiativeParsedDates struct {
	StartDate *time.Time
	EndDate   *time.Time
}
