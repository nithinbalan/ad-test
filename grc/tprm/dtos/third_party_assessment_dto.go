package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type ThirdPartyAssessmentListDTO struct {
	ID                          uint    `json:"id"`
	ReferenceNumber             string  `json:"reference_number"`
	Title                       string  `json:"title"`
	ThirdPartyRegisterID        uint    `json:"third_party_register_id"`
	ThirdPartyRegisterTitle     string  `json:"third_party_register_title"`
	QuestionSetID               uint    `json:"question_set_id"`
	QuestionSetTitleEn          string  `json:"question_set_title_en"`
	QuestionSetTitleAr          string  `json:"question_set_title_ar"`
	Message                     *string `json:"message"`
	DueDate                     *string `json:"due_date"`
	CreatedAt                   string  `json:"created_at"`
	UpdatedAt                   string  `json:"updated_at"`
	CreatedByID                 uint    `json:"created_by_id"`
	CreatedUserFirstNameEn      string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string  `json:"created_user_last_name_ar"`
	CreatedUserInitial          string  `json:"created_user_initial"`
	CreatedUserColor            string  `json:"created_user_color"`
	CreatedUserEmail            string  `json:"created_user_email"`
	CreatedUserProfileImagePath string  `json:"created_user_profile_image_path"`
	UpdatedByID                 *uint   `json:"updated_by_id"`
	UpdatedUserFirstNameEn      *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          *string `json:"updated_user_initial"`
	UpdatedUserColor            *string `json:"updated_user_color"`
	UpdatedUserEmail            *string `json:"updated_user_email"`
	UpdatedUserProfileImagePath *string `json:"updated_user_profile_image_path"`
}

type ThirdPartyAssessmentFormDTO struct {
	Title                string  `json:"title" binding:"required"`
	ThirdPartyRegisterID uint    `json:"third_party_register_id" binding:"required"`
	QuestionSetID        uint    `json:"question_set_id" binding:"required"`
	Message              *string `json:"message"`
	DueDate              *string `json:"due_date" binding:"required"`
}

type ThirdPartyAssessmentDetailDTO struct {
	ID                 uint                            `json:"id"`
	ReferenceNumber    string                          `json:"reference_number"`
	Title              string                          `json:"title"`
	ThirdPartyRegister ThirdPartyRegisterDetailDTO     `json:"third_party_register"`
	QuestionSet        QuestionSetDetailDTO            `json:"question_set"`
	Message            *string                         `json:"message"`
	DueDateString      *string                         `json:"due_date"`
	CreatedAt          time.Time                       `json:"created_at"`
	UpdatedAt          *time.Time                      `json:"updated_at"`
	CreatedBy          organizationDTOs.UserMinimalDTO `json:"created_by"`
	UpdatedBy          organizationDTOs.UserMinimalDTO `json:"updated_by"`
}
type ThirdPartyAssessmentListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	DateFrom    string
	DateTo      string
	Export      bool
}
type ThirdPartyAssessmentListResponseDTO struct {
	Items      []*ThirdPartyAssessmentListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO   `json:"pagination"`
}

type ThirdPartyAssessmentParsedDates struct {
	DueDate *time.Time `json:"due_date"`
}
