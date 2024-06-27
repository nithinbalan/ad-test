package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type ComplianceStatusUpdateListDTO struct {
	ID                          uint                                     `json:"id"`
	ComplianceTitleEn           string                                   `json:"compliance_title_en"`
	ComplianceTitleAr           string                                   `json:"compliance_title_ar"`
	ComplianceDescriptionEn     string                                   `json:"compliance_description_en"`
	ComplianceDescriptionAr     string                                   `json:"compliance_description_ar"`
	ComplianceStatusID          string                                   `json:"compliance_status_id"`
	ComplianceStatusTitleEn     string                                   `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr     string                                   `json:"compliance_status_title_ar"`
	ComplianceStatusColorCode   string                                   `json:"compliance_status_color_code"`
	Comment                     string                                   `json:"comment"`
	CreatedAt                   string                                   `json:"created_at"`
	UpdatedAt                   *string                                  `json:"updated_at"`
	CreatedUserFirstNameEn      string                                   `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string                                   `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string                                   `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       *string                                  `json:"created_user_last_name_ar"`
	CreatedUserInitial          string                                   `json:"created_user_initial"`
	CreatedUserColor            string                                   `json:"created_user_color"`
	CreatedUserEmail            string                                   `json:"created_user_email"`
	CreatedUserProfileImagePath string                                   `json:"created_user_profile_image_path"`
	UpdatedUserFirstNameEn      *string                                  `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      *string                                  `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       *string                                  `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       *string                                  `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          *string                                  `json:"updated_user_initial"`
	UpdatedUserColor            *string                                  `json:"updated_user_color"`
	UpdatedUserEmail            *string                                  `json:"updated_user_email"`
	UpdatedUserProfileImagePath *string                                  `json:"updated_user_profile_image_path"`
	Documents                   []*models.ComplianceStatusUpdateDocument `json:"documents" gorm:"-"`
}

type ComplianceStatusUpdateFormDTO struct {
	ComplianceStatusID uint                          `json:"compliance_status_id" binding:"required"`
	Comment            string                        `json:"comment" binding:"required"`
	Documents          []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
}

type ComplianceStatusUpdateDetailDTO struct {
	ID uint `json:"id"`
	// ComplianceRegister *ComplianceRegisterDetailDTO             `json:"compliance_register"`
	ComplianceStatus *ComplianceStatusDTO                     `json:"compliance_status"`
	Documents        []*models.ComplianceStatusUpdateDocument `json:"documents"`
	Comment          string                                   `json:"comment"`
	CreatedAt        time.Time                                `json:"created_at"`
	UpdatedAt        *time.Time                               `json:"updated_at"`
	CreatedByUser    organizationDTOs.UserMinimalDTO          `json:"created_by"`
	UpdatedByUser    *organizationDTOs.UserMinimalDTO         `json:"updated_by"`
}
