package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"

	"time"
)

type PotentialRiskListDTO struct {
	ID                           uint    `json:"id"`
	RiskCategoryTitleEn          string  `json:"risk_category_title_en"`
	RiskCategoryTitleAr          string  `json:"risk_category_title_ar"`
	SeverityLevelTitleEn         string  `json:"severity_level_title_en"`
	SeverityLevelTitleAr         string  `json:"severity_level_title_ar"`
	SeverityLevelColorCode       string  `json:"severity_level_color_code"`
	PotentialRiskStatusTitleEn   string  `json:"potential_risk_status_title_en"`
	PotentialRiskStatusTitleAr   string  `json:"potential_risk_status_title_ar"`
	PotentialRiskStatusColorCode string  `json:"potential_risk_status_color_code"`
	Title                        string  `json:"title"`
	Description                  *string `json:"description"`
	RiskSource                   string  `json:"risk_source"`
	CreatedAt                    string  `json:"created_at"`
	UpdatedAt                    string  `json:"updated_at"`
	CreatedUserID                uint    `json:"created_user_id"`
	CreatedUserFirstNameEn       string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr       string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn        string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr        string  `json:"created_user_last_name_ar"`
	CreatedUserInitial           string  `json:"created_user_initial"`
	CreatedUserColor             string  `json:"created_user_color"`
	CreatedUserEmail             string  `json:"created_user_email"`
	CreatedUserProfileImagePath  string  `json:"created_user_profile_image_path"`
	UpdatedUserID                *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn       string  `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr       string  `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn        string  `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr        string  `json:"updated_user_last_name_ar"`
	UpdatedUserInitial           string  `json:"updated_user_initial"`
	UpdatedUserColor             string  `json:"updated_user_color"`
	UpdatedUserEmail             string  `json:"updated_user_email"`
	UpdatedUserProfileImagePath  string  `json:"updated_user_profile_image_path"`
}

type PotentialRiskListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	Export      bool
}

type PotentialRiskListResponseDTO struct {
	Items      []*PotentialRiskListDTO      `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type PotentialRiskFormDTO struct {
	RiskCategoryID  uint    `json:"risk_category_id" binding:"required"`
	SeverityLevelID uint    `json:"severity_level_id" binding:"required"`
	Title           string  `json:"title" binding:"required,max=100"`
	Description     *string `json:"description"`
	RiskSource      *string `json:"risk_source"`
}

type PotentialRiskDetailsDTO struct {
	ID                  uint                             `json:"id"`
	RiskCategory        RiskCategoryDTO                  `json:"risk_category"`
	SeverityLevel       masterDTOs.SeverityLevelDTO      `json:"severity_level"`
	PotentialRiskStatus PotentialRiskStatusDTO           `json:"potential_risk_status"`
	Title               string                           `json:"title"`
	Description         *string                          `json:"description"`
	RiskSource          string                           `json:"risk_source"`
	CreatedAt           time.Time                        `json:"created_at"`
	UpdatedAt           *time.Time                       `json:"updated_at"`
	CreatedByUser       organizationDTOs.UserMinimalDTO  `json:"created_by_user"`
	UpdatedByUser       *organizationDTOs.UserMinimalDTO `json:"updated_by_user"`
}
