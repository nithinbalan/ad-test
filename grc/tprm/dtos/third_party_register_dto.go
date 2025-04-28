package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"time"
)

type ThirdPartyRegisterListDTO struct {
	ID                              uint    `json:"id"`
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
	CreatedAt                       string  `json:"created_at"`
	UpdatedAt                       string  `json:"updated_at"`
	CreatedByID                     uint    `json:"created_by_id"`
	CreatedUserFirstNameEn          string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr          string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn           string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr           string  `json:"created_user_last_name_ar"`
	CreatedUserInitial              string  `json:"created_user_initial"`
	CreatedUserColor                string  `json:"created_user_color"`
	CreatedUserEmail                string  `json:"created_user_email"`
	CreatedUserProfileImagePath     string  `json:"created_user_profile_image_path"`
	UpdatedByID                     *uint   `json:"updated_by_id"`
	UpdatedUserFirstNameEn          string  `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr          string  `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn           string  `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr           string  `json:"updated_user_last_name_ar"`
	UpdatedUserInitial              string  `json:"updated_user_initial"`
	UpdatedUserColor                string  `json:"updated_user_color"`
	UpdatedUserEmail                string  `json:"updated_user_email"`
	UpdatedUserProfileImagePath     string  `json:"updated_user_profile_image_path"`
}

type ThirdPartyRegisterFormDTO struct {
	Title              string  `json:"title" binding:"required"`
	ThirdPartyTypeID   uint    `json:"third_party_type_id" binding:"required"`
	CountryID          uint    `json:"country_id" binding:"required"`
	Address            *string `json:"address"`
	Website            *string `json:"website"`
	OrganizationUnitID *uint   `json:"organization_unit_id"`
	RepresentativeName *string `json:"representative_name"`
	RepresentativeRole *string `json:"representative_role"`
	Email              *string `json:"email"`
	PhoneNumber        *string `json:"phone_number"`
	ResponsibleUserID  *uint   `json:"responsible_user_id"`
	OnboardingDate     *string `json:"onboarding_date"`
	ServiceStartDate   *string `json:"service_start_date"`
	ServiceEndDate     *string `json:"service_end_date"`
	SeverityLevelID    *uint   `json:"severity_level_id"`
}

type ThirdPartyRegisterDetailDTO struct {
	ID                     uint                             `json:"id"`
	ReferenceNumber        string                           `json:"reference_number"`
	Title                  string                           `json:"title"`
	ThirdPartyType         ThirdPartyTypeDTO                `json:"third_party_type"`
	Country                masterDTOs.CountryDTO            `json:"country"`
	Address                *string                          `json:"address"`
	Website                *string                          `json:"website"`
	OrganizationUnit       settingsDTOs.OrganizationUnitDTO `json:"organization_unit"`
	RepresentativeName     *string                          `json:"representative_name"`
	RepresentativeRole     *string                          `json:"representative_role"`
	Email                  *string                          `json:"email"`
	PhoneNumber            *string                          `json:"phone_number"`
	ResponsibleUser        organizationDTOs.UserMinimalDTO  `json:"responsible_user"`
	OnboardingDateString   *string                          `json:"onboarding_date"`
	ServiceStartDateString *string                          `json:"service_start_date"`
	ServiceEndDateString   *string                          `json:"service_end_date"`
	SeverityLevel          *masterDTOs.SeverityLevelDTO     `json:"severity_level"`
	CreatedAt              time.Time                        `json:"created_at"`
	UpdatedAt              *time.Time                       `json:"updated_at"`
	CreatedBy              organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedBy              organizationDTOs.UserMinimalDTO  `json:"updated_by"`
}
type ThirdPartyRegisterListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	DateFrom    string
	DateTo      string
}
type ThirdPartyRegisterListResponseDTO struct {
	Items      []*ThirdPartyRegisterListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type ThirdPartyRegisterParsedDates struct {
	OnboardingDate   *time.Time `json:"onboarding_date"`
	ServiceStartDate *time.Time `json:"service_start_date"`
	ServiceEndDate   *time.Time `json:"service_end_date"`
}
