package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
)

type StrategyProfileActionPlanStatusUpdateListDTO struct {
	ID                             uint                                                `json:"id"`
	StrategyProfileActionPlanTitle string                                              `json:"test_and_exercise_recovery_action_plan_title"`
	ActionPlanStatusTitleEn        string                                              `json:"action_plan_status_title_en"`
	ActionPlanStatusTitleAr        string                                              `json:"action_plan_status_title_ar"`
	ActionPlanStatusColorCode      string                                              `json:"action_plan_status_color_code"`
	ActionPlanStatusType           string                                              `json:"action_plan_status_type"`
	Comment                        string                                              `json:"comment"`
	Percentage                     string                                              `json:"percentage"`
	CreatedAt                      string                                              `json:"created_at"`
	UpdatedAt                      string                                              `json:"updated_at"`
	CreatedUserID                  uint                                                `json:"created_user_id"`
	CreatedUserFirstNameEn         string                                              `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr         string                                              `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn          string                                              `json:"created_user_last_name_en"`
	CreatedUserLastNameAr          *string                                             `json:"created_user_last_name_ar"`
	CreatedUserInitial             string                                              `json:"created_user_initial"`
	CreatedUserColor               string                                              `json:"created_user_color"`
	CreatedUserEmail               string                                              `json:"created_user_email"`
	CreatedUserProfileImagePath    string                                              `json:"created_user_profile_image_path"`
	UpdatedUserID                  *uint                                               `json:"updated_user_id"`
	UpdatedUserFirstNameEn         string                                              `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr         string                                              `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn          string                                              `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr          string                                              `json:"updated_user_last_name_ar"`
	UpdatedUserInitial             string                                              `json:"updated_user_initial"`
	UpdatedUserColor               string                                              `json:"updated_user_color"`
	UpdatedUserEmail               string                                              `json:"updated_user_email"`
	UpdatedUserProfileImagePath    string                                              `json:"updated_user_profile_image_path"`
	Documents                      []*StrategyProfileActionPlanStatusUpdateDocumentDTO `json:"documents" gorm:"-"`
}

type StrategyProfileActionPlanStatusUpdateFormDTO struct {
	ActionPlanStatusID uint                          `json:"action_plan_status_id" binding:"required"`
	Comment            string                        `json:"comment" binding:"required"`
	Percentage         *int                          `json:"percentage" binding:"required"`
	Documents          []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
}

type StrategyProfileActionPlanStatusUpdateDocumentDTO struct {
	ID           uint   `json:"id"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	Size         int    `json:"size"`
	Token        string `json:"token"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
}

type StrategyProfileActionPlanReviewUserStatusUpdateFormDTO struct {
	Comment string `json:"comment"`
}
