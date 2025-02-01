package dtos

import (
	"time"

	organizationDTOs "jethur/grc/organization/dtos"
)

type ComplianceActionPlanCommentListDTO struct {
	ID                          uint   `json:"id"`
	Comment                     string `json:"comment" binding:"required"`
	ComplianceActionPlanID      uint   `json:"compliance_action_plan_id"`
	CreatedAt                   string `json:"created_at"`
	CreatedBYID                 uint   `json:"created_by_id"`
	CreatedUserFirstNameEn      string `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string `json:"created_user_last_name_ar"`
	CreatedUserInitial          string `json:"created_user_initial"`
	CreatedUserColor            string `json:"created_user_color"`
	CreatedUserEmail            string `json:"created_user_email"`
	CreatedUserProfileImagePath string `json:"created_user_profile_image_path"`
}

type ComplianceActionPlanCommentFormDTO struct {
	Comment string `json:"comment" binding:"required,max=1000"`
}

type ComplianceActionPlanCommentDetailDTO struct {
	ID                     uint                            `json:"id"`
	Comment                string                          `json:"comment"`
	ComplianceActionPlanID uint                            `json:"compliance_action_plan_id"`
	CreatedAt              time.Time                       `json:"created_at"`
	CreatedByUser          organizationDTOs.UserMinimalDTO `json:"created_by"`
}
