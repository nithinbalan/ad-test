package dtos

import (
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type ProjectClosureListDTO struct {
	ID                           uint    `json:"id"`
	ProjectID                    uint    `json:"project_id"`
	OverallProjectOutcome        *string `json:"overall_project_outcome"`
	ClosureReason                *string `json:"closure_reason"`
	ClosureComments              *string `json:"closure_comments"`
	LessonLearned                *string `json:"lesson_learned"`
	ApproverID                   uint    `json:"approver_id"`
	IsApproved                   bool    `json:"is_approved"`
	ApproverUserFirstNameEn      string  `json:"approver_user_first_name_en"`
	ApproverUserFirstNameAr      string  `json:"approver_user_first_name_ar"`
	ApproverUserLastNameEn       string  `json:"approver_user_last_name_en"`
	ApproverUserLastNameAr       string  `json:"approver_user_last_name_ar"`
	ApproverUserInitial          string  `json:"approver_user_initial"`
	ApproverUserColor            string  `json:"approver_user_color"`
	ApproverUserEmail            string  `json:"approver_user_email"`
	ApproverUserProfileImagePath string  `json:"approver_user_profile_image_path"`
	CreatedUserID                uint    `json:"created_user_id"`
	CreatedUserFirstNameEn       string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr       string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn        string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr        string  `json:"created_user_last_name_ar"`
	CreatedUserInitial           string  `json:"created_user_initial"`
	CreatedUserColor             string  `json:"created_user_color"`
	CreatedUserEmail             string  `json:"created_user_email"`
	CreatedUserProfileImagePath  string  `json:"created_user_profile_image_path"`
	CreatedAt                    string  `json:"created_at"`
	UpdatedAt                    *string `json:"updated_at"`
}

type ProjectClosureDetailDTO struct {
	ID                    uint                             `json:"id"`
	ProjectID             uint                             `json:"project_id"`
	OverallProjectOutcome *string                          `json:"overall_project_outcome"`
	ClosureReason         *string                          `json:"closure_reason"`
	ClosureComments       *string                          `json:"closure_comments"`
	LessonLearned         *string                          `json:"lesson_learned"`
	ApproverID            uint                             `json:"approver_id"`
	IsApproved            bool                             `json:"is_approved"`
	ApproverUser          *organizationDTOs.UserMinimalDTO `json:"approver_user"`
	CreatedByUser         *organizationDTOs.UserMinimalDTO `json:"created_by"`
	UpdatedByUser         *organizationDTOs.UserMinimalDTO `json:"updated_by"`
	CreatedAt             time.Time                        `json:"created_at"`
	UpdatedAt             *time.Time                       `json:"updated_at"`
}

type ProjectClosureFormDTO struct {
	OverallProjectOutcome *string `json:"overall_project_outcome"`
	ClosureReason         *string `json:"closure_reason"`
	ClosureComments       *string `json:"closure_comments"`
	LessonLearned         *string `json:"lesson_learned"`
	ApproverID            uint    `json:"approver_id" binding:"required"`
}
