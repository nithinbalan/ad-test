package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type ThirdPartyQuestionLibraryListDTO struct {
	ID                                uint   `json:"id"`
	ReferenceNumber                   string `json:"reference_number"`
	QuestionTitleEn                   string `json:"question_title_en"`
	QuestionTitleAr                   string `json:"question_title_ar"`
	ThirdPartyQuestionCategoryID      uint   `json:"third_party_question_category_id"`
	ThirdPartyQuestionCategoryTitleEn string `json:"third_party_question_category_title_en"`
	ThirdPartyQuestionCategoryTitleAr string `json:"third_party_question_category_title_ar"`
	Mandatory                         bool   `json:"mandatory"`
	ThirdPartyQuestionTypeID          uint   `json:"third_party_question_type_id"`
	ThirdPartyQuestionTypeTitleEn     string `json:"third_party_question_type_title_en"`
	ThirdPartyQuestionTypeTitleAr     string `json:"third_party_question_type_title_ar"`
	CreatedAt                         string `json:"created_at"`
	UpdatedAt                         string `json:"updated_at"`
}

type ThirdPartyQuestionLibraryFormDTO struct {
	QuestionTitleEn              string                                         `json:"question_title_en" binding:"required"`
	QuestionTitleAr              string                                         `json:"question_title_ar" binding:"required"`
	ThirdPartyQuestionCategoryID uint                                           `json:"third_party_question_category_id" binding:"required"`
	Mandatory                    bool                                           `json:"mandatory"`
	ThirdPartyQuestionTypeID     uint                                           `json:"third_party_question_type_id" binding:"required"`
	AnswerOptions                []ThirdPartyQuestionLibraryAnswerOptionFormDTO `json:"answer_options"`
}

type ThirdPartyQuestionLibraryDetailDTO struct {
	ID                         uint                          `json:"id"`
	ReferenceNumber            string                        `json:"reference_number"`
	QuestionTitleEn            string                        `json:"question_title_en"`
	QuestionTitleAr            string                        `json:"question_title_ar"`
	ThirdPartyQuestionCategory ThirdPartyQuestionCategoryDTO `json:"third_party_question_category"`
	ThirdPartyQuestionType     ThirdPartyQuestionTypeDTO     `json:"third_party_question_type"`
	Mandatory                  bool                          `json:"mandatory"`
	CreatedAt                  time.Time                     `json:"created_at"`
	UpdatedAt                  *time.Time                    `json:"updated_at"`
}

type ThirdPartyQuestionLibraryListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}

type ThirdPartyQuestionLibraryListResponseDTO struct {
	Items      []*ThirdPartyQuestionLibraryListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO        `json:"pagination"`
}

type ThirdPartyQuestionLibraryAnswerOptionFormDTO struct {
	AnswerOptionTitleEn *string `json:"answer_option_title_en"`
	AnswerOptionTitleAr *string `json:"answer_option_title_ar"`
	RiskLevelID         *uint   `json:"risk_level_id"`
	RangeFrom           *int    `json:"range_from"`
	RangeTo             *int    `json:"range_to"`
}
