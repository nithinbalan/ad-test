package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"time"
)

type AssetRegisterListDTO struct {
	ID                                         uint       `json:"id"`
	ReferenceNumber                            string     `json:"reference_number"`
	TitleEn                                    string     `json:"title_en"`
	TitleAr                                    string     `json:"title_ar"`
	AssetCategoryTitleEn                       *string    `json:"asset_category_title_en"`
	AssetCategoryTitleAr                       *string    `json:"asset_category_title_ar"`
	AssetSubCategoryTitleEn                    *string    `json:"asset_sub_category_title_en"`
	AssetSubCategoryTitleAr                    *string    `json:"asset_sub_category_title_ar"`
	AssetStatusTitleEn                         string     `json:"asset_status_title_en"`
	AssetStatusTitleAr                         string     `json:"asset_status_title_ar"`
	AssetStatusColorCode                       string     `json:"asset_status_color_code"`
	ResponsibleUnitTitleEn                     string     `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr                     string     `json:"responsible_unit_title_ar"`
	ResponsibleUserID                          uint       `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn                 string     `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr                 string     `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn                  string     `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr                  string     `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial                     string     `json:"responsible_user_initial"`
	ResponsibleUserColor                       string     `json:"responsible_user_color"`
	ResponsibleUserEmail                       string     `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath            string     `json:"responsible_user_profile_image_path"`
	Location                                   *string    `json:"location"`
	Description                                *string    `json:"description"`
	AcquisitionDate                            string     `json:"acquisition_date"`
	AssignmentDate                             *string    `json:"assignment_date"`
	WarrantyExpiryDate                         *string    `json:"warranty_expiry_date"`
	Manufacturer                               *string    `json:"manufacturer"`
	Model                                      *string    `json:"model"`
	SerialNumber                               *string    `json:"serial_number"`
	TotalScore                                 *float64   `json:"total_score"`
	AssetCriticalityMatrixRatingScaleTitleEn   string     `json:"asset_criticality_matrix_rating_scale_title_en"`
	AssetCriticalityMatrixRatingScaleTitleAr   string     `json:"asset_criticality_matrix_rating_scale_title_ar"`
	AssetCriticalityMatrixRatingScaleColorCode string     `json:"asset_criticality_matrix_rating_scale_color_code"`
	Rto                                        *int       `json:"rto"`
	RtoUnitId                                  *uint      `json:"rto_unit_id"`
	Mtpd                                       *int       `json:"mtpd"`
	MtpdUnitId                                 *uint      `json:"mtpd_unit_id"`
	Rpo                                        *int       `json:"rpo"`
	RpoUnitId                                  *uint      `json:"rpo_unit_id"`
	TechnicalOwnerID                           *uint      `json:"technical_owner_id"`
	TechnicalOwnerFirstNameEn                  *string    `json:"technical_owner_first_name_en"`
	TechnicalOwnerFirstNameAr                  *string    `json:"technical_owner_first_name_ar"`
	TechnicalOwnerLastNameEn                   *string    `json:"technical_owner_last_name_en"`
	TechnicalOwnerLastNameAr                   *string    `json:"technical_owner_last_name_ar"`
	AdministrativeOwnerID                      *uint      `json:"administrative_owner_id"`
	AdministrativeOwnerFirstNameEn             *string    `json:"administrative_owner_first_name_en"`
	AdministrativeOwnerFirstNameAr             *string    `json:"administrative_owner_first_name_ar"`
	AdministrativeOwnerLastNameEn              *string    `json:"administrative_owner_last_name_en"`
	AdministrativeOwnerLastNameAr              *string    `json:"administrative_owner_last_name_ar"`
	AssignedUserID                             *uint      `json:"assigned_user_id"`
	AssignedUserFirstNameEn                    *string    `json:"assigned_user_first_name_en"`
	AssignedUserFirstNameAr                    *string    `json:"assigned_user_first_name_ar"`
	AssignedUserLastNameEn                     *string    `json:"assigned_user_last_name_en"`
	AssignedUserLastNameAr                     *string    `json:"assigned_user_last_name_ar"`
	Hostname                                   *string    `json:"hostname"`
	IPAddress                                  *string    `json:"ip_address"`
	MACAddress                                 *string    `json:"mac_address"`
	FQDN                                       *string    `json:"fqdn"`
	Region                                     *string    `json:"region"`
	OperatingSystem                            *string    `json:"operating_system"`
	DeviceCategoryTitleEn                      *string    `json:"device_category_title_en"`
	DeviceCategoryTitleAr                      *string    `json:"device_category_title_ar"`
	CriticalityClassificationTitleEn           *string    `json:"criticality_classification_title_en"`
	CriticalityClassificationTitleAr           *string    `json:"criticality_classification_title_ar"`
	ConfidentialityClassificationTitleEn       *string    `json:"confidentiality_classification_title_en"`
	ConfidentialityClassificationTitleAr       *string    `json:"confidentiality_classification_title_ar"`
	AssetEnvironmentTitleEn                    *string    `json:"asset_environment_title_en"`
	AssetEnvironmentTitleAr                    *string    `json:"asset_environment_title_ar"`
	InventoryTime                              *string    `json:"inventory_time"`
	DiscoveryTime                              *string    `json:"discovery_time"`
	LastResponseTime                           *string    `json:"last_response_time"`
	CreatedUserID                              uint       `json:"created_user_id"`
	CreatedUserFirstNameEn                     string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr                     string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn                      string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr                      string     `json:"created_user_last_name_ar"`
	CreatedUserInitial                         string     `json:"created_user_initial"`
	CreatedUserColor                           string     `json:"created_user_color"`
	CreatedUserEmail                           string     `json:"created_user_email"`
	CreatedUserProfileImagePath                string     `json:"created_user_profile_image_path"`
	CreatedAt                                  time.Time  `json:"created_at"`
	UpdatedUserID                              *uint      `json:"updated_user_id"`
	UpdatedUserFirstNameEn                     *string    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr                     *string    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn                      *string    `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr                      *string    `json:"updated_user_last_name_ar"`
	UpdatedUserInitial                         string     `json:"updated_user_initial"`
	UpdatedUserColor                           string     `json:"updated_user_color"`
	UpdatedUserEmail                           string     `json:"updated_user_email"`
	UpdatedUserProfileImagePath                string     `json:"updated_user_profile_image_path"`
	UpdatedAt                                  *time.Time `json:"updated_at"`
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
	IsArchived          *bool
}

