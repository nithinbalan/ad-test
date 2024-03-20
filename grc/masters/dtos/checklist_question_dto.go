package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
)

type ChecklistQuestionListDTO struct {
	ID               uint   `json:"id"`
	QuestionEn       string `json:"question_en"`
	QuestionAr       string `json:"question_ar"`
	ChecklistTitleEn string `json:"checklist_title_en"`
	ChecklistTitleAr string `json:"checklist_title_ar"`
	ChecklistID      uint   `json:"checklist_id"`
}

type ChecklistQuestionListResponseDTO struct {
	Items      []*ChecklistQuestionListDTO  `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type ChecklistQuestionFormDTO struct {
	ChecklistID       uint   `json:"checklist_id" binding:"required"`
	QuestionLibraryID uint   `json:"question_library_id"`
	QuestionEn        string `json:"question_en" binding:"required,max=255"`
	QuestionAr        string `json:"question_ar" binding:"required,max=255"`
}
