package dtos

type SLACheckQuestionGroupListDTO struct {
	ID             uint   `json:"id"`
	SLACheckID     uint   `json:"sla_check_id"`
	TitleEn        string `json:"title_en"`
	TitleAr        string `json:"title_ar"`
	TotalCount     int    `json:"total_count"`
	CompletedCount int    `json:"completed_count"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
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
