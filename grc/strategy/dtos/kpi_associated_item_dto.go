package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
)

type KPIControlListDTO struct {
	ControlID                            uint    `json:"control_id"`
	ControlReferenceNumber               string  `json:"control_reference_number"`
	ControlCode                          *string `json:"control_code"`
	ControlLastTestDate                  *string `json:"control_last_test_date"`
	ControlCategoryTitleEn               string  `json:"control_category_title_en"`
	ControlCategoryTitleAr               string  `json:"control_category_title_ar"`
	ControlTypeTitleEn                   string  `json:"control_type_title_en"`
	ControlTypeTitleAr                   string  `json:"control_type_title_ar"`
	ControlTitleEn                       string  `json:"control_title_en"`
	ControlTitleAr                       string  `json:"control_title_ar"`
	ControlDescription                   *string `json:"control_description"`
	ControlEfficiencyStatusTitleEn       string  `json:"control_efficiency_status_title_en"`
	ControlEfficiencyStatusTitleAr       string  `json:"control_efficiency_status_title_ar"`
	ControlEfficiencyStatusType          string  `json:"control_efficiency_status_type"`
	ControlEfficiencyStatusColorCode     string  `json:"control_efficiency_status_color_code"`
	ControlImplementationStatusTitleEn   string  `json:"control_implementation_status_title_en"`
	ControlImplementationStatusTitleAr   string  `json:"control_implementation_status_title_ar"`
	ControlImplementationStatusType      string  `json:"control_implementation_status_type"`
	ControlImplementationStatusColorCode string  `json:"control_implementation_status_color_code"`
	ControlOwnerUnitTitleEn              *string `json:"control_owner_unit_title_en"`
	ControlOwnerUnitTitleAr              *string `json:"control_owner_unit_title_ar"`
	ControlOwnerID                       *uint   `json:"control_owner_id"`
	ControlOwnerFirstNameEn              *string `json:"control_owner_first_name_en"`
	ControlOwnerFirstNameAr              *string `json:"control_owner_first_name_ar"`
	ControlOwnerLastNameEn               *string `json:"control_owner_last_name_en"`
	ControlOwnerLastNameAr               *string `json:"control_owner_last_name_ar"`
	ControlOwnerInitial                  string  `json:"control_owner_initial"`
	ControlOwnerColor                    string  `json:"control_owner_color"`
	ControlOwnerEmail                    string  `json:"control_owner_email"`
	ControlOwnerProfileImagePath         string  `json:"control_owner_profile_image_path"`
}

