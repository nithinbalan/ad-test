package dtos

import (
	"time"
)

// ProcessListDTO defines the data transfer object for listing processes
type UserDocumentListDTO struct {
	ID                          uint       `json:"id"`
	UserID                      uint       `json:"user_id"`
	UserDocumentTypeID          string     `json:"user_document_type_id"`
	UserDocumentTypeTitleEn     string     `json:"user_document_type_title_en"`
	UserDocumentTypeTitleAr     string     `json:"user_document_type_title_ar"`
	UserDocumentSubTypeID       string     `json:"user_document_sub_type_id"`
	UserDocumentSubTypeTitleEn  string     `json:"user_document_sub_type_title_en"`
	UserDocumentSubTypeTitleAr  string     `json:"user_document_sub_type_title_ar"`
	ExpiryDate                  string     `json:"expiry_date"`
	IssuedBy                    string     `json:"issued_by"`
	Remarks                     *string    `json:"remarks"`
	URL                         string     `json:"url"`
	Name                        string     `json:"name"`
	OriginalName                string     `json:"original_name"`
	Type                        string     `json:"type"`
	Size                        int        `json:"size"`
	Token                       string     `json:"token"`
	CreatedUserID               uint       `json:"created_user_id"`
	CreatedUserFirstNameEn      string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string     `json:"created_user_last_name_ar"`
	CreatedUserInitial          string     `json:"created_user_initial"`
	CreatedUserColor            string     `json:"created_user_color"`
	CreatedUserEmail            string     `json:"created_user_email"`
	CreatedUserProfileImagePath string     `json:"created_user_profile_image_path"`
	UpdatedUserID               *uint      `json:"updated_user_id"`
	UpdatedUserFirstNameEn      *string    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      *string    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       *string    `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       *string    `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          string     `json:"updated_user_initial"`
	UpdatedUserColor            string     `json:"updated_user_color"`
	UpdatedUserEmail            string     `json:"updated_user_email"`
	UpdatedUserProfileImagePath string     `json:"updated_user_profile_image_path"`
	CreatedAt                   time.Time  `json:"created_at"`
	UpdatedAt                   *time.Time `json:"updated_at"`
}

type UserDocumentListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}

type UserDocumentFormDTO struct {
	UserDocumentTypeID    uint    `json:"user_document_type_id" binding:"required"`
	UserDocumentSubTypeID uint    `json:"user_document_sub_type_id" binding:"required"`
	IssuedBy              string  `json:"issued_by" binding:"required"`
	ExpiryDate            string  `json:"expiry_date" binding:"required"`
	Remarks               *string `json:"remarks"`
	FileName              string  `json:"file_name" binding:"required"`
	OriginalName          string  `json:"original_name" binding:"required"`
	Type                  string  `json:"type" binding:"required"`
	Size                  int     `json:"size" binding:"required"`
}

type UserDocumentDetailDTO struct {
	ID                  uint                         `json:"id"`
	UserDocumentType    UserDocumentTypeDetailDTO    `json:"user_document_type"`
	UserDocumentSubType UserDocumentSubTypeDetailDTO `json:"user_document_sub_type"`
	ExpiryDateString    string                       `json:"expiry_date"`
	IssuedBy            string                       `json:"issued_by"`
	Remarks             *string                      `json:"remarks"`
	URL                 string                       `json:"url"`
	Name                string                       `json:"name"`
	OriginalName        string                       `json:"original_name"`
	Token               string                       `json:"token"`
	Type                string                       `json:"type"`
	Size                int                          `json:"size"`
	CreatedByUser       UserMinimalDTO               `json:"created_by"`
	UpdatedByUser       *UserMinimalDTO              `json:"updated_by"`
}

type UserDocumentParsedDates struct {
	ExpiryDate *time.Time
}
