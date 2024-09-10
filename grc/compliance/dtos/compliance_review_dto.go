package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
)

type ComplianceReviewListDTO struct {
	ID                           uint                               `json:"id"`
	ComplianceTitleEn            string                             `json:"compliance_title_en"`
	ComplianceTitleAr            string                             `json:"compliance_title_ar"`
	ComplianceDescriptionEn      string                             `json:"compliance_description_en"`
	ComplianceDescriptionAr      string                             `json:"compliance_description_ar"`
	Comment                      string                             `json:"comment"`
	ReviewedAt                   string                             `json:"reviewed_at"`
	ReviewedUserID               uint                               `json:"reviewed_user_id"`
	ReviewedUserFirstNameEn      string                             `json:"reviewed_user_first_name_en"`
	ReviewedUserFirstNameAr      string                             `json:"reviewed_user_first_name_ar"`
	ReviewedUserLastNameEn       string                             `json:"reviewed_user_last_name_en"`
	ReviewedUserLastNameAr       string                             `json:"reviewed_user_last_name_ar"`
	ReviewedUserInitial          string                             `json:"reviewed_user_initial"`
	ReviewedUserColor            string                             `json:"reviewed_user_color"`
	ReviewedUserEmail            string                             `json:"reviewed_user_email"`
	ReviewedUserProfileImagePath string                             `json:"reviewed_user_profile_image_path"`
	Documents                    []*models.ComplianceReviewDocument `json:"documents" gorm:"-"`
}

type ComplianceReviewFormDTO struct {
	Comment   string                        `json:"comment" binding:"required"`
	Documents []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
}
