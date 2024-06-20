package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
)

type ControlReviewListDTO struct {
	ID                           uint                            `json:"id"`
	ControlTitleEn               string                          `json:"control_title_en"`
	ControlTitleAr               string                          `json:"control_title_ar"`
	ControlDescription           string                          `json:"control_description"`
	ControlCode                  string                          `json:"control_code"`
	Comment                      string                          `json:"comment"`
	ReviewedAt                   string                          `json:"reviewed_at"`
	ReviewedUserFirstNameEn      string                          `json:"reviewed_user_first_name_en"`
	ReviewedUserFirstNameAr      string                          `json:"reviewed_user_first_name_ar"`
	ReviewedUserLastNameEn       string                          `json:"reviewed_user_last_name_en"`
	ReviewedUserLastNameAr       string                          `json:"reviewed_user_last_name_ar"`
	ReviewedUserInitial          string                          `json:"reviewed_user_initial"`
	ReviewedUserColor            string                          `json:"reviewed_user_color"`
	ReviewedUserEmail            string                          `json:"reviewed_user_email"`
	ReviewedUserProfileImagePath string                          `json:"reviewed_user_profile_image_path"`
	Documents                    []*models.ControlReviewDocument `json:"documents" gorm:"-"`
}

type ControlReviewFormDTO struct {
	Comment   string                        `json:"comment" binding:"required"`
	Documents []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
}
