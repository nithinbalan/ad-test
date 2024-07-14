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
	TitleEn string `json:"title_en" binding:"required"`
	TitleAr string `json:"title_ar" binding:"required"`
}

type SLACheckQuestionGroupDetailDTO struct {
	ID             uint       `json:"id"`
	TitleEn        string     `json:"title_en"`
	TitleAr        string     `json:"title_ar"`
	TotalCount     int        `json:"total_count"`
	CompletedCount int        `json:"completed_count"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}
