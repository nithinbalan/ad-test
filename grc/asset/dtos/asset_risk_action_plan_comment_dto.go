package dtos

import (
	"time"

	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
)

type AssetRiskActionPlanCommentListDTO struct {
	ID                          uint                                     `json:"id"`
	Comment                     string                                   `json:"comment" binding:"required"`
	RiskActionPlanID            uint                                     `json:"risk_action_plan_id"`
	CreatedAt                   string                                   `json:"created_at"`
	CreatedByID                 uint                                     `json:"created_by_id"`
	CreatedUserFirstNameEn      string                                   `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string                                   `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string                                   `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string                                   `json:"created_user_last_name_ar"`
	CreatedUserInitial          string                                   `json:"created_user_initial"`
	CreatedUserColor            string                                   `json:"created_user_color"`
	CreatedUserEmail            string                                   `json:"created_user_email"`
	CreatedUserProfileImagePath string                                   `json:"created_user_profile_image_path"`
	Documents                   []*AssetRiskActionPlanCommentDocumentDTO `json:"documents" gorm:"-"`
}

type AssetRiskActionPlanCommentFormDTO struct {
	Comment   string                        `json:"comment" binding:"required,max=1000"`
	Documents []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
}

type AssetRiskActionPlanCommentDetailDTO struct {
	ID               uint                                     `json:"id"`
	Comment          string                                   `json:"comment"`
	RiskActionPlanID uint                                     `json:"risk_action_plan_id"`
	CreatedAt        time.Time                                `json:"created_at"`
	CreatedByUser    organizationDTOs.UserMinimalDTO          `json:"created_by"`
	Documents        []*AssetRiskActionPlanCommentDocumentDTO `json:"documents"`
}

type AssetRiskActionPlanCommentDocumentDTO struct {
	ID           uint   `json:"id"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	Size         int    `json:"size"`
	Token        string `json:"token"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
}
