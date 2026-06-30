package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type ProcessInternalDependencyListDTO struct {
	ID                          uint    `json:"id"`
	ProcessID                   uint    `json:"process_id"`
	SupportingUnitID            uint    `json:"supporting_unit_id"`
	SupportingUnitTitleEn       *string `json:"supporting_unit_title_en"`
	SupportingUnitTitleAr       *string `json:"supporting_unit_title_ar"`
	DependencyTierID            uint    `json:"dependency_tier_id"`
	DependencyTierTitleEn       *string `json:"dependency_tier_title_en"`
	DependencyTierTitleAr       *string `json:"dependency_tier_title_ar"`
	DependencyPercentageFrom    *int    `json:"dependency_percentage_from"`
	DependencyPercentageTo      *int    `json:"dependency_percentage_to"`
	RequiredSupport             *string `json:"required_support"`
	Notes                       *string `json:"notes"`
	CreatedBYID                 *uint   `json:"created_by_id"`
	CreatedUserFirstNameEn      string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string  `json:"created_user_last_name_ar"`
	CreatedUserInitial          string  `json:"created_user_initial"`
	CreatedUserColor            string  `json:"created_user_color"`
	CreatedUserEmail            string  `json:"created_user_email"`
	CreatedUserProfileImagePath string  `json:"created_user_profile_image_path"`
	UpdatedBYID                 *uint   `json:"updated_by_id"`
	UpdatedUserFirstNameEn      string  `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      string  `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       string  `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       string  `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          string  `json:"updated_user_initial"`
	UpdatedUserColor            string  `json:"updated_user_color"`
	UpdatedUserEmail            string  `json:"updated_user_email"`
	UpdatedUserProfileImagePath string  `json:"updated_user_profile_image_path"`
	CreatedAt                   string  `json:"created_at"`
	UpdatedAt                   string  `json:"updated_at"`
}

type ProcessInternalDependencyDetailDTO struct {
	ID                       uint                      `json:"id"`
	ProcessID                uint                      `json:"process_id"`
	SupportingUnitID         uint                      `json:"supporting_unit_id"`
	SupportingUnit           *OrganizationUnit1DTO     `json:"supporting_unit"`
	DependencyTierID         uint                      `json:"dependency_tier_id"`
	DependencyTier           *ProcessDependencyTierDTO `json:"dependency_tier"`
	DependencyPercentageFrom *int                      `json:"dependency_percentage_from"`
	DependencyPercentageTo   *int                      `json:"dependency_percentage_to"`
	RequiredSupport          *string                   `json:"required_support"`
	Notes                    *string                   `json:"notes"`
	CreatedByUser            *UserMinimalDTO           `json:"created_by"`
	UpdatedByUser            *UserMinimalDTO           `json:"updated_by"`
	CreatedAt                time.Time                 `json:"created_at"`
	UpdatedAt                *time.Time                `json:"updated_at"`
}

type ProcessInternalDependencyFormDTO struct {
	SupportingUnitID         uint    `json:"supporting_unit_id" binding:"required"`
	DependencyTierID         uint    `json:"dependency_tier_id" binding:"required"`
	DependencyPercentageFrom *int    `json:"dependency_percentage_from" binding:"omitempty"`
	DependencyPercentageTo   *int    `json:"dependency_percentage_to" binding:"omitempty"`
	RequiredSupport          *string `json:"required_support"`
	Notes                    *string `json:"notes"`
}

type ProcessInternalDependencyListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
	Export   bool
}

type ProcessInternalDependencyListResponseDTO struct {
	Items      []*ProcessInternalDependencyListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO        `json:"pagination"`
}
