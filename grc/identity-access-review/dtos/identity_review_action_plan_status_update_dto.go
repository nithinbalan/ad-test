package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
)

type IdentityReviewActionPlanStatusUpdateListDTO struct {
	ID                         uint                                                   `json:"id"`
	IdentityReviewActionPlanID uint                                                   `json:"identity_review_action_plan_id"`
	ActionPlanTitle            string                                                 `json:"action_plan_title"`
	ActionPlanStatusID         uint                                                   `json:"action_plan_status_id"`
	ActionPlanStatusTitleEn    string                                                 `json:"action_plan_status_title_en"`
	ActionPlanStatusTitleAr    string                                                 `json:"action_plan_status_title_ar"`
	ActionPlanStatusColorCode  string                                                 `json:"action_plan_status_color_code"`
	ActionPlanStatusType       string                                                 `json:"action_plan_status_type"`
	Comment                    string                                                 `json:"comment"`
	Percentage                 *uint                                                  `json:"percentage"`
	CreatedAt                  string                                                 `json:"created_at"`
	UpdatedAt                  *string                                                `json:"updated_at"`
	CreatedUserID              uint                                                   `json:"created_user_id"`
	CreatedUserFirstNameEn     string                                                 `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr     string                                                 `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn      string                                                 `json:"created_user_last_name_en"`
	CreatedUserLastNameAr      string                                                 `json:"created_user_last_name_ar"`
	Documents                  []*models.IdentityReviewActionPlanStatusUpdateDocument `json:"documents" gorm:"-"`
}

type IdentityReviewActionPlanStatusUpdateFormDTO struct {
	ActionPlanStatusID uint                          `json:"action_plan_status_id" binding:"required"`
	Comment            string                        `json:"comment" binding:"required"`
	Percentage         *uint                         `json:"percentage" binding:"required,min=0,max=100"`
	Documents          []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
}

type IdentityReviewActionPlanReviewFormDTO struct {
	Comment string `json:"comment" binding:"required"`
}
