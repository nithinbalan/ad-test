package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type ProcessExternalDependencyListDTO struct {
	ID                          uint    `json:"id"`
	ProcessID                   uint    `json:"process_id"`
	EntityName                  string  `json:"entity_name"`
	EntityTypeID                *uint   `json:"entity_type_id"`
	EntityTypeTitleEn           *string `json:"entity_type_title_en"`
	EntityTypeTitleAr           *string `json:"entity_type_title_ar"`
	DependencyTierID            *uint   `json:"dependency_tier_id"`
	DependencyTierTitleEn       *string `json:"dependency_tier_title_en"`
	DependencyTierTitleAr       *string `json:"dependency_tier_title_ar"`
	NeededSupportDescription    *string `json:"needed_support_description"`
	VendorRTO                   *int    `json:"vendor_rto"`
	VendorRTOUnitID             *uint   `json:"vendor_rto_unit_id"`
	VendorRTOUnitTitleEn        *string `json:"vendor_rto_unit_title_en"`
	VendorRTOUnitTitleAr        *string `json:"vendor_rto_unit_title_ar"`
	ContactNumber               *string `json:"contact_number"`
	ContactEmail                *string `json:"contact_email"`
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

type ProcessExternalDependencyDetailDTO struct {
	ID                       uint                                    `json:"id"`
	ProcessID                uint                                    `json:"process_id"`
	EntityName               string                                  `json:"entity_name"`
	EntityTypeID             *uint                                   `json:"entity_type_id"`
	EntityType               *ProcessExternalDependencyEntityTypeDTO `json:"entity_type"`
	NeededSupportDescription *string                                 `json:"needed_support_description"`
	DependencyTierID         *uint                                   `json:"dependency_tier_id"`
	DependencyTier           *ProcessDependencyTierDTO               `json:"dependency_tier"`
	SLARegisterID            *uint                                   `json:"sla_register_id"`
	VendorRTO                *int                                    `json:"vendor_rto"`
	VendorRTOUnit            *BusinessContinuityDurationUnitDTO      `json:"vendor_rto_unit"`
	InternalStakeholder      *string                                 `json:"internal_stakeholder"`
	ContactNumber            *string                                 `json:"contact_number"`
	ContactEmail             *string                                 `json:"contact_email"`
	CreatedByUser            *UserMinimalDTO                         `json:"created_by"`
	UpdatedByUser            *UserMinimalDTO                         `json:"updated_by"`
	CreatedAt                time.Time                               `json:"created_at"`
	UpdatedAt                *time.Time                              `json:"updated_at"`
}

type ProcessExternalDependencyFormDTO struct {
	EntityName               string  `json:"entity_name" binding:"required"`
	EntityTypeID             *uint   `json:"entity_type_id" binding:"omitempty"`
	NeededSupportDescription *string `json:"needed_support_description"`
	DependencyTierID         *uint   `json:"dependency_tier_id" binding:"omitempty"`
	SLARegisterID            *uint   `json:"sla_register_id" binding:"omitempty"`
	InternalStakeholder      *string `json:"internal_stakeholder"`
	ContactNumber            *string `json:"contact_number"`
	ContactEmail             *string `json:"contact_email"`
}

type ProcessExternalDependencyListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
	Export   bool
}

type ProcessExternalDependencyListResponseDTO struct {
	Items      []*ProcessExternalDependencyListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO        `json:"pagination"`
}
