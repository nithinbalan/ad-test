package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type AssetMaintenanceListDTO struct {
	ID                          uint       `json:"id"`
	AssetRegisterID             uint       `json:"asset_register_id"`
	Type                        string     `json:"type"`
	PerformedDate               string     `json:"performed_date"`
	Title                       string     `json:"title"`
	Description                 *string    `json:"description"`
	PerformedByID               uint       `json:"performed_by_id"`
	PerformedByFirstNameEn      string     `json:"performed_by_first_name_en"`
	PerformedByFirstNameAr      string     `json:"performed_by_first_name_ar"`
	PerformedByLastNameEn       string     `json:"performed_by_last_name_en"`
	PerformedByLastNameAr       string     `json:"performed_by_last_name_ar"`
	PerformedByInitial          string     `json:"performed_by_initial"`
	PerformedByColor            string     `json:"performed_by_color"`
	PerformedByEmail            string     `json:"performed_by_email"`
	PerformedByProfileImagePath string     `json:"performed_by_profile_image_path"`
	Outcome                     *string    `json:"outcome"`
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

type AssetMaintenanceListQueryParamsDTO struct {
	Sort              string
	SortBy            string
	Keyword           string
	Page              int
	PageSize          int
	Filters           map[string][]string
	QuickFilter       []string
	PerformedDateFrom string
	PerformedDateTo   string
	Export            bool
}

type AssetMaintenanceListResponseDTO struct {
	Items      []*AssetMaintenanceListDTO   `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type AssetMaintenanceFormDTO struct {
	AssetRegisterID uint    `json:"asset_register_id" binding:"required"`
	Type            string  `json:"type" binding:"required"`
	PerformedDate   string  `json:"performed_date" binding:"required,datetime=2006-01-02"`
	Title           string  `json:"title" binding:"required"`
	Description     *string `json:"description" binding:"omitempty"`
	PerformedByID   uint    `json:"performed_by_id" binding:"required"`
	Outcome         *string `json:"outcome" binding:"omitempty"`
}

type AssetMaintenanceParsedDates struct {
	PerformedDate *time.Time
}

type AssetMaintenanceDetailDTO struct {
	ID              uint                             `json:"id"`
	AssetRegisterID uint                             `json:"asset_register_id"`
	Type            string                           `json:"type"`
	PerformedDate   string                           `json:"performed_date"`
	Title           string                           `json:"title"`
	Description     *string                          `json:"description"`
	PerformedBy     organizationDTOs.UserMinimalDTO  `json:"performed_by"`
	Outcome         *string                          `json:"outcome"`
	CreatedAt       time.Time                        `json:"created_at"`
	UpdatedAt       *time.Time                       `json:"updated_at"`
	CreatedBy       organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedBy       *organizationDTOs.UserMinimalDTO `json:"updated_by"`
}
