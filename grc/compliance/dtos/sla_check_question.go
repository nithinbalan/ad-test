package dtos

type SLACheckDTO struct {
	ID                     uint    `json:"id"`
	SLARegisterID          uint    `json:"sla_register_id"`
	CreatedAt              string  `json:"created_at"`
	UpdatedAt              *string `json:"updated_at,omitempty"`
	CreatedUserFirstNameEn string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn  string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr  *string `json:"created_user_last_name_ar"`
	UpdatedUserFirstNameEn string  `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr string  `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn  string  `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr  string  `json:"updated_user_last_name_ar"`
}

type SLACheckQuestionFormDTO struct {
	QuestionGroupID      uint   `json:"question_group_id"`
	QuestionGroupTitleEn string `json:"question_group_title_en"`
	QuestionGroupTitleAr string `json:"question_group_title_ar"`
	QuestionTitleEn      string `json:"question_title_en" binding:"required"`
	QuestionTitleAr      string `json:"question_title_ar" binding:"required"`
}

type SLACheckUpdateQuestionGroupDTO struct {
	QuestionGroupTitleEn string `json:"question_group_title_en" binding:"required"`
	QuestionGroupTitleAr string `json:"question_group_title_ar" binding:"required"`
}

type SLACheckUpdateQuestionDTO struct {
	QuestionTitleEn string `json:"question_title_en" binding:"required"`
	QuestionTitleAr string `json:"question_title_ar" binding:"required"`
}
