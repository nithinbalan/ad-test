package dtos

type ComplianceActionPlanStatusUpdateListDTO struct {
	ID                        uint    `json:"id"`
	ComplianceActionPlanTitle string  `json:"compliance_action_plan_title"`
	ActionPlanStatusTitleEn   string  `json:"action_plan_status_title_en"`
	ActionPlanStatusTitleAr   string  `json:"action_plan_status_title_ar"`
	ActionPlanStatusColorCode string  `json:"action_plan_status_color_code"`
	Comment                   string  `json:"comment"`
	Percentage                string  `json:"percentage"`
	CreatedAt                 string  `json:"created_at"`
	UpdatedAt                 string  `json:"updated_at"`
	CreatedUserFirstNameEn    string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr    string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn     string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr     *string `json:"created_user_last_name_ar"`
	UpdatedUserFirstNameEn    string  `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr    string  `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn     string  `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr     string  `json:"updated_user_last_name_ar"`
}

type ComplianceActionPlanStatusUpdateFormDTO struct {
	ActionPlanStatusID uint   `json:"action_plan_status_id" binding:"required"`
	Comment            string `json:"comment" binding:"required"`
	Percentage         *uint  `json:"percentage" binding:"required"`
}