type AssetRegisterListResponseDTO struct {
	Items      []*AssetRegisterListDTO      `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type AssetRegisterFormDTO struct {
	TitleEn            string                        `json:"title_en" binding:"required"`
	TitleAr            string                        `json:"title_ar" binding:"required"`
	AssetCategoryID    uint                          `json:"asset_category_id" binding:"required"`
	AssetSubCategoryID uint                          `json:"asset_sub_category_id" binding:"required"`
	AssetStatusID      uint                          `json:"asset_status_id" binding:"required"`
	ResponsibleUserID  uint                          `json:"responsible_user_id" binding:"required"`
	Location           *string                       `json:"location" binding:"omitempty"`
	Description        *string                       `json:"description" binding:"omitempty"`
	AcquisitionDate    *string                       `json:"acquisition_date" binding:"omitempty,datetime=2006-01-02"`
	AssignmentDate     *string                       `json:"assignment_date" binding:"omitempty,datetime=2006-01-02"`
	WarrantyExpiryDate *string                       `json:"warranty_expiry_date" binding:"omitempty,datetime=2006-01-02"`
	Manufacturer       *string                       `json:"manufacturer" binding:"omitempty"`
	Model              *string                       `json:"model" binding:"omitempty"`
	SerialNumber       *string                       `json:"serial_number" binding:"omitempty"`
	Documents          []*commonDTOs.DocumentFormDTO `json:"documents"`

	TechnicalOwnerID                *uint   `json:"technical_owner_id" binding:"omitempty"`
	AdministrativeOwnerID           *uint   `json:"administrative_owner_id" binding:"omitempty"`
	Hostname                        *string `json:"hostname" binding:"omitempty"`
	IPAddress                       *string `json:"ip_address" binding:"omitempty"`
	MACAddress                      *string `json:"mac_address" binding:"omitempty"`
	FQDN                            *string `json:"fqdn" binding:"omitempty"`
	DeviceCategoryID                *uint   `json:"device_category_id" binding:"omitempty"`
	CriticalityClassificationID     *uint   `json:"criticality_classification_id" binding:"omitempty"`
	ConfidentialityClassificationID *uint   `json:"confidentiality_classification_id" binding:"omitempty"`
	Region                          *string `json:"region" binding:"omitempty"`
	AssignedUserID                  *uint   `json:"assigned_user_id" binding:"omitempty"`
	OperatingSystem                 *string `json:"operating_system" binding:"omitempty"`
	AssetEnvironmentID              *uint   `json:"asset_environment_id" binding:"omitempty"`
	InventoryTime                   *string `json:"inventory_time" binding:"omitempty,datetime=2006-01-02T15:04:05"`
	DiscoveryTime                   *string `json:"discovery_time" binding:"omitempty,datetime=2006-01-02T15:04:05"`
	LastResponseTime                *string `json:"last_response_time" binding:"omitempty,datetime=2006-01-02T15:04:05"`
}

type AssetRegisterParsedDates struct {
	AcquisitionDate    *time.Time
	AssignmentDate     *time.Time
	WarrantyExpiryDate *time.Time
	InventoryTime      *time.Time
	DiscoveryTime      *time.Time
	LastResponseTime   *time.Time
}

type AssetRegisterDetailDTO struct {
	ID                                uint                                  `json:"id"`
	ReferenceNumber                   string                                `json:"reference_number"`
	TitleEn                           string                                `json:"title_en"`
	TitleAr                           string                                `json:"title_ar"`
	AssetCategory                     *AssetCategoryDTO                     `json:"asset_category"`
	AssetSubCategory                  *AssetSubCategoryDTO                  `json:"asset_sub_category"`
	AssetStatus                       AssetStatusDTO                        `json:"asset_status"`
	ResponsibleUser                   organizationDTOs.UserMinimalDTO       `json:"responsible_user"`
	ResponsibleUnit                   settingsDTOs.OrganizationUnitDTO      `json:"responsible_unit"`
	Location                          *string                               `json:"location"`
	Description                       *string                               `json:"description"`
	AcquisitionDateString             *string                               `json:"acquisition_date"`
	Controls                          []*AssetControlDTO                    `json:"controls"`
	AssignmentDateString              *string                               `json:"assignment_date"`
	WarrantyExpiryDateString          *string                               `json:"warranty_expiry_date"`
	Manufacturer                      *string                               `json:"manufacturer"`
	Model                             *string                               `json:"model"`
	SerialNumber                      *string                               `json:"serial_number"`
	AssetCriticalityMatrix            *AssetCriticalityMatrixDetailDTO      `json:"asset_criticality_matrix"`
	AssetCriticalityMatrixRatingScale *AssetCriticalityMatrixRatingScaleDTO `json:"asset_criticality_matrix_rating_scale"`
	TotalScore                        *float64                              `json:"total_score"`
	Documents                         []*AssetDocumentDTO                   `json:"documents"`
	CreatedAt                         time.Time                             `json:"created_at"`
	UpdatedAt                         *time.Time                            `json:"updated_at"`
	CreatedByUser                     organizationDTOs.UserMinimalDTO       `json:"created_by"`
	UpdatedByUser                     *organizationDTOs.UserMinimalDTO      `json:"updated_by"`
	Kris                              []*AssetKRIDTO                        `json:"kris"`

	Hostname                      *string                                `json:"hostname"`
	IPAddress                     *string                                `json:"ip_address"`
	MACAddress                    *string                                `json:"mac_address"`
	FQDN                          *string                                `json:"fqdn"`
	Region                        *string                                `json:"region"`
	OperatingSystem               *string                                `json:"operating_system"`
	IsArchived                    bool                                   `json:"is_archived"`
	InventoryTimeString           *string                                `json:"inventory_time"`
	DiscoveryTimeString           *string                                `json:"discovery_time"`
	LastResponseTimeString        *string                                `json:"last_response_time"`
	TechnicalOwner                *organizationDTOs.UserMinimalDTO       `json:"technical_owner"`
	AdministrativeOwner           *organizationDTOs.UserMinimalDTO       `json:"administrative_owner"`
	AssignedUser                  *organizationDTOs.UserMinimalDTO       `json:"assigned_user"`
	DeviceCategory                *AssetDeviceCategoryDTO                `json:"device_category"`
	CriticalityClassification     *AssetCriticalityMatrixFactorRatingDTO `json:"criticality_classification"`
	ConfidentialityClassification *ConfidentialityClassificationDTO      `json:"confidentiality_classification"`
	AssetEnvironment              *AssetEnvironmentDTO                   `json:"asset_environment"`
}

type AssetKRIDTO struct {
	ID              *uint `json:"id"`
	AssetRegisterID uint  `json:"risk_register_id"`
	KriID           uint  `json:"kri_id"`
}

type AssetDocumentDTO struct {
	ID           *uint   `json:"id"`
	URL          *string `json:"url"`
	Name         string  `json:"name"`
	OriginalName string  `json:"original_name"`
	Token        string  `json:"token"`
	Type         string  `json:"type"`
	Size         int     `json:"size"`
}
type AssetControlDTO struct {
	ID        *uint                             `json:"id"`
	ProcessID uint                              `json:"process_id"`
	ControlID uint                              `json:"control_id"`
	Control   organizationDTOs.ControlDetailDTO `json:"control"`
}
