package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type ThirdPartyAssessmentListDTO struct {
	ID                                  uint    `json:"id"`
	ReferenceNumber                     string  `json:"reference_number"`
	Title                               string  `json:"title"`
	ThirdPartyRegisterID                uint    `json:"third_party_register_id"`
	ThirdPartyRegisterTitle             string  `json:"third_party_register_title"`
	QuestionSetID                       uint    `json:"question_set_id"`
	QuestionSetTitleEn                  string  `json:"question_set_title_en"`
	QuestionSetTitleAr                  string  `json:"question_set_title_ar"`
	Message                             *string `json:"message"`
	DueDate                             *string `json:"due_date"`
	UID                                 *string `json:"uid"`
	Link                                *string `json:"link"`
	ThirdPartyAssessmentStatusTitleEn   string  `json:"third_party_assessment_status_title_en"`
	ThirdPartyAssessmentStatusTitleAr   string  `json:"third_party_assessment_status_title_ar"`
	ThirdPartyAssessmentStatusType      string  `json:"third_party_assessment_status_type"`
	ThirdPartyAssessmentStatusColorCode string  `json:"third_party_assessment_status_color_code"`
	RiskLevelTitleEn                    *string `json:"risk_level_title_en"`
	RiskLevelTitleAr                    *string `json:"risk_level_title_ar"`
	RiskLevelColorCode                  *string `json:"risk_level_color_code"`
	RiskLevelType                       *string `json:"risk_level_type"`
	CreatedAt                           string  `json:"created_at"`
	UpdatedAt                           string  `json:"updated_at"`
	CreatedByID                         uint    `json:"created_by_id"`
	CreatedUserFirstNameEn              string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr              string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn               string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr               string  `json:"created_user_last_name_ar"`
	CreatedUserInitial                  string  `json:"created_user_initial"`
	CreatedUserColor                    string  `json:"created_user_color"`
	CreatedUserEmail                    string  `json:"created_user_email"`
	CreatedUserProfileImagePath         string  `json:"created_user_profile_image_path"`
	UpdatedByID                         *uint   `json:"updated_by_id"`
	UpdatedUserFirstNameEn              *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr              *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn               *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr               *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial                  *string `json:"updated_user_initial"`
	UpdatedUserColor                    *string `json:"updated_user_color"`
	UpdatedUserEmail                    *string `json:"updated_user_email"`
	UpdatedUserProfileImagePath         *string `json:"updated_user_profile_image_path"`
}

type ThirdPartyAssessmentFormDTO struct {
	Title                string  `json:"title" binding:"required"`
	ThirdPartyRegisterID uint    `json:"third_party_register_id" binding:"required"`
	QuestionSetID        uint    `json:"question_set_id" binding:"required"`
	Message              *string `json:"message"`
	DueDate              *string `json:"due_date" binding:"required"`
	IsSentEmail          bool    `json:"is_sent_email"`
}

type ThirdPartyAssessmentDetailDTO struct {
	ID                         uint                             `json:"id"`
	ReferenceNumber            string                           `json:"reference_number"`
	Title                      string                           `json:"title"`
	ThirdPartyRegister         ThirdPartyRegisterDetailDTO      `json:"third_party_register"`
	QuestionSet                QuestionSetDetailDTO             `json:"question_set"`
	Message                    *string                          `json:"message"`
	DueDateString              *string                          `json:"due_date"`
	UID                        *string                          `json:"uid"`
	Link                       *string                          `json:"link"`
	ThirdPartyAssessmentStatus ThirdPartyAssessmentStatusDTO    `json:"third_party_assessment_status"`
	NextReviewLevel            *int                             `json:"next_review_level"`
	SubmittedAt                *time.Time                       `json:"submitted_at"`
	Sections                   []TPRMAssessmentSectionDTO       `json:"categories"`
	IsSend                     bool                             `json:"is_send"`
	CreatedAt                  time.Time                        `json:"created_at"`
	UpdatedAt                  *time.Time                       `json:"updated_at"`
	CreatedByUser              organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser              *organizationDTOs.UserMinimalDTO `json:"updated_by"`
	SubmittedByUser            *organizationDTOs.UserMinimalDTO `json:"submitted_by"`
	ImportedByUser             *organizationDTOs.UserMinimalDTO `json:"imported_by"`
	ImportedAt                 time.Time                        `json:"imported_at"`
	SentByUser                 *organizationDTOs.UserMinimalDTO `json:"send_by"`
	SentAt                     time.Time                        `json:"send_at"`
	Score                      int                              `json:"score"`
	RiskLevel                  organizationDTOs.RiskLevelDTO    `json:"risk_level"`
}

type TPRMAssessmentSectionDTO struct {
	ID            uint                              `json:"id"`
	TitleEn       string                            `json:"category_title_en"`
	TitleAr       string                            `json:"category_title_ar"`
	Order         int                               `json:"order"`
	Questions     []TPRMAssessmentQuestionDetailDTO `json:"questions"`
	QuestionCount int                               `json:"question_count"`
}

type TPRMAssessmentQuestionDetailDTO struct {
	ID              uint   `json:"id"`
	ReferenceNumber string `json:"reference_number"`
	QuestionTitleEn string `json:"question_title_en"`
	QuestionTitleAr string `json:"question_title_ar"`
	Order           int    `json:"order"`
	Mandatory       bool   `json:"mandatory"`

	QuestionTypeTitleEn string                                           `json:"question_type_title_en"`
	QuestionTypeTitleAr string                                           `json:"question_type_title_ar"`
	QuestionTypeType    string                                           `json:"question_type_type"`
	AnswerOptions       []ThirdPartyQuestionLibraryAnswerOptionDetailDTO `json:"answer_options"`

	Answer TPRMAssessmentAnswerDTO `json:"answer"`
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

type TPRMAssessmentAnswerFormDTO struct {
	Question  string `json:"question"`
	Mandatory string `json:"mandatory"`
	Answer    string `json:"answer"`
	Comment   string `json:"comment"`
}

type TPRMAssessmentAnswerDTO struct {
	ID            uint                             `json:"id"`
	Comment       string                           `json:"comment"`
	AnswerOptions []*TPRMAssessmentAnswerOptionDTO `json:"answer_options"`
}

type TPRMAssessmentAnswerOptionDTO struct {
	ID           uint                                           `json:"id"`
	AnswerOption ThirdPartyQuestionLibraryAnswerOptionDetailDTO `json:"answer_option"`
}

type AnswerSummaryItemDTO struct {
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	Count     int64  `json:"count"`
	ColorCode string `json:"color_code"`
}

type ThirdPartyAssessmentAnswerFormDTO struct {
	OptionId uint    `json:"option_id" binding:"required"`
	Comment  *string `json:"comment"`
}

type AnswerDocumentFormDTO struct {
	Documents []*commonDTOs.DocumentFormDTO `json:"documents" binding:"required"`
}
type AnswerDocumentListDTO struct {
	ID           uint       `json:"id"`
	URL          string     `json:"url"`
	Name         string     `json:"name"`
	OriginalName string     `json:"original_name"`
	Type         string     `json:"type"`
	Size         int        `json:"size"`
	Token        string     `json:"token"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}
