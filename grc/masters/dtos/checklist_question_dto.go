package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
)

type ChecklistQuestionListDTO struct {
	ID                   uint   `json:"id"`
	QuestionEn           string `json:"question_en"`
	QuestionAr           string `json:"question_ar"`
	ChecklistTitleEn     string `json:"checklist_title_en"`
	ChecklistTitleAr     string `json:"checklist_title_ar"`
	ChecklistID          uint   `json:"checklist_id"`
	QuestionGroupTitleEn string `json:"cquestion_group_title_en"`
	QuestionGroupTitleAr string `json:"cquestion_group_title_ar"`
	QuestionGroupID      uint   `json:"cquestion_group_id"`
}

type ChecklistQuestionListResponseDTO struct {
	Items      []*ChecklistQuestionListDTO  `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type ChecklistQuestionFormDTO struct {
	ChecklistID     uint   `json:"checklist_id" binding:"required"`
	QuestionGroupID uint   `json:"question_group_id" binding:"required"`
	QuestionEn      string `json:"question_en" binding:"required,max=255"`
	QuestionAr      string `json:"question_ar" binding:"required,max=255"`
}
