package dtos

import "time"

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

type SLACheckQuestionDetailDTO struct {
	ID                      uint                     `json:"id"`
	SLACheckQuestionGroupID uint                     `json:"sla_check_question_group_id"`
	TitleEn                 string                   `json:"title_en"`
	TitleAr                 string                   `json:"title_ar"`
	Answer                  *SLACheckAnswerDetailDTO `json:"answer"`
	CreatedAt               time.Time                `json:"created_at"`
	UpdatedAt               *time.Time               `json:"updated_at"`
}

type SLACheckQuestionGroupFormDTO struct {
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
}

type SLACheckQuestionFormDTO struct {
	QuestionGroupID uint   `json:"question_group_id"`
	TitleEn         string `json:"title_en" binding:"required"`
	TitleAr         string `json:"title_ar" binding:"required"`
}

type SLACheckUpdateQuestionGroupDTO struct {
	QuestionGroupTitleEn string `json:"question_group_title_en" binding:"required"`
	QuestionGroupTitleAr string `json:"question_group_title_ar" binding:"required"`
}

type SLACheckUpdateQuestionDTO struct {
	QuestionTitleEn string `json:"question_title_en" binding:"required"`
	QuestionTitleAr string `json:"question_title_ar" binding:"required"`
}
