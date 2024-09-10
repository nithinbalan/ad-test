package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type DocumentAuditLogListDTO struct {
	ID                            uint      `json:"id"`
	DocumentID                    uint      `json:"document_id"`
	DocumentVersionID             uint      `json:"document_version_id"`
	UserID                        uint      `json:"user_id"`
	CreatedBYID                   uint      `json:"created_by_id"`
	CreatedAt                     time.Time `json:"created_at"`
	DocumentAuditLogTypeTitleEn   string    `json:"document_audit_log_type_title_en"`
	DocumentAuditLogTypeTitleAr   string    `json:"document_audit_log_type_title_ar"`
	DocumentAuditLogTypeColorCode string    `json:"document_audit_log_type_color_code"`
	DocumentTitleEn               string    `json:"document_title_en"`
	DocumentTitleAr               string    `json:"document_title_ar"`
	DocumentReferenceNumber       string    `json:"document_reference_number"`
	DocumentURL                   string    `json:"document_url"`
	DocumentName                  string    `json:"document_name"`
	DocumentSize                  int       `json:"document_size"`
	DocumentToken                 string    `json:"document_token"`
	DocumentOriginalName          string    `json:"document_original_name"`
	DocumentType                  string    `json:"document_type"`
	DocumentVersion               string    `json:"document_version"`
	UserFirstNameEn               string    `json:"user_first_name_en"`
	UserFirstNameAr               string    `json:"user_first_name_ar"`
	UserLastNameEn                string    `json:"user_last_name_en"`
	UserLastNameAr                string    `json:"user_last_name_ar"`
	UserEmail                     string    `json:"user_email"`
	UserInitial                   string    `json:"user_initial"`
	UserColor                     string    `json:"user_color"`
	UserProfileImagePath          string    `json:"user_profile_image_path"`
	CreatedUserID                 uint      `json:"created_user_id"`
	CreatedUserFirstNameEn        string    `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr        string    `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn         string    `json:"created_user_last_name_en"`
	CreatedUserLastNameAr         string    `json:"created_user_last_name_ar"`
	CreatedUserInitial            string    `json:"created_user_initial"`
	CreatedUserColor              string    `json:"created_user_color"`
	CreatedUserEmail              string    `json:"created_user_email"`
	CreatedUserProfileImagePath   string    `json:"created_user_profile_image_path"`
}

type DocumentAuditLogListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}

type DocumentAuditLogListResponseDTO struct {
	Items      []*DocumentAuditLogListDTO   `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type DocumentAuditLogFormDTO struct {
	ID                     uint      `json:"id,omitempty"`
	DocumentID             uint      `json:"document_id" binding:"required"`
	DocumentVersionID      uint      `json:"document_version_id" binding:"required"`
	UserID                 uint      `json:"user_id" binding:"required"`
	CreatedByID            uint      `json:"created_by_id" binding:"required"`
	CreatedAt              time.Time `json:"created_at" binding:"required"`
	DocumentAuditLogTypeID uint      `json:"document_audit_log_type_id" binding:"required"`
}
