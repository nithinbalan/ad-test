package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"time"
)

type AssetRegisterListDTO struct {
	ID                              uint       `json:"id"`
	ReferenceNumber                 string     `json:"reference_number"`
	TitleEn                         string     `json:"title_en"`
	TitleAr                         string     `json:"title_ar"`
	AssetTypeTitleEn                *string    `json:"asset_type_title_en"`
	AssetTypeTitleAr                *string    `json:"asset_type_title_ar"`
	AssetCategoryTitleEn            *string    `json:"asset_category_title_en"`
	AssetCategoryTitleAr            *string    `json:"asset_category_title_ar"`
	AssetStatusTitleEn              string     `json:"asset_status_title_en"`
	AssetStatusTitleAr              string     `json:"asset_status_title_ar"`
	AssetStatusColorCode            string     `json:"asset_status_color_code"`
	ResponsibleUnitTitleEn          string     `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr          string     `json:"responsible_unit_title_ar"`
	ResponsibleUserID               uint       `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string     `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string     `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string     `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string     `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string     `json:"responsible_user_initial"`
	ResponsibleUserColor            string     `json:"responsible_user_color"`
	ResponsibleUserEmail            string     `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string     `json:"responsible_user_profile_image_path"`
	Location                        *string    `json:"location"`
	Description                     *string    `json:"description"`
	AcquisitionDate                 string     `json:"acquisition_date"`
	CreatedUserID                   uint       `json:"created_user_id"`
	CreatedUserFirstNameEn          string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr          string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn           string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr           string     `json:"created_user_last_name_ar"`
	CreatedUserInitial              string     `json:"created_user_initial"`
	CreatedUserColor                string     `json:"created_user_color"`
	CreatedUserEmail                string     `json:"created_user_email"`
	CreatedUserProfileImagePath     string     `json:"created_user_profile_image_path"`
	CreatedAt                       time.Time  `json:"created_at"`
	UpdatedUserID                   *uint      `json:"updated_user_id"`
	UpdatedUserFirstNameEn          *string    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr          *string    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn           *string    `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr           *string    `json:"updated_user_last_name_ar"`
	UpdatedUserInitial              string     `json:"updated_user_initial"`
	UpdatedUserColor                string     `json:"updated_user_color"`
	UpdatedUserEmail                string     `json:"updated_user_email"`
	UpdatedUserProfileImagePath     string     `json:"updated_user_profile_image_path"`
	UpdatedAt                       *time.Time `json:"updated_at"`
}

type AssetRegisterListQueryParamsDTO struct {
	Sort                string
	SortBy              string
	Keyword             string
	Status              string
	Page                int
	PageSize            int
	Filters             map[string][]string
	QuickFilter         []string
	AcquisitionDateFrom string
	AcquisitionDateTo   string
	Export              bool
}

type AssetRegisterListResponseDTO struct {
	Items      []*AssetRegisterListDTO      `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type AssetRegisterFormDTO struct {
	TitleEn           string  `json:"title_en" binding:"required"`
	TitleAr           string  `json:"title_ar" binding:"required"`
	AssetTypeID       uint    `json:"asset_type_id" binding:"required"`
	AssetCategoryID   uint    `json:"asset_category_id" binding:"required"`
	AssetStatusID     uint    `json:"asset_status_id" binding:"required"`
	ResponsibleUserID uint    `json:"responsible_user_id" binding:"required"`
	Location          *string `json:"location" binding:"omitempty"`
	Description       *string `json:"description" binding:"omitempty"`
	AcquisitionDate   *string `json:"acquisition_date" binding:"omitempty,datetime=2006-01-02"`
	ControlIds        []*uint `json:"control_ids"`
}

type AssetRegisterParsedDates struct {
	AcquisitionDate *time.Time
}

type AssetRegisterDetailDTO struct {
	ID                    uint                             `json:"id"`
	ReferenceNumber       string                           `json:"reference_number"`
	TitleEn               string                           `json:"title_en"`
	TitleAr               string                           `json:"title_ar"`
	AssetType             *AssetTypeDTO                    `json:"asset_type"`
	AssetCategory         *AssetCategoryDTO                `json:"asset_category"`
	AssetStatus           AssetStatusDTO                   `json:"asset_status"`
	ResponsibleUser       organizationDTOs.UserMinimalDTO  `json:"responsible_user"`
	ResponsibleUnit       settingsDTOs.OrganizationUnitDTO `json:"responsible_unit"`
	Location              *string                          `json:"location"`
	Description           *string                          `json:"description"`
	AcquisitionDateString *string                          `json:"acquisition_date"`
	Controls              []*AssetControlDTO               `json:"controls"`
	CreatedAt             time.Time                        `json:"created_at"`
	UpdatedAt             *time.Time                       `json:"updated_at"`
	CreatedByUser         organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser         *organizationDTOs.UserMinimalDTO `json:"updated_by"`
}

type AssetControlDTO struct {
	ID        *uint                             `json:"id"`
	ProcessID uint                              `json:"process_id"`
	ControlID uint                              `json:"control_id"`
	Control   organizationDTOs.ControlDetailDTO `json:"control"`
}
