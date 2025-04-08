package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"

	"time"
)

type AssetControlAnalysisListDTO struct {
	ID                                   uint       `json:"id"`
	ControlID                            uint       `json:"control_id"`
	ControlCode                          string     `json:"control_code"`
	ControlReferenceNumber               string     `json:"control_reference_number"`
	ControlDescription                   string     `json:"control_description"`
	ControlTitleEn                       string     `json:"control_title_en"`
	ControlTitleAr                       string     `json:"control_title_ar"`
	ControlTypeTitleEn                   string     `json:"control_type_title_en"`
	ControlTypeTitleAr                   string     `json:"control_type_title_ar"`
	ControlCategoryTitleEn               string     `json:"control_category_title_en"`
	ControlCategoryTitleAr               string     `json:"control_category_title_ar"`
	OwnerUnitTitleEn                     string     `json:"owner_unit_title_en"`
	OwnerUnitTitleAr                     string     `json:"owner_unit_title_ar"`
	OwnerID                              uint       `json:"owner_id"`
	OwnerFirstNameEn                     string     `json:"owner_first_name_en"`
	OwnerFirstNameAr                     string     `json:"owner_first_name_ar"`
	OwnerLastNameEn                      string     `json:"owner_last_name_en"`
	OwnerLastNameAr                      string     `json:"owner_last_name_ar"`
	OwnerInitial                         string     `json:"owner_initial"`
	OwnerColor                           string     `json:"owner_color"`
	OwnerEmail                           string     `json:"owner_email"`
	OwnerProfileImagePath                string     `json:"owner_profile_image_path"`
	ControlImplementationStatusTitleEn   *string    `json:"control_implementation_status_title_en"`
	ControlImplementationStatusTitleAr   *string    `json:"control_implementation_status_title_ar"`
	ControlImplementationStatusColorCode *string    `json:"control_implementation_status_color_code"`
	ControlEfficiencyStatusTitleEn       *string    `json:"control_efficiency_status_title_en"`
	ControlEfficiencyStatusTitleAr       *string    `json:"control_efficiency_status_title_ar"`
	ControlEfficiencyStatusColorCode     *string    `json:"control_efficiency_status_color_code"`
	Comments                             *string    `json:"comments"`
	CreatedUserID                        uint       `json:"created_user_id"`
	CreatedUserFirstNameEn               string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr               string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn                string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr                string     `json:"created_user_last_name_ar"`
	CreatedUserInitial                   string     `json:"created_user_initial"`
	CreatedUserColor                     string     `json:"created_user_color"`
	CreatedUserEmail                     string     `json:"created_user_email"`
	CreatedUserProfileImage              string     `json:"created_user_profile_image_path"`
	CreatedAt                            string     `json:"created_at"`
	UpdatedUserID                        *uint      `json:"updated_user_id"`
	UpdatedUserFirstNameEn               *string    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr               *string    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn                *string    `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr                *string    `json:"updated_user_last_name_ar"`
	UpdatedUserInitial                   string     `json:"updated_user_initial"`
	UpdatedUserColor                     string     `json:"updated_user_color"`
	UpdatedUserEmail                     string     `json:"updated_user_email"`
	UpdatedUserProfileImagePath          string     `json:"updated_user_profile_image_path"`
	UpdatedAt                            *time.Time `json:"updated_at"`

	Documents []*AssetControlAnalysisDocumentDTO `json:"documents" gorm:"-"`
}

type AssetControlAnalysisDocumentDTO struct {
	ID           uint   `json:"id"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	Size         int    `json:"size"`
	Token        string `json:"token"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
}

type AssetControlAnalysisDetailDTO struct {
	ID                            uint                                             `json:"id"`
	ControlID                     uint                                             `json:"control_id"`
	ControlEfficiencyStatusID     uint                                             `json:"control_efficiency_status_id"`
	ControlImplementationStatusID uint                                             `json:"control_implementation_status_id"`
	Comments                      *string                                          `json:"comments"`
	Control                       *organizationDTOs.ControlDetailDTO               `json:"control"`
	ControlImplementationStatus   *organizationDTOs.ControlImplementationStatusDTO `json:"control_implementation_status"`
	ControlEfficiencyStatus       *organizationDTOs.ControlEfficiencyStatusDTO     `json:"control_efficiency_status"`
	CreatedAt                     time.Time                                        `json:"created_at"`
	UpdatedAt                     *time.Time                                       `json:"updated_at"`
	CreatedByUser                 organizationDTOs.UserMinimalDTO                  `json:"created_by"`
	UpdatedByUser                 *organizationDTOs.UserMinimalDTO                 `json:"updated_by"`
}

type AssetControlFormDTO struct {
	ControlIds []uint `json:"control_ids"`
}

type AssetControlUpdateFormDTO struct {
	ControlEfficiencyStatusID     uint                          `json:"control_efficiency_status_id"`
	ControlImplementationStatusID uint                          `json:"control_implementation_status_id"`
	Comments                      *string                       `json:"comments"`
	Documents                     []*commonDTOs.DocumentFormDTO `json:"documents"`
}

type AssetControlAnalysisListQueryParamsDTO struct {
	Sort    string `json:"sort_order"`
	SortBy  string `json:"sort_by"`
	Keyword string `json:"keyword"`
}
