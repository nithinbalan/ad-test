package dtos

import (
	"time"
)

type PolicyClauseListDTO struct {
	ID                          uint                   `json:"id"`
	ClauseNumber                *string                `json:"clause_number"`
	ClauseTitleEn               string                 `json:"clause_title_en"`
	ClauseTitleAr               string                 `json:"clause_title_ar"`
	Description                 *string                `json:"description"`
	CreatedUserID               uint                   `json:"created_user_id"`
	CreatedUserFirstNameEn      string                 `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string                 `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string                 `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string                 `json:"created_user_last_name_ar"`
	CreatedUserInitial          string                 `json:"created_user_initial"`
	CreatedUserColor            string                 `json:"created_user_color"`
	CreatedUserEmail            string                 `json:"created_user_email"`
	CreatedUserProfileImagePath *string                `json:"created_user_profile_image_path"`
	UpdatedUserID               *uint                  `json:"updated_user_id"`
	UpdatedUserFirstNameEn      string                 `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      string                 `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       string                 `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       string                 `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          string                 `json:"updated_user_initial"`
	UpdatedUserColor            string                 `json:"updated_user_color"`
	UpdatedUserEmail            string                 `json:"updated_user_email"`
	UpdatedUserProfileImagePath *string                `json:"updated_user_profile_image_path"`
	CreatedAt                   time.Time              `json:"created_at"`
	UpdatedAt                   *time.Time             `json:"updated_at"`
	DocumentVersionClauseID     *uint                  `json:"document_version_clause_id"`
	Children                    []*PolicyClauseListDTO `json:"children" gorm:"-"`
	CommentsCount               int                    `json:"comments_count"`
}

type PolicyClauseFormDTO struct {
	ClauseTitleEn    string  `json:"clause_title_en" binding:"required"`
	ClauseTitleAr    string  `json:"clause_title_ar" binding:"required"`
	Description      *string `json:"description"`
	PreviousClauseID *uint   `json:"previous_clause_id"`
	ParentClauseID   *uint   `json:"parent_clause_id"`
}

type PolicyClauseListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
}

type PolicyClauseDetailsDTO struct {
	ID            uint                     `json:"id"`
	ClauseNumber  *string                  `json:"clause_number"`
	ClauseTitleEn string                   `json:"clause_title_en"`
	ClauseTitleAr string                   `json:"clause_title_ar"`
	Description   *string                  `json:"description"`
	ParentClause  *DocumentParentClauseDTO `json:"parent_clause"`
}

type DocumentParentClauseDTO struct {
	ID            uint    `json:"id"`
	ClauseNumber  *string `json:"clause_number"`
	ClauseTitleEn string  `json:"clause_title_en"`
	ClauseTitleAr string  `json:"clause_title_ar"`
	Description   *string `json:"description"`
}

type ClauseHistoryDTO struct {
	ID                          *uint     `json:"id"`
	DocumentVersionClauseID     uint      `json:"document_version_clause_id"`
	DocumentVersionClauseNumber string    `json:"document_version_clause_number"`
	OldValue                    *string   `json:"old_value"`
	NewValue                    *string   `json:"new_value"`
	ChangedField                *string   `json:"changedField"`
	CreatedUserID               uint      `json:"created_user_id"`
	CreatedUserFirstNameEn      string    `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string    `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string    `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string    `json:"created_user_last_name_ar"`
	CreatedUserInitial          string    `json:"created_user_initial"`
	CreatedUserColor            string    `json:"created_user_color"`
	CreatedUserEmail            string    `json:"created_user_email"`
	CreatedUserProfileImagePath *string   `json:"created_user_profile_image_path"`
	CreatedAt                   time.Time `json:"created_at"`
}
