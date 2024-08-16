package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
)

type MeetingActionPlanStatusUpdateListDTO struct {
	ID                        uint    `json:"id"`
	MeetingActionPlanTitle    string  `json:"meeting_action_plan_title"`
	ActionPlanStatusTitleEn   string  `json:"action_plan_status_title_en"`
	ActionPlanStatusTitleAr   string  `json:"action_plan_status_title_ar"`
	ActionPlanStatusColorCode string  `json:"action_plan_status_color_code"`
	Comment                   string  `json:"comment"`
	Percentage                string  `json:"percentage"`
	CreatedAt                 string  `json:"created_at"`
	UpdatedAt                 string  `json:"updated_at"`
	CreatedUserID             uint    `json:"created_user_id"`
	CreatedUserFirstNameEn    string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr    string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn     string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr     *string `json:"created_user_last_name_ar"`
	UpdatedUserFirstNameEn    string  `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr    string  `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn     string  `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr     string  `json:"updated_user_last_name_ar"`
}

type MeetingActionPlanStatusUpdateFormDTO struct {
	ActionPlanStatusID uint                          `json:"action_plan_status_id" binding:"required"`
	Comment            string                        `json:"comment" binding:"required"`
	Percentage         uint                          `json:"percentage" binding:"required"`
	Documents          []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
}
