package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type RiskClosureListDTO struct {
	ID                          uint    `json:"id"`
	RiskRegisterID              uint    `json:"risk_register_id"`
	Justification               string  `json:"justification"`
	LessonLearned               *string `json:"lesson_learned"`
	WhatWorkedWell              *string `json:"what_worked_well"`
	Improvements                *string `json:"improvements"`
	Takeaways                   *string `json:"takeaways"`
	CreatedAt                   string  `json:"created_at"`
	UpdatedAt                   *string `json:"updated_at"`
	CreatedUserID               uint    `json:"created_user_id"`
	CreatedUserFirstNameEn      string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string  `json:"created_user_last_name_ar"`
	CreatedUserInitial          string  `json:"created_user_initial"`
	CreatedUserColor            string  `json:"created_user_color"`
	CreatedUserEmail            string  `json:"created_user_email"`
	CreatedUserProfileImagePath string  `json:"created_user_profile_image_path"`
	UpdatedUserID               *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn      *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          *string `json:"updated_user_initial"`
	UpdatedUserColor            *string `json:"updated_user_color"`
	UpdatedUserEmail            *string `json:"updated_user_email"`
	UpdatedUserProfileImagePath *string `json:"updated_user_profile_image_path"`
	CurrentRiskRatingID         *uint   `json:"current_risk_rating_id"`
	CurrentRiskRatingTitleEn    *string `json:"current_risk_rating_title_en"`
	CurrentRiskRatingTitleAr    *string `json:"current_risk_rating_title_ar"`
	CurrentRiskRatingColorCode  *string `json:"current_risk_rating_color_code"`
}

type RiskClosureDetailsDTO struct {
	ID             uint                             `json:"id"`
	RiskRegisterID uint                             `json:"risk_register_id"`
	Justification  string                           `json:"justification"`
	LessonLearned  *string                          `json:"lesson_learned"`
	WhatWorkedWell *string                          `json:"what_worked_well"`
	Improvements   *string                          `json:"improvements"`
	Takeaways      *string                          `json:"takeaways"`
	CreatedAt      time.Time                        `json:"created_at"`
	CreatedByUser  organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedAt      *time.Time                       `json:"updated_at,omitempty"`
	UpdatedByUser  *organizationDTOs.UserMinimalDTO `json:"updated_by"`
	RiskRegister   RiskRegisterDetailDTO            `json:"risk_register"`
	Documents      []*RiskClosureDocumentDTO        `json:"documents"`
}

type RiskClosureDocumentDTO struct {
	ID            uint   `json:"id"`
	RiskClosureID uint   `json:"risk_closure_id"`
	URL           string `json:"url"`
	Name          string `json:"name"`
	Size          int    `json:"size"`
	Token         string `json:"token"`
	OriginalName  string `json:"original_name"`
	Type          string `json:"type"`
}

type RiskClosureFormDTO struct {
	Justification  string                        `json:"justification" binding:"required"`
	LessonLearned  *string                       `json:"lesson_learned"`
	WhatWorkedWell *string                       `json:"what_worked_well"`
	Improvements   *string                       `json:"improvements"`
	Takeaways      *string                       `json:"takeaways"`
	Documents      []*commonDTOs.DocumentFormDTO `json:"documents"`
}
