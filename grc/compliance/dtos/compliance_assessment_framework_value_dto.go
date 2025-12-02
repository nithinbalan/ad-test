package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type ComplianceAssessmentFrameworkValueListDTO struct {
	ID                                   uint    `json:"id"`
	ComplianceAssessmentFrameworkID      uint    `json:"compliance_assessment_framework_id"`
	ComplianceAssessmentFrameworkTitleEn string  `json:"compliance_assessment_framework_title_en"`
	ComplianceAssessmentFrameworkTitleAr string  `json:"compliance_assessment_framework_title_ar"`
	TitleEn                              *string `json:"title_en"`
	TitleAr                              *string `json:"title_ar"`
	Description                          *string `json:"description"`
	ComplianceStatusID                   uint    `json:"compliance_status_id"`
	ComplianceStatusTitleEn              string  `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr              string  `json:"compliance_status_title_ar"`
	ComplianceStatusType                 string  `json:"compliance_status_type"`
	ComplianceStatusColorCode            string  `json:"compliance_status_color_code"`
	Value                                *int    `json:"value"`
	ScoreMax                             *int    `json:"score_max"`
	ScoreMin                             *int    `json:"score_min"`
	IsApplicable                         *bool   `json:"is_applicable"`
	CreatedUserID                        uint    `json:"created_user_id"`
	CreatedUserFirstNameEn               string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr               string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn                string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr                string  `json:"created_user_last_name_ar"`
	CreatedUserInitial                   string  `json:"created_user_initial"`
	CreatedUserColor                     string  `json:"created_user_color"`
	CreatedUserEmail                     string  `json:"created_user_email"`
	CreatedUserProfileImagePath          string  `json:"created_user_profile_image_path"`
	UpdatedUserID                        *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn               *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr               *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn                *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr                *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial                   *string `json:"updated_user_initial"`
	UpdatedUserColor                     *string `json:"updated_user_color"`
	UpdatedUserEmail                     *string `json:"updated_user_email"`
	UpdatedUserProfileImagePath          *string `json:"updated_user_profile_image_path"`
}

type ComplianceAssessmentFrameworkValueListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
	Filters  map[string][]string
}

type ComplianceAssessmentFrameworkValueListResponseDTO struct {
	Items      []*ComplianceAssessmentFrameworkValueListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO                 `json:"pagination"`
}

type ComplianceAssessmentFrameworkValueDetailDTO struct {
	ID                              uint                                           `json:"id"`
	ComplianceAssessmentFrameworkID uint                                           `json:"compliance_assessment_framework_id"`
	ComplianceAssessmentFramework   ComplianceAssessmentFrameworkDetailResponseDTO `json:"compliance_assessment_framework"`
	TitleEn                         *string                                        `json:"title_en"`
	TitleAr                         *string                                        `json:"title_ar"`
	Description                     *string                                        `json:"description"`
	Value                           *int                                           `json:"value"`
	ScoreMax                        *int                                           `json:"score_max"`
	ScoreMin                        *int                                           `json:"score_min"`
	IsApplicable                    *bool                                          `json:"is_applicable"`
	ComplianceStatusID              uint                                           `json:"compliance_status_id"`
	ComplianceStatus                ComplianceStatusDTO                            `json:"compliance_status"`
	CreatedByUser                   organizationDTOs.UserMinimalDTO                `json:"created_by"`
	UpdatedByUser                   *organizationDTOs.UserMinimalDTO               `json:"updated_by"`
	CreatedAt                       time.Time                                      `json:"created_at"`
	UpdatedAt                       *time.Time                                     `json:"updated_at"`
}

type ComplianceAssessmentFrameworkValueFormDTO struct {
	TitleEn            *string `json:"title_en" binding:"required"`
	TitleAr            *string `json:"title_ar" binding:"required"`
	Description        *string `json:"description"`
	Value              *int    `json:"value"`
	ScoreMax           *int    `json:"score_max"`
	ScoreMin           *int    `json:"score_min"`
	IsApplicable       bool    `json:"is_applicable"`
	ComplianceStatusID uint    `json:"compliance_status_id" binding:"required"`
}
