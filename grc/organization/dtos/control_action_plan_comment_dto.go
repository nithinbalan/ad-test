package dtos

import (
	"time"

	commonDTOs "jethur/grc/common/dtos"
)

type ControlActionPlanCommentListDTO struct {
	ID                          uint                                   `json:"id"`
	Comment                     string                                 `json:"comment" binding:"required"`
	ControlActionPlanID         uint                                   `json:"control_action_plan_id"`
	CreatedAt                   string                                 `json:"created_at"`
	CreatedBYID                 uint                                   `json:"created_by_id"`
	CreatedUserFirstNameEn      string                                 `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string                                 `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string                                 `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string                                 `json:"created_user_last_name_ar"`
	CreatedUserInitial          string                                 `json:"created_user_initial"`
	CreatedUserColor            string                                 `json:"created_user_color"`
	CreatedUserEmail            string                                 `json:"created_user_email"`
	CreatedUserProfileImagePath string                                 `json:"created_user_profile_image_path"`
	Documents                   []*ControlActionPlanCommentDocumentDTO `json:"documents" gorm:"-"`
}

type ControlActionPlanCommentFormDTO struct {
	Comment   string                        `json:"comment" binding:"required,max=1000"`
	Documents []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
}

type ControlActionPlanCommentDetailDTO struct {
	ID                  uint                                   `json:"id"`
	Comment             string                                 `json:"comment"`
	ControlActionPlanID uint                                   `json:"control_action_plan_id"`
	CreatedAt           time.Time                              `json:"created_at"`
	CreatedByUser       UserMinimalDTO                         `json:"created_by"`
	Documents           []*ControlActionPlanCommentDocumentDTO `json:"documents"`
}

type ControlActionPlanCommentDocumentDTO struct {
	ID           uint   `json:"id"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	Size         int    `json:"size"`
	Token        string `json:"token"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
}
