package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
)

type PenetrationTestFindingStatusUpdateFormDTO struct {
	ActionPlanStatusID uint                          `json:"action_plan_status_id" binding:"required"`
	Comment            string                        `json:"comment" binding:"required"`
	Percentage         *uint                         `json:"percentage" binding:"required,min=0,max=100"`
	Documents          []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
}

type PenetrationTestFindingStatusUpdateListDTO struct {
	ID                          uint                                                 `json:"id"`
	PenetrationTestFindingID    uint                                                 `json:"penetration_test_finding_id"`
	ActionPlanStatusID          uint                                                 `json:"action_plan_status_id"`
	ActionPlanStatusTitleEn     string                                               `json:"action_plan_status_title_en"`
	ActionPlanStatusTitleAr     string                                               `json:"action_plan_status_title_ar"`
	ActionPlanStatusColorCode   string                                               `json:"action_plan_status_color_code"`
	ActionPlanStatusType        string                                               `json:"action_plan_status_type"`
	Percentage                  *uint                                                `json:"percentage"`
	Comment                     string                                               `json:"comment"`
	CreatedAt                   string                                               `json:"created_at"`
	UpdatedAt                   *string                                              `json:"updated_at"`
	CreatedUserID               uint                                                 `json:"created_user_id"`
	CreatedUserFirstNameEn      string                                               `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string                                               `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string                                               `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string                                               `json:"created_user_last_name_ar"`
	CreatedUserInitial          string                                               `json:"created_user_initial"`
	CreatedUserColor            string                                               `json:"created_user_color"`
	CreatedUserEmail            string                                               `json:"created_user_email"`
	CreatedUserProfileImagePath string                                               `json:"created_user_profile_image_path"`
	UpdatedUserID               *uint                                                `json:"updated_user_id"`
	UpdatedUserFirstNameEn      *string                                              `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      *string                                              `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       *string                                              `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       *string                                              `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          *string                                              `json:"updated_user_initial"`
	UpdatedUserColor            *string                                              `json:"updated_user_color"`
	UpdatedUserEmail            *string                                              `json:"updated_user_email"`
	UpdatedUserProfileImagePath *string                                              `json:"updated_user_profile_image_path"`
	Documents                   []*models.PenetrationTestFindingStatusUpdateDocument `json:"documents" gorm:"-"`
}
