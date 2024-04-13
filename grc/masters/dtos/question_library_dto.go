package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
)

type QuestionLibraryFormDTO struct {
	QuestionGroupID uint   `json:"question_group_id" binding:"required"`
	QuestionEn      string `json:"question_en" binding:"required,max=500"`
	QuestionAr      string `json:"question_ar" binding:"required,max=500"`
	Status          string `json:"status" binding:"required,max=10,oneof=active inactive"`
}

type QuestionLibraryListDTO struct {
	ID                   uint   `json:"id"`
	QuestionGroupID      uint   `json:"question_group_id"`
	QuestionEn           string `json:"question_en"`
	QuestionAr           string `json:"question_ar"`
	Status               string `json:"status"`
	QuestionGroupTitleAr string `json:"question_group_title_ar"`
	QuestionGroupTitleEn string `json:"question_group_title_en"`
}

type QuestionLibraryListResponseDTO struct {
	Items      []*QuestionLibraryListDTO    `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}
