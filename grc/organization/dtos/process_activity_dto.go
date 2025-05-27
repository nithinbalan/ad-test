package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	settingsDTOs "jethur/grc/settings/dtos"

	"time"
)

type ProcessActivityListDTO struct {
	ID                uint       `json:"id"`
	TitleEn           string     `json:"title_en"`
	TitleAr           string     `json:"title_ar"`
	DescriptionEn     *string    `json:"description_en"`
	DescriptionAr     *string    `json:"description_ar"`
	Sequence          int        `json:"sequence"`
	EstimatedDuration *string    `json:"estimated_duration"`
	ExpectedOutput    string     `json:"expected_output"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at"`

	CreatedUserID           uint    `json:"created_user_id"`
	CreatedUserFirstNameEn  string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr  string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn   string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr   string  `json:"created_user_last_name_ar"`
	CreatedUserInitial      string  `json:"created_user_initial"`
	CreatedUserColor        string  `json:"created_user_color"`
	CreatedUserEmail        string  `json:"created_user_email"`
	CreatedUserProfileImage string  `json:"created_user_profile_image_path"`
	UpdatedUserID           *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn  *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr  *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn   *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr   *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial      string  `json:"updated_user_initial"`
	UpdatedUserColor        string  `json:"updated_user_color"`
	UpdatedUserEmail        string  `json:"updated_user_email"`
	UpdatedUserProfileImage string  `json:"updated_user_profile_image_path"`
}

type ProcessActivityListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	Export      bool
	DateFrom    string
	DateTo      string
}

type ProcessActivityListResponseDTO struct {
	Items      []*ProcessActivityListDTO    `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type ProcessActivityFormDTO struct {
	TitleEn           string  `json:"title_en" binding:"required"`
	TitleAr           string  `json:"title_ar" binding:"required"`
	DescriptionEn     *string `json:"description_en"`
	DescriptionAr     *string `json:"description_ar"`
	Sequence          int     `json:"sequence" binding:"required"`
	EstimatedDuration *string `json:"estimated_duration"`
	ExpectedOutput    string  `json:"expected_output" binding:"required"`
}

type ProcessActivityDetailDTO struct {
	ID                uint                      `json:"id"`
	TitleEn           string                    `json:"title_en"`
	TitleAr           string                    `json:"title_ar"`
	DescriptionEn     *string                   `json:"description_en,omitempty"`
	DescriptionAr     *string                   `json:"description_ar,omitempty"`
	Sequence          int                       `json:"sequence"`
	EstimatedDuration *string                   `json:"estimated_duration,omitempty"`
	ExpectedOutput    string                    `json:"expected_output"`
	CreatedAt         time.Time                 `json:"created_at"`
	UpdatedAt         *time.Time                `json:"updated_at"`
	CreatedByUser     UserMinimalDTO            `json:"created_by_user"`
	UpdatedByUser     *UserMinimalDTO           `json:"updated_by_user,omitempty"`
	Process           *ProcessMinimalDTO        `json:"process"`
	Users             []*ProcessActivityUserDTO `json:"users"`
}

type ProcessActivityUserDTO struct {
	ID              uint                             `json:"id"`
	UserID          uint                             `json:"user_id"`
	RoleType        string                           `json:"role_type"` // responsible, accountable, etc.
	User            UserMinimalDTO                   `json:"user"`
	ResponsibleUnit settingsDTOs.OrganizationUnitDTO `json:"responsible_unit"`
}
type ProcessMinimalDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
}

type ProcessActivityUserRACIDTO struct {
	ResponsibleUserIDs []uint `json:"responsible_user_ids"`
	AccountableUserIDs []uint `json:"accountable_user_ids"`
	ConsultedUserIDs   []uint `json:"consulted_user_ids"`
	InformedUserIDs    []uint `json:"informed_user_ids"`
}
