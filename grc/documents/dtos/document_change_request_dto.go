package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

// List DTO
type DocumentChangeRequestListDTO struct {
	ID                            uint    `json:"id"`
	ReferenceNumber               string  `json:"reference_number"`
	DocumentID                    uint    `json:"document_id"`
	DocumentTitleEn               string  `json:"document_title_en"`
	DocumentTitleAr               string  `json:"document_title_ar"`
	DocumentChangeRequestTypeID   uint    `json:"document_change_request_type_id"`
	DocumentChangeRequestTypeEn   string  `json:"document_change_request_type_en"`
	DocumentChangeRequestTypeAr   string  `json:"document_change_request_type_ar"`
	DocumentChangeRequestStatusID uint    `json:"document_change_request_status_id"`
	DocumentChangeRequestStatusEn string  `json:"document_change_request_status_en"`
	DocumentChangeRequestStatusAr string  `json:"document_change_request_status_ar"`
	ReasonForChange               *string `json:"reason_for_change"`
	Version                       float64 `json:"version"`

	CreatedByUserID           uint   `json:"created_by_user_id"`
	CreatedByUserFirstNameEn  string `json:"created_by_user_first_name_en"`
	CreatedByUserFirstNameAr  string `json:"created_by_user_first_name_ar"`
	CreatedByUserLastNameEn   string `json:"created_by_user_last_name_en"`
	CreatedByUserLastNameAr   string `json:"created_by_user_last_name_ar"`
	CreatedByUserEmail        string `json:"created_by_user_email"`
	CreatedByUserProfileImage string `json:"created_by_user_profile_image"`

	UpdatedByUserID           *uint   `json:"updated_by_user_id"`
	UpdatedByUserFirstNameEn  *string `json:"updated_by_user_first_name_en"`
	UpdatedByUserFirstNameAr  *string `json:"updated_by_user_first_name_ar"`
	UpdatedByUserLastNameEn   *string `json:"updated_by_user_last_name_en"`
	UpdatedByUserLastNameAr   *string `json:"updated_by_user_last_name_ar"`
	UpdatedByUserEmail        *string `json:"updated_by_user_email"`
	UpdatedByUserProfileImage *string `json:"updated_by_user_profile_image"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// Details DTO
type DocumentChangeRequestDetailsDTO struct {
	ID                            uint    `json:"id"`
	ReferenceNumber               string  `json:"reference_number"`
	DocumentID                    uint    `json:"document_id"`
	DocumentTitleEn               string  `json:"document_title_en"`
	DocumentTitleAr               string  `json:"document_title_ar"`
	DocumentChangeRequestTypeID   uint    `json:"document_change_request_type_id"`
	DocumentChangeRequestTypeEn   string  `json:"document_change_request_type_en"`
	DocumentChangeRequestTypeAr   string  `json:"document_change_request_type_ar"`
	DocumentChangeRequestStatusID uint    `json:"document_change_request_status_id"`
	DocumentChangeRequestStatusEn string  `json:"document_change_request_status_en"`
	DocumentChangeRequestStatusAr string  `json:"document_change_request_status_ar"`
	ReasonForChange               *string `json:"reason_for_change"`
	Version                       float64 `json:"version"`

	SubmittedByUser organizationDTOs.UserMinimalDTO  `json:"submitted_by_user"`
	PublishedByUser *organizationDTOs.UserMinimalDTO `json:"published_by_user,omitempty"`

	Documents []*DocumentChangeRequestDocumentDTO `json:"documents"`

	CreatedByUser organizationDTOs.UserMinimalDTO  `json:"created_by_user"`
	UpdatedByUser *organizationDTOs.UserMinimalDTO `json:"updated_by_user,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// Document DTO
type DocumentChangeRequestDocumentDTO struct {
	ID           *uint   `json:"id,omitempty"`
	URL          *string `json:"url"`
	Name         string  `json:"file_name"`
	OriginalName string  `json:"original_name"`
	Type         string  `json:"type"`
	Size         int     `json:"size"`
}

// Query Params DTO
type DocumentChangeRequestQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}

// List Response DTO
type DocumentChangeRequestListResponseDTO struct {
	Items      []*DocumentChangeRequestListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO    `json:"pagination"`
}

type DocumentChangeRequestFormDTO struct {
	DocumentChangeRequestTypeID     uint                       `json:"document_change_request_type_id" binding:"required"`
	DocumentID                      uint                       `json:"document_id" binding:"required"`
	DocumentChangeRequestInitiateID *uint                      `json:"document_change_request_initiate_id"`
	ReasonForChange                 *string                    `json:"reason_for_change"`
	Document                        commonDTOs.DocumentFormDTO `json:"document" binding:"required"`
}
