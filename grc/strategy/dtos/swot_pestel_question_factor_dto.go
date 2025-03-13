package dtos

import (
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type SwotPestelQuestionFactorListDTO struct {
	ID                          uint    `json:"id"`
	SwotPestelQuestionID        uint    `json:"swot_pestel_question_id"`
	SwotPestelCategoryID        uint    `json:"swot_pestel_category_id"`
	Description                 string  `json:"description"`
	SwotPestelCategoryTitleEn   string  `json:"swot_pestel_category_title_en"`
	SwotPestelCategoryTitleAr   string  `json:"swot_pestel_category_title_ar"`
	SwotPestelCategoryColorCode string  `json:"swot_pestel_category_color_code"`
	SwotPestelCategoryType      string  `json:"swot_pestel_category_type"`
	SwotPestelQuestionTitleEn   string  `json:"swot_pestel_question_title_en"`
	SwotPestelQuestionTitleAr   string  `json:"swot_pestel_question_title_ar"`
	CreatedUserID               uint    `json:"created_user_id"`
	CreatedUserFirstNameEn      string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string  `json:"created_user_last_name_ar"`
	CreatedUserInitial          string  `json:"created_user_initial"`
	CreatedUserColor            string  `json:"created_user_color"`
	CreatedUserEmail            string  `json:"created_user_email"`
	CreatedUserProfileImagePath string  `json:"created_user_profile_image_path"`
	UpdatedUserID               *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn      *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          *string `json:"updated_user_initial"`
	UpdatedUserColor            *string `json:"updated_user_color"`
	UpdatedUserEmail            *string `json:"updated_user_email"`
	UpdatedUserProfileImagePath *string `json:"updated_user_profile_image_path"`
}

type SwotPestelQuestionFactorFormDTO struct {
	SwotPestelCategoryID uint    `json:"swot_pestel_category_id" binding:"required"`
	Description          *string `json:"description" binding:"required"`
}

type SwotPestelQuestionFactorListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	Export      bool
}

type SwotPestelQuestionFactorDetailDTO struct {
	ID                   uint                             `json:"id"`
	SwotPestelQuestionID uint                             `json:"swot_pestel_question_id"`
	SwotPestelQuestion   SwotPestelQuestionDTO            `json:"swot_pestel_question"`
	SwotPestelCategoryID uint                             `json:"swot_pestel_category_id"`
	SwotPestelCategory   SwotPestelCategoryDTO            `json:"swot_pestel_category"`
	Description          *string                          `json:"description"`
	CreatedAt            time.Time                        `json:"created_at"`
	UpdatedAt            *time.Time                       `json:"updated_at"`
	CreatedByUser        organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser        *organizationDTOs.UserMinimalDTO `json:"updated_by"`
}
