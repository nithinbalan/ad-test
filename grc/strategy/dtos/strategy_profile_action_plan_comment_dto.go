package dtos

import (
	"time"

	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
)

type StrategyProfileActionPlanCommentListDTO struct {
	ID                          uint                                           `json:"id"`
	Comment                     string                                         `json:"comment" binding:"required"`
	StrategyProfileActionPlanID uint                                           `json:"test_and_exercise_action_plan_id"`
	CreatedAt                   string                                         `json:"created_at"`
	CreatedBYID                 uint                                           `json:"created_by_id"`
	CreatedUserFirstNameEn      string                                         `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string                                         `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string                                         `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string                                         `json:"created_user_last_name_ar"`
	CreatedUserInitial          string                                         `json:"created_user_initial"`
	CreatedUserColor            string                                         `json:"created_user_color"`
	CreatedUserEmail            string                                         `json:"created_user_email"`
	CreatedUserProfileImagePath string                                         `json:"created_user_profile_image_path"`
	Documents                   []*StrategyProfileActionPlanCommentDocumentDTO `json:"documents" gorm:"-"`
}

type StrategyProfileActionPlanCommentFormDTO struct {
	Comment   string                        `json:"comment" binding:"required,max=1000"`
	Documents []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
}

type StrategyProfileActionPlanCommentDetailDTO struct {
	ID                          uint                                           `json:"id"`
	Comment                     string                                         `json:"comment"`
	StrategyProfileActionPlanID uint                                           `json:"test_and_exercise_action_plan_id"`
	CreatedAt                   time.Time                                      `json:"created_at"`
	CreatedByUser               organizationDTOs.UserMinimalDTO                `json:"created_by"`
	Documents                   []*StrategyProfileActionPlanCommentDocumentDTO `json:"documents"`
}

type StrategyProfileActionPlanCommentDocumentDTO struct {
	ID           uint   `json:"id"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	Size         int    `json:"size"`
	Token        string `json:"token"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
}
