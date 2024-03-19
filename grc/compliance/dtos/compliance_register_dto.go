package dtos

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
)

type ComplianceRegisterListDTO struct {
	ID                               uint   `json:"id"`
	ComplianceCategoryTitleEn        string `json:"compliance_category_title_en"`
	ComplianceCategoryTitleAr        string `json:"compliance_category_title_ar"`
	ComplianceStatusTitleEn          string `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr          string `json:"compliance_status_title_ar"`
	ComplianceStatusColorCode        string `json:"compliance_status_color_code"`
	ResponsibleUnitTitleEn           string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr           string `json:"responsible_unit_title_ar"`
	ResponsibleUserFirstNameEn       string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr       string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn        string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr        string `json:"responsible_user_last_name_ar"`
	ReviewFrequencyTitleEn           string `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr           string `json:"review_frequency_title_ar"`
	ReviewResponsibleUserFirstNameEn string `json:"review_responsible_user_first_name_en"`
	ReviewResponsibleUserFirstNameAr string `json:"review_responsible_user_first_name_ar"`
	ReviewResponsibleUserLastNameEn  string `json:"review_responsible_user_last_name_en"`
	ReviewResponsibleUserLastNameAr  string `json:"review_responsible_user_last_name_ar"`
	SeverityLevelTitleEn             string `json:"severity_level_title_en"`
	SeverityLevelTitleAr             string `json:"severity_level_title_ar"`
	SeverityLevelColorCode           string `json:"severity_level_color_code"`
	ReferenceNumber                  string `json:"reference_number"`
	StandardCode                     string `json:"standard_code"`
	TitleEn                          string `json:"title_en"`
	TitleAr                          string `json:"title_ar"`
	DescriptionEn                    string `json:"description_en"`
	DescriptionAr                    string `json:"description_ar"`
	Version                          string `json:"version"`
	Deadline                         string `json:"deadline"`
	IssueDate                        string `json:"issue_date"`
	ExpiryDate                       string `json:"expiry_date"`
	NextReviewDate                   string `json:"next_review_date"`
	CreatedAt                        string `json:"created_at"`
	UpdatedAt                        string `json:"updated_at"`
	CreatedUserFirstNameEn           string `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr           string `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn            string `json:"created_user_last_name_en"`
	CreatedUserLastNameAr            string `json:"created_user_last_name_ar"`
	UpdatedUserFirstNameEn           string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr           string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn            string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr            string `json:"updated_user_last_name_ar"`
	LifecycleStatusTitleEn           string `json:"lifecycle_status_title_en"`
	LifecycleStatusTitleAr           string `json:"lifecycle_status_title_ar"`
	LifecycleStatusColorCode         string `json:"lifecycle_status_color_code"`
}

type ComplianceRegisterListResponseDTO struct {
	Items      []*ComplianceRegisterListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type ComplianceRegisterFormDTO struct {
	ComplianceCategoryID    uint   `json:"compliance_category_id" binding:"required"`
	ComplianceStatusID      uint   `json:"compliance_status_id" binding:"required"`
	ResponsibleUnitID       uint   `json:"responsible_unit_id" binding:"required"`
	ResponsibleUserID       uint   `json:"responsible_user_id" binding:"required"`
	ReviewFrequencyID       uint   `json:"review_frequency_id" binding:"required"`
	ReviewResponsibleUserID uint   `json:"review_responsible_user_id"`
	SeverityLevelID         uint   `json:"severity_level_id"`
	TitleEn                 string `json:"title_en" binding:"required,max=255"`
	TitleAr                 string `json:"title_ar" binding:"required,max=255"`
	DescriptionEn           string `json:"description_en"`
	DescriptionAr           string `json:"description_ar"`
	Version                 string `json:"version" binding:"omitempty,max=45"`
	Deadline                string `json:"deadline" binding:"omitempty,datetime=2006-01-02"`
	IssueDate               string `json:"issue_date" binding:"omitempty,datetime=2006-01-02"`
	ExpiryDate              string `json:"expiry_date" binding:"omitempty,datetime=2006-01-02"`
	NextReviewDate          string `json:"next_review_date" binding:"omitempty,datetime=2006-01-02"`
	StandardCode            string `json:"standard_code" binding:"omitempty,max=45"`
}

func (dto *ComplianceRegisterFormDTO) Validate() error {
	if !utility.IsValidDate(dto.Deadline) {
		return utility.NewValidationError(errors.New("invalid deadline"), "invalid_date", "deadline")
	}
	if !utility.IsValidDate(dto.IssueDate) {
		return utility.NewValidationError(errors.New("invalid issue date"), "invalid_date", "issue date")
	}
	if !utility.IsValidDate(dto.ExpiryDate) {
		return utility.NewValidationError(errors.New("invalid expiry date"), "invalid_date", "expiry date")
	}
	if !utility.IsValidDate(dto.NextReviewDate) {
		return utility.NewValidationError(errors.New("invalid next review date"), "invalid_date", "next review date")
	}
	return nil
}
