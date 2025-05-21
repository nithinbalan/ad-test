package dtos

import (
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	tprmDTOs "jethur/grc/tprm/dtos"
	"time"
)

type StakeholderListDTO struct {
	ID                          int     `json:"id"`
	TitleEn                     string  `json:"title_en"`
	TitleAr                     string  `json:"title_ar"`
	TypeID                      *uint   `json:"type_id"`
	TypeTitleEn                 *string `json:"type_title_en"`
	TypeTitleAr                 *string `json:"type_title_ar"`
	StakeholderCategoryID       *uint   `json:"stakeholder_category_id"`
	StakeholderCategoryTitleEn  *string `json:"stakeholder_category_title_en"`
	StakeholderCategoryTitleAr  *string `json:"stakeholder_category_title_ar"`
	StakeholderCategoryType     *string `json:"stakeholder_category_type"`
	Role                        *string `json:"role"`
	InfluenceLevelID            *uint   `json:"influence_level_id"`
	InfluenceLevelTitleEn       *string `json:"influence_level_title_en"`
	InfluenceLevelTitleAr       *string `json:"influence_level_title_ar"`
	InfluenceLevelType          *string `json:"influence_level_type"`
	ImpactLevelID               *uint   `json:"impact_level_id"`
	ImpactLevelTitleEn          *string `json:"impact_level_title_en"`
	ImpactLevelTitleAr          *string `json:"impact_level_title_ar"`
	ImpactLevelType             *string `json:"impact_level_type"`
	ContactInformation          *string `json:"contact_information"`
	OrganizationUnitID          *uint   `json:"organization_unit_id"`
	OrganizationUnitTitleEn     *string `json:"organization_unit_title_en"`
	OrganizationUnitTitleAr     *string `json:"organization_unit_title_ar"`
	ThirdPartyRegisterID        *uint   `json:"third_party_register_id"`
	ThirdPartyRegisterTitle     *string `json:"third_party_register_title"`
	ClientID                    *uint   `json:"client_id"`
	ClientTitle                 *string `json:"client_title"`
	UserID                      *uint   `json:"user_id"`
	UserFirstNameEn             string  `json:"user_first_name_en"`
	UserFirstNameAr             string  `json:"user_first_name_ar"`
	UserLastNameEn              string  `json:"user_last_name_en"`
	UserLastNameAr              string  `json:"user_last_name_ar"`
	UserInitial                 string  `json:"user_initial"`
	UserColor                   string  `json:"user_color"`
	UserEmail                   string  `json:"user_email"`
	UserProfileImagePath        string  `json:"user_profile_image_path"`
	CreatedUserID               uint    `json:"created_user_id"`
	CreatedUserFirstNameEn      string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string  `json:"created_user_last_name_ar"`
	CreatedUserInitial          string  `json:"created_user_initial"`
	CreatedUserColor            string  `json:"created_user_color"`
	CreatedUserEmail            string  `json:"created_user_email"`
	CreatedUserProfileImagePath string  `json:"created_user_profile_image_path"`
	UpdatedUserID               *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn      *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          *string `json:"updated_user_initial"`
	UpdatedUserColor            *string `json:"updated_user_color"`
	UpdatedUserEmail            *string `json:"updated_user_email"`
	UpdatedUserProfileImagePath *string `json:"updated_user_profile_image_path"`
}

type StakeholderFormDTO struct {
	TitleEn               string  `json:"title_en" binding:"required"`
	TitleAr               string  `json:"title_ar" binding:"required"`
	TypeID                *uint   `json:"type_id" binding:"required"`
	StakeholderCategoryID *uint   `json:"stakeholder_category_id" binding:"required"`
	Role                  *string `json:"role"`
	InfluenceLevelID      *uint   `json:"influence_level_id"`
	ImpactLevelID         *uint   `json:"impact_level_id"`
	ContactInformation    *string `json:"contact_information"`
	OrganizationUnitID    *uint   `json:"organization_unit_id"`
	ThirdPartyRegisterID  *uint   `json:"third_party_register_id"`
	ClientID              *uint   `json:"client_id"`
	UserID                *uint   `json:"user_id"`
}

type StakeholderListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	Export      bool
}

type StakeholderDetailDTO struct {
	ID                  uint                                       `json:"id"`
	TitleEn             string                                     `json:"title_en"`
	TitleAr             string                                     `json:"title_ar"`
	Type                *organizationDTOs.ProcessDependencyTypeDTO `json:"type"`
	StakeholderCategory *StakeholderCategoryDTO                    `json:"stakeholder_category"`
	Role                *string                                    `json:"role"`
	InfluenceLevel      *masterDTOs.SeverityLevelDTO               `json:"influence_level"`
	ImpactLevel         *masterDTOs.SeverityLevelDTO               `json:"impact_level"`
	ContactInformation  *string                                    `json:"contact_information"`
	OrganizationUnit    *settingsDTOs.OrganizationUnitDTO          `json:"organization_unit"`
	ThirdPartyRegister  *tprmDTOs.ThirdPartyRegisterDetailDTO      `json:"third_party_register"`
	Client              *tprmDTOs.ThirdPartyRegisterDetailDTO      `json:"client"`
	User                *organizationDTOs.UserMinimalDTO           `json:"user"`
	CreatedAt           time.Time                                  `json:"created_at"`
	UpdatedAt           *time.Time                                 `json:"updated_at"`
	CreatedByUser       organizationDTOs.UserMinimalDTO            `json:"created_by"`
	UpdatedByUser       *organizationDTOs.UserMinimalDTO           `json:"updated_by"`
}
