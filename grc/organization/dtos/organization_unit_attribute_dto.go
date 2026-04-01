package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type OrganizationUnitAttributeListDTO struct {
	ID                          uint       `json:"id"`
	DepartmentID                uint       `json:"department_id"`
	DepartmentTitleEn           string     `json:"department_title_en"`
	DepartmentTitleAr           string     `json:"department_title_ar"`
	Title                       *string    `json:"title"`
	Description                 *string    `json:"description"`
	CreatedAt                   time.Time  `json:"created_at"`
	UpdatedAt                   *time.Time `json:"updated_at"`
	CreatedUserID               uint       `json:"created_user_id"`
	CreatedUserFirstNameEn      string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string     `json:"created_user_last_name_ar"`
	CreatedUserInitial          string     `json:"created_user_initial"`
	CreatedUserColor            string     `json:"created_user_color"`
	CreatedUserEmail            string     `json:"created_user_email"`
	CreatedUserProfileImagePath string     `json:"created_user_profile_image_path"`
	UpdatedUserID               *uint      `json:"updated_user_id"`
	UpdatedUserFirstNameEn      *string    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      *string    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       *string    `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       *string    `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          string     `json:"updated_user_initial"`
	UpdatedUserColor            string     `json:"updated_user_color"`
	UpdatedUserEmail            string     `json:"updated_user_email"`
	UpdatedUserProfileImagePath string     `json:"updated_user_profile_image_path"`
}

type OrganizationUnitAttributeListQueryParamsDTO struct {
	Sort     string              `form:"sort_order"`
	SortBy   string              `form:"sort_by"`
	Keyword  string              `form:"keyword"`
	Page     int                 `form:"page"`
	PageSize int                 `form:"page_size"`
	Filters  map[string][]string `form:"-"`
}

type OrganizationUnitAttributeListResponseDTO struct {
	Items      []*OrganizationUnitAttributeListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO        `json:"pagination"`
}

type OrganizationUnitAttributeDetailDTO struct {
	ID            uint                 `json:"id"`
	DepartmentID  uint                 `json:"department_id"`
	Department    *OrganizationUnitDTO `json:"department"`
	Title         *string              `json:"title"`
	Description   *string              `json:"description"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     *time.Time           `json:"updated_at"`
	CreatedByUser UserMinimalDTO       `json:"created_by"`
	UpdatedByUser *UserMinimalDTO      `json:"updated_by"`
}

type OrganizationUnitAttributeFormDTO struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}
