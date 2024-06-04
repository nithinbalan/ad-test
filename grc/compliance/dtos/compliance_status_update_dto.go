package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
)

type ComplianceStatusUpdateListDTO struct {
	ID                        uint    `json:"id"`
	ComplianceTitleEn         string  `json:"compliance_title_en"`
	ComplianceTitleAr         string  `json:"compliance_title_ar"`
	ComplianceDescriptionEn   string  `json:"compliance_description_en"`
	ComplianceDescriptionAr   string  `json:"compliance_description_ar"`
	ComplianceStatusTitleEn   string  `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr   string  `json:"compliance_status_title_ar"`
	ComplianceStatusColorCode string  `json:"compliance_status_color_code"`
	Comment                   string  `json:"comment"`
	CreatedAt                 string  `json:"created_at"`
	UpdatedAt                 string  `json:"updated_at"`
	CreatedUserFirstNameEn    string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr    string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn     string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr     *string `json:"created_user_last_name_ar"`
	UpdatedUserFirstNameEn    string  `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr    string  `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn     string  `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr     string  `json:"updated_user_last_name_ar"`
}

type ComplianceStatusUpdateFormDTO struct {
	ComplianceStatusID uint                          `json:"compliance_status_id" binding:"required"`
	Comment            string                        `json:"comment" binding:"required"`
	Documents          []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
}
