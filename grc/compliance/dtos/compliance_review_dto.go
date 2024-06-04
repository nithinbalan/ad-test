package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
)

type ComplianceReviewListDTO struct {
	ID                      uint           `json:"id"`
	ComplianceTitleEn       string         `json:"compliance_title_en"`
	ComplianceTitleAr       string         `json:"compliance_title_ar"`
	ComplianceDescriptionEn string         `json:"compliance_description_en"`
	ComplianceDescriptionAr string         `json:"compliance_description_ar"`
	Comment                 string         `json:"comment"`
	ReviewedAt              string         `json:"reviewed_at"`
	ReviewedUserFirstNameEn string         `json:"reviewed_user_first_name_en"`
	ReviewedUserFirstNameAr string         `json:"reviewed_user_first_name_ar"`
	ReviewedUserLastNameEn  string         `json:"reviewed_user_last_name_en"`
	ReviewedUserLastNameAr  string         `json:"reviewed_user_last_name_ar"`
	// Documents               []*DocumentDTO `json:"documents" binding:"omitempty"`
}

type ComplianceReviewFormDTO struct {
	Comment   string                        `json:"comment" binding:"required"`
	Documents []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
}

type DocumentDTO struct {
	URL          *string `json:"url"`
	Name         string  `json:"name"`
	OriginalName string  `json:"original_name"`
	Type         string  `json:"type"`
	Token        string  `json:"token"`
	Size         int     `json:"size"`
}
