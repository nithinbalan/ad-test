package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type DocumentChangeRequestInitiateListDTO struct {
	ID                                       uint    `json:"id"`
	DocumentsID                              uint    `json:"documents_id"`
	DocumentTitle                            string  `json:"document_title"`
	DocumentChangeRequestInitiateTypeID      uint    `json:"document_change_request_initiate_type_id"`
	DocumentChangeRequestInitiateTypeEn      string  `json:"document_change_request_initiate_type_en"`
	DocumentChangeRequestInitiateTypeAr      string  `json:"document_change_request_initiate_type_ar"`
	ReasonForChange                          *string `json:"reason_for_change"`
	DocumentChangeRequestInitiateStatusID    uint    `json:"document_change_request_initiate_status_id"`
	DocumentChangeRequestInitiateStatusEn    string  `json:"document_change_request_initiate_status_en"`
	DocumentChangeRequestInitiateStatusAr    string  `json:"document_change_request_initiate_status_ar"`
	DocumentChangeRequestInitiateStatusColor string  `json:"document_change_request_initiate_status_color"`

	CreatedByUserID               uint   `json:"created_by_user_id"`
	CreatedByUserFirstNameEn      string `json:"created_by_user_first_name_en"`
	CreatedByUserFirstNameAr      string `json:"created_by_user_first_name_ar"`
	CreatedByUserLastNameEn       string `json:"created_by_user_last_name_en"`
	CreatedByUserLastNameAr       string `json:"created_by_user_last_name_ar"`
	CreatedByUserEmail            string `json:"created_by_user_email"`
	CreatedByUserProfileImagePath string `json:"created_by_user_profile_image_path"`

	UpdatedByUserID               *uint   `json:"updated_by_user_id"`
	UpdatedByUserFirstNameEn      *string `json:"updated_by_user_first_name_en"`
	UpdatedByUserFirstNameAr      *string `json:"updated_by_user_first_name_ar"`
	UpdatedByUserLastNameEn       *string `json:"updated_by_user_last_name_en"`
	UpdatedByUserLastNameAr       *string `json:"updated_by_user_last_name_ar"`
	UpdatedByUserEmail            *string `json:"updated_by_user_email"`
	UpdatedByUserProfileImagePath *string `json:"updated_by_user_profile_image_path"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type DocumentChangeRequestInitiateDetailsDTO struct {
	ID                                       uint    `json:"id"`
	DocumentsID                              uint    `json:"documents_id"`
	DocumentTitle                            string  `json:"document_title"`
	DocumentChangeRequestInitiateTypeID      uint    `json:"document_change_request_initiate_type_id"`
	DocumentChangeRequestInitiateTypeEn      string  `json:"document_change_request_initiate_type_en"`
	DocumentChangeRequestInitiateTypeAr      string  `json:"document_change_request_initiate_type_ar"`
	ReasonForChange                          *string `json:"reason_for_change"`
	DocumentChangeRequestInitiateStatusID    uint    `json:"document_change_request_initiate_status_id"`
	DocumentChangeRequestInitiateStatusEn    string  `json:"document_change_request_initiate_status_en"`
	DocumentChangeRequestInitiateStatusAr    string  `json:"document_change_request_initiate_status_ar"`
	DocumentChangeRequestInitiateStatusColor string  `json:"document_change_request_initiate_status_color"`

	CreatedByUser organizationDTOs.UserMinimalDTO  `json:"created_by_user"`
	UpdatedByUser *organizationDTOs.UserMinimalDTO `json:"updated_by_user"`

	Documents []*DocumentChangeRequestInitiateDocumentDTO `json:"documents"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type DocumentChangeRequestInitiateFormDTO struct {
	DocumentChangeRequestInitiateTypeID uint                          `json:"document_change_request_initiate_type_id" binding:"required"`
	ReasonForChange                     *string                       `json:"reason_for_change" binding:"omitempty"`
	Documents                           []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
}

type DocumentChangeRequestInitiateDocumentDTO struct {
	ID           *uint   `json:"id,omitempty"`
	URL          *string `json:"url"`
	Name         string  `json:"file_name"`
	OriginalName string  `json:"original_name"`
	Type         string  `json:"type"`
	Size         int     `json:"size"`
}

type DocumentChangeRequestInitiateListResponseDTO struct {
	Items      []*DocumentChangeRequestInitiateListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO            `json:"pagination"`
}

type DocumentChangeRequestInitiateQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}
