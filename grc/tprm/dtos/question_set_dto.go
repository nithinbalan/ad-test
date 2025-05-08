package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type QuestionSetListDTO struct {
	ID                          uint   `json:"id"`
	TitleEn                     string `json:"title_en"`
	TitleAr                     string `json:"title_ar"`
	Description                 string `json:"description"`
	CreatedAt                   string `json:"created_at"`
	UpdatedAt                   string `json:"updated_at"`
	CreatedByID                 uint   `json:"created_by_id"`
	CreatedUserFirstNameEn      string `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string `json:"created_user_last_name_ar"`
	CreatedUserInitial          string `json:"created_user_initial"`
	CreatedUserColor            string `json:"created_user_color"`
	CreatedUserEmail            string `json:"created_user_email"`
	CreatedUserProfileImagePath string `json:"created_user_profile_image_path"`
	UpdatedByID                 *uint  `json:"updated_by_id"`
	UpdatedUserFirstNameEn      string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          string `json:"updated_user_initial"`
	UpdatedUserColor            string `json:"updated_user_color"`
	UpdatedUserEmail            string `json:"updated_user_email"`
	UpdatedUserProfileImagePath string `json:"updated_user_profile_image_path"`
	TotalQuestions              int    `json:"total_questions"`
}

type QuestionSetFormDTO struct {
	TitleEn     string  `json:"title_en" binding:"required"`
	TitleAr     string  `json:"title_ar" binding:"required"`
	Description *string `json:"description"`
}

type QuestionSetQuestionFormDTO struct {
	ThirdPartyQuestionLibraryID uint   `json:"third_party_question_library_id" binding:"required"`
	QuestionTitleEn             string `json:"question_title_en" binding:"required"`
	QuestionTitleAr             string `json:"question_title_ar" binding:"required"`
}

type QuestionSetListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	DateFrom    string
	DateTo      string
	Export      bool
}

type QuestionSetListResponseDTO struct {
	Items      []*QuestionSetListDTO        `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type QuestionSetQuestionDetailDTO struct {
	ID                          uint   `json:"id"`
	QuestionSetID               uint   `json:"question_set_id"`
	ThirdPartyQuestionLibraryID uint   `json:"third_party_question_library_id"`
	QuestionTitleEn             string `json:"question_title_en"`
	QuestionTitleAr             string `json:"question_title_ar"`

	// Related data from library
	ReferenceNumber              string `json:"reference_number"`
	Mandatory                    bool   `json:"mandatory"`
	ThirdPartyQuestionCategoryID uint   `json:"third_party_question_category_id"`
	ThirdPartyQuestionTypeID     uint   `json:"third_party_question_type_id"`

	// Category & Type Titles
	QuestionCategoryTitleEn string `json:"question_category_title_en"`
	QuestionCategoryTitleAr string `json:"question_category_title_ar"`
	QuestionTypeTitleEn     string `json:"question_type_title_en"`
	QuestionTypeTitleAr     string `json:"question_type_title_ar"`
	QuestionTypeType        string `json:"question_type_type"`

	// Answer Options
	AnswerOptions []ThirdPartyQuestionLibraryAnswerOptionDetailDTO `json:"answer_options"`
}

type QuestionSetDetailDTO struct {
	ID             uint                          `json:"id"`
	TitleEn        string                        `json:"title_en"`
	TitleAr        string                        `json:"title_ar"`
	Description    *string                       `json:"description"`
	CreatedAt      time.Time                     `json:"created_at"`
	UpdatedAt      *time.Time                    `json:"updated_at"`
	Categories     []QuestionSetCategoryGroupDTO `json:"categories"`
	TotalQuestions int                           `json:"total_questions"`
}

type QuestionSetCategoryGroupDTO struct {
	CategoryTitleEn string                         `json:"category_title_en"`
	CategoryTitleAr string                         `json:"category_title_ar"`
	Questions       []QuestionSetQuestionDetailDTO `json:"questions"`
}

type QuestionSetAddQuestionsDTO struct {
	QuestionLibraryIDs []uint `json:"question_library_ids" binding:"required"` // required non-empty array
}