type KPIAssociatedItemControlListResponseDTO struct {
	Items      []*KPIControlListDTO         `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type KPIControlMapDTO struct {
	ControlIDs []uint `json:"control_ids" binding:"required"`
}

type KPIProcessListDTO struct {
	ProcessID                uint    `json:"process_id"`
	ProcessGroupTitleEn      string  `json:"process_group_title_en"`
	ProcessGroupTitleAr      string  `json:"process_group_title_ar"`
	OwnerUnitTitleEn         *string `json:"owner_unit_title_en"`
	OwnerUnitTitleAr         *string `json:"owner_unit_title_ar"`
	OwnerID                  *uint   `json:"owner_id"`
	OwnerFirstNameEn         *string `json:"owner_first_name_en"`
	OwnerFirstNameAr         *string `json:"owner_first_name_ar"`
	OwnerLastNameEn          *string `json:"owner_last_name_en"`
	OwnerLastNameAr          *string `json:"owner_last_name_ar"`
	OwnerInitial             string  `json:"owner_initial"`
	OwnerColor               string  `json:"owner_color"`
	OwnerEmail               string  `json:"owner_email"`
	OwnerProfileImagePath    string  `json:"owner_profile_image_path"`
	LifecycleStatusTitleEn   string  `json:"lifecycle_status_title_en"`
	LifecycleStatusTitleAr   string  `json:"lifecycle_status_title_ar"`
	LifecycleStatusColorCode string  `json:"lifecycle_status_color_code"`
	ReferenceNumber          string  `json:"reference_number"`
	TitleEn                  string  `json:"title_en"`
	TitleAr                  string  `json:"title_ar"`
	Description              *string `json:"description"`
}

type KPIAssociatedItemProcessListResponseDTO struct {
	Items      []*KPIProcessListDTO         `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type KPIProcessMapDTO struct {
	ProcessIDs []uint `json:"process_ids" binding:"required"`
}

type KPIAssetListDTO struct {
	AssetRegisterID                 uint    `json:"asset_register_id"`
	ReferenceNumber                 string  `json:"reference_number"`
	TitleEn                         string  `json:"title_en"`
	TitleAr                         string  `json:"title_ar"`
	AssetCategoryTitleEn            *string `json:"asset_category_title_en"`
	AssetCategoryTitleAr            *string `json:"asset_category_title_ar"`
	AssetStatusTitleEn              string  `json:"asset_status_title_en"`
	AssetStatusTitleAr              string  `json:"asset_status_title_ar"`
	AssetStatusColorCode            string  `json:"asset_status_color_code"`
	ResponsibleUnitTitleEn          string  `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr          string  `json:"responsible_unit_title_ar"`
	ResponsibleUserID               uint    `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string  `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string  `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string  `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string  `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string  `json:"responsible_user_initial"`
	ResponsibleUserColor            string  `json:"responsible_user_color"`
	ResponsibleUserEmail            string  `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string  `json:"responsible_user_profile_image_path"`
	Location                        *string `json:"location"`
	Description                     *string `json:"description"`
	AcquisitionDate                 string  `json:"acquisition_date"`
}

type KPIAssociatedItemAssetListResponseDTO struct {
	Items      []*KPIAssetListDTO           `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type KPIAssetMapDTO struct {
	AssetIDs []uint `json:"asset_ids" binding:"required"`
}

type KPIComplianceListDTO struct {
	ComplianceRegisterID                  uint   `json:"compliance_register_id"`
	ComplianceCategoryTitleEn             string `json:"compliance_category_title_en"`
	ComplianceCategoryTitleAr             string `json:"compliance_category_title_ar"`
	ComplianceStatusTitleEn               string `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr               string `json:"compliance_status_title_ar"`
	ComplianceStatusColorCode             string `json:"compliance_status_color_code"`
	ResponsibleUnitTitleEn                string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr                string `json:"responsible_unit_title_ar"`
	ResponsibleUserID                     uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn            string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr            string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn             string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr             string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial                string `json:"responsible_user_initial"`
	ResponsibleUserColor                  string `json:"responsible_user_color"`
	ResponsibleUserEmail                  string `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath       string `json:"responsible_user_profile_image_path"`
	ReviewFrequencyTitleEn                string `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr                string `json:"review_frequency_title_ar"`
	ReviewResponsibleUserID               uint   `json:"review_responsible_user_id"`
	ReviewResponsibleUserFirstNameEn      string `json:"review_responsible_user_first_name_en"`
	ReviewResponsibleUserFirstNameAr      string `json:"review_responsible_user_first_name_ar"`
	ReviewResponsibleUserLastNameEn       string `json:"review_responsible_user_last_name_en"`
	ReviewResponsibleUserLastNameAr       string `json:"review_responsible_user_last_name_ar"`
	ReviewResponsibleUserInitial          string `json:"review_responsible_user_initial"`
	ReviewResponsibleUserColor            string `json:"review_responsible_user_color"`
	ReviewResponsibleUserEmail            string `json:"review_responsible_user_email"`
	ReviewResponsibleUserProfileImagePath string `json:"review_responsible_user_profile_image_path"`
	SeverityLevelTitleEn                  string `json:"severity_level_title_en"`
	SeverityLevelTitleAr                  string `json:"severity_level_title_ar"`
	SeverityLevelColorCode                string `json:"severity_level_color_code"`
	ComplianceRegisterReferenceNumber     string `json:"compliance_register_reference_number"`
	ComplianceRegisterStandardCode        string `json:"compliance_register_standard_code"`
	ComplianceRegisterTitleEn             string `json:"compliance_register_title_en"`
	ComplianceRegisterTitleAr             string `json:"compliance_register_title_ar"`
	ComplianceRegisterDescriptionEn       string `json:"compliance_register_description_en"`
	ComplianceRegisterDescriptionAr       string `json:"compliance_register_description_ar"`
	ComplianceRegisterVersion             string `json:"compliance_register_version"`
	Deadline                              string `json:"deadline"`
	IssueDate                             string `json:"issue_date"`
	ExpiryDate                            string `json:"expiry_date"`
	NextReviewDate                        string `json:"next_review_date"`
	LastReviewedDate                      string `json:"last_reviewed_date"`
}

type KPIAssociatedItemComplianceListResponseDTO struct {
	Items      []*KPIComplianceListDTO      `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type KPIComplianceMapDTO struct {
	ComplianceRegisterIDs []uint `json:"compliance_register_ids" binding:"required"`
}

type KPIRiskMapDTO struct {
	RiskIDs []uint `json:"risk_ids" binding:"required"`
}

type KPIThirdPartyListDTO struct {
	ThirdPartyRegisterID            uint    `json:"third_party_register_id"`
	ReferenceNumber                 string  `json:"reference_number"`
	Title                           string  `json:"title"`
	ThirdPartyTypeID                uint    `json:"third_party_type_id"`
	ThirdPartyTypeTitleEn           string  `json:"third_party_type_title_en"`
	ThirdPartyTypeTitleAr           string  `json:"third_party_type_title_ar"`
	CountryID                       uint    `json:"country_id"`
	CountryTitleEn                  string  `json:"country_title_en"`
	CountryTitleAr                  string  `json:"country_title_ar"`
	Address                         *string `json:"address"`
	Website                         *string `json:"website"`
	OrganizationUnitID              *uint   `json:"organization_unit_id"`
	OrganizationUnitTitleEn         string  `json:"organization_unit_title_en"`
	OrganizationUnitTitleAr         string  `json:"organization_unit_title_ar"`
	RepresentativeName              *string `json:"representative_name"`
	RepresentativeRole              *string `json:"representative_role"`
	Email                           *string `json:"email"`
	PhoneNumber                     *string `json:"phone_number"`
	ResponsibleUserID               *uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string  `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string  `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string  `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string  `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string  `json:"responsible_user_initial"`
	ResponsibleUserColor            string  `json:"responsible_user_color"`
	ResponsibleUserEmail            string  `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string  `json:"responsible_user_profile_image_path"`
	OnboardingDate                  *string `json:"onboarding_date"`
	ServiceStartDate                *string `json:"service_start_date"`
	ServiceEndDate                  *string `json:"service_end_date"`
	SeverityLevelID                 *uint   `json:"severity_level_id"`
	SeverityLevelTitleEn            string  `json:"severity_level_title_en"`
	SeverityLevelTitleAr            string  `json:"severity_level_title_ar"`
	SeverityLevelColorCode          string  `json:"severity_level_color_code"`
}

type KPIAssociatedItemThirdPartyListResponseDTO struct {
	Items      []*KPIThirdPartyListDTO      `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type KPIThirdPartyMapDTO struct {
	ThirdPartyRegisterIDs []uint `json:"third_party_register_ids" binding:"required"`
}

type KPIPolicyListDTO struct {
	ID                   uint    `json:"id"`
	DocumentID           uint    `json:"document_id"`
	DocumentTitleEn      string  `json:"document_title_en"`
	DocumentTitleAr      string  `json:"document_title_ar"`
	DocumentURL          string  `json:"document_url"`
	DocumentName         string  `json:"document_name"`
	DocumentSize         int     `json:"document_size"`
	DocumentToken        string  `json:"document_token"`
	DocumentOriginalName *string `json:"document_original_name"`
}

type KPIAssociatedItemPolicyListResponseDTO struct {
	Items      []*KPIPolicyListDTO          `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type KPIPolicyMapDTO struct {
	DocumentIDs []uint `json:"document_ids" binding:"required"`
}

type KPIAssociatedItemListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Value       int
	QuickFilter []string
}
