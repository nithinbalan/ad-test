package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type ComplianceAssessmentFrameworkListDTO struct {
	ID                                             uint       `json:"id"`
	TitleEn                                        string     `json:"title_en"`
	TitleAr                                        string     `json:"title_ar"`
	Description                                    *string    `json:"description"`
	ComplianceAssessmentFrameworkTypeID            uint       `json:"compliance_assessment_framework_type_id"`
	ComplianceAssessmentFrameworkTypeTitleEn       string     `json:"compliance_assessment_framework_type_title_en"`
	ComplianceAssessmentFrameworkTypeTitleAr       string     `json:"compliance_assessment_framework_type_title_ar"`
	ComplianceAssessmentFrameworkTypeDescriptionEn string     `json:"compliance_assessment_framework_type_description_en"`
	ComplianceAssessmentFrameworkTypeDescriptionAr string     `json:"compliance_assessment_framework_type_description_ar"`
	Status                                         string     `json:"status"`
	CreatedUserID                                  uint       `json:"created_user_id"`
	CreatedUserFirstNameEn                         string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr                         string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn                          string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr                          string     `json:"created_user_last_name_ar"`
	CreatedUserInitial                             string     `json:"created_user_initial"`
	CreatedUserColor                               string     `json:"created_user_color"`
	CreatedUserEmail                               string     `json:"created_user_email"`
	CreatedUserProfileImagePath                    string     `json:"created_user_profile_image_path"`
	UpdatedUserID                                  *uint      `json:"updated_user_id"`
	UpdatedUserFirstNameEn                         *string    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr                         *string    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn                          *string    `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr                          *string    `json:"updated_user_last_name_ar"`
	UpdatedUserInitial                             *string    `json:"updated_user_initial"`
	UpdatedUserColor                               *string    `json:"updated_user_color"`
	UpdatedUserEmail                               *string    `json:"updated_user_email"`
	UpdatedUserProfileImagePath                    *string    `json:"updated_user_profile_image_path"`
	CreatedAt                                      time.Time  `json:"created_at"`
	UpdatedAt                                      *time.Time `json:"updated_at"`
}

type ComplianceAssessmentFrameworkListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
	Filters map[string][]string
}

type ComplianceAssessmentFrameworkListResponseDTO struct {
	Items      []*ComplianceAssessmentFrameworkListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO            `json:"pagination"`
}

type ComplianceAssessmentFrameworkFormDTO struct {
	TitleEn                   string  `json:"title_en" binding:"required"`
	TitleAr                   string  `json:"title_ar" binding:"required"`
	Description               *string `json:"description"`
	AssessmentFrameworkTypeID uint    `json:"compliance_assessment_framework_type_id" binding:"required"`
}

type ComplianceAssessmentFrameworkDetailResponseDTO struct {
	ID                      uint                                           `json:"id"`
	TitleEn                 string                                         `json:"title_en"`
	TitleAr                 string                                         `json:"title_ar"`
	Description             *string                                        `json:"description"`
	Status                  string                                         `json:"status"`
	AssessmentFrameworkType *ComplianceAssessmentFrameworkTypeDTO          `json:"compliance_assessment_framework_type"`
	Values                  []*ComplianceAssessmentFrameworkValueDetailDTO `json:"values"`
	CreatedAt               time.Time                                      `json:"created_at"`
	UpdatedAt               *time.Time                                     `json:"updated_at"`
}
