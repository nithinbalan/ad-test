package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type IncidentComplianceRequirementListDTO struct {
	ComplianceRequirementID              uint   `json:"compliance_requirement_id"`
	ComplianceClauseTitle                string `json:"compliance_clause_title"`
	ComplianceRequirementCategoryTitleEn string `json:"compliance_requirement_category_title_en"`
	ComplianceRequirementCategoryTitleAr string `json:"compliance_requirement_category_title_ar"`
	ComplianceRequirementTitle           string `json:"compliance_requirement_title"`
	ComplianceStatusTitleEn              string `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr              string `json:"compliance_status_title_ar"`
	ComplianceStatusColorCode            string `json:"compliance_status_color_code"`
}

type IncidentSLAListDTO struct {
	ReferenceNumber                 string `json:"reference_number"`
	LifecycleStatusID               uint   `json:"lifecycle_status_id"`
	LifecycleStatusTitleEn          string `json:"lifecycle_status_title_en"`
	LifecycleStatusTitleAr          string `json:"lifecycle_status_title_ar"`
	LifecycleStatusColorCode        string `json:"lifecycle_status_color_code"`
	TitleEn                         string `json:"title_en"`
	TitleAr                         string `json:"title_ar"`
	Description                     string `json:"description"`
	Client                          string `json:"client"`
	SLACategoryTitleEn              string `json:"sla_category_title_en"`
	SLACategoryTitleAr              string `json:"sla_category_title_ar"`
	SeverityLevelTitleEn            string `json:"severity_level_title_en"`
	SeverityLevelTitleAr            string `json:"severity_level_title_ar"`
	ResponsibleUnitTitleEn          string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr          string `json:"responsible_unit_title_ar"`
	ResponsibleUserID               uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string `json:"responsible_user_initial"`
	ResponsibleUserColor            string `json:"responsible_user_color"`
	ResponsibleUserEmail            string `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string `json:"responsible_user_profile_image_path"`
	ComplianceStatusTitleEn         string `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr         string `json:"compliance_status_title_ar"`
	ComplianceStatusColorCode       string `json:"compliance_status_color_code"`
}

type IncidentAssetListDTO struct {
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

type IncidentAssociatedItemListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	QuickFilter []string
}

type IncidentAssociatedItemComplianceRequirementListResponseDTO struct {
	Items      []*IncidentComplianceRequirementListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO            `json:"pagination"`
}

type IncidentAssociatedItemSLAListResponseDTO struct {
	Items      []*IncidentSLAListDTO        `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type IncidentAssociatedItemAssetListResponseDTO struct {
	Items      []*IncidentAssetListDTO      `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}
