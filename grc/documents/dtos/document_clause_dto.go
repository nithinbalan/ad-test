package dtos

import (
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type DocumentClauseListDTO struct {
	ID                          uint                                      `json:"id"`
	ClauseNumber                *string                                   `json:"clause_number"`
	ClauseTitleEn               string                                    `json:"clause_title_en"`
	ClauseTitleAr               string                                    `json:"clause_title_ar"`
	Description                 *string                                   `json:"description"`
	CreatedUserID               uint                                      `json:"created_user_id"`
	CreatedUserFirstNameEn      string                                    `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string                                    `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string                                    `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string                                    `json:"created_user_last_name_ar"`
	CreatedUserInitial          string                                    `json:"created_user_initial"`
	CreatedUserColor            string                                    `json:"created_user_color"`
	CreatedUserEmail            string                                    `json:"created_user_email"`
	CreatedUserProfileImagePath *string                                   `json:"created_user_profile_image_path"`
	UpdatedUserID               *uint                                     `json:"updated_user_id"`
	UpdatedUserFirstNameEn      string                                    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      string                                    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       string                                    `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       string                                    `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          string                                    `json:"updated_user_initial"`
	UpdatedUserColor            string                                    `json:"updated_user_color"`
	UpdatedUserEmail            string                                    `json:"updated_user_email"`
	UpdatedUserProfileImagePath *string                                   `json:"updated_user_profile_image_path"`
	CreatedAt                   time.Time                                 `json:"created_at"`
	UpdatedAt                   *time.Time                                `json:"updated_at"`
	Controls                    []*organizationDTOs.ControlMinimalListDTO `json:"controls"`
}

type DocumentClauseFormDTO struct {
	ClauseNumber   *string `json:"clause_number"`
	ClauseTitleEn  string  `json:"clause_title_en" binding:"required"`
	ClauseTitleAr  string  `json:"clause_title_ar" binding:"required"`
	Description    *string `json:"description"`
	ParantClauseID *uint   `json:"parant_clause_id"`
	ControlIds     []*uint `json:"control_ids"`
}

type DocumentClauseListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
}

type DocumentClauseDetailsDTO struct {
	ID               uint                     `json:"id"`
	ClauseNumber     *string                  `json:"clause_number"`
	ClauseTitleEn    string                   `json:"clause_title_en"`
	ClauseTitleAr    string                   `json:"clause_title_ar"`
	Description      *string                  `json:"description"`
	ParantClause     *DocumentParantClauseDTO `json:"parant_clause"`
	DocumentControls []*DocumentControlDTO    `json:"document_controls"`
}

type DocumentParantClauseDTO struct {
	ID            uint    `json:"id"`
	ClauseNumber  *string `json:"clause_number"`
	ClauseTitleEn string  `json:"clause_title_en"`
	ClauseTitleAr string  `json:"clause_title_ar"`
	Description   *string `json:"description"`
}

type DocumentControlDTO struct {
	ID                      uint                                    `json:"id"`
	DocumentVersionClauseID *uint                                   `json:"document_version_clause_id"`
	Control                 *organizationDTOs.ControlMinimalListDTO `json:"controls"`
}
