package dtos

type SLACheckAnswerListDTO struct {
	ID                      uint    `json:"id"`
	SLACheckQuestionTitleEn string  `json:"sla_check_question_title_en"`
	SLACheckQuestionTitleAr string  `json:"sla_check_question_title_ar"`
	ComplianceStatusTitleEn string  `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr string  `json:"compliance_status_title_ar"`
	ComplianceStatusColor   string  `json:"compliance_status_color_code"`
	Comment                 string  `json:"comment"`
	CreatedAt               string  `json:"created_at"`
	UpdatedAt               *string `json:"updated_at,omitempty"`
	CreatedUserFirstNameEn  string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr  string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn   string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr   *string `json:"created_user_last_name_ar,omitempty"`
	UpdatedUserFirstNameEn  string  `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr  string  `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn   string  `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr   string  `json:"updated_user_last_name_ar"`
}

type SLACheckAnswerFormDTO struct {
	ComplianceStatusID uint   `json:"compliance_status_id" binding:"required"`
	Comment            string `json:"comment" binding:"required"`
}

type AssessmentSummaryDTO struct {
	TotalQuestions int          `json:"totalQuestions"`
	TotalAnswers   int          `json:"totalAnswers"`
	AnswerCounts   map[uint]int `json:"answerCounts"`
}
