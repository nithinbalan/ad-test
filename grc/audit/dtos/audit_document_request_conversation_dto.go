package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
)

type AuditDocumentRequestConversationListDTO struct {
	ID                            uint                                           `json:"id"`
	AuditDocumentRequestID        uint                                           `json:"audit_document_request_id"`
	RequestedBYID                 uint    										 `json:"requested_by_id"`
	Comment                       string                                         `json:"comment"`
	CommentedAt                   string                                         `json:"commented_at"`
	CommentedUserID               uint                                           `json:"commented_user_id"`
	CommentedUserFirstNameEn      string                                         `json:"commented_user_first_name_en"`
	CommentedUserFirstNameAr      string                                         `json:"commented_user_first_name_ar"`
	CommentedUserLastNameEn       string                                         `json:"commented_user_last_name_en"`
	CommentedUserLastNameAr       string                                         `json:"commented_user_last_name_ar"`
	CommentedUserInitial          string                                         `json:"commented_user_initial"`
	CommentedUserColor            string                                         `json:"commented_user_color"`
	CommentedUserEmail            string                                         `json:"commented_user_email"`
	CommentedUserProfileImagePath string                                         `json:"commented_user_profile_image_path"`
	Documents                     []*AuditDocumentRequestConversationDocumentDTO `json:"documents" gorm:"-"`
}

type AuditDocumentRequestConversationFormDTO struct {
	Comment   string                        `json:"comment" binding:"required,max=500"`
	Documents []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
	IsClose   *bool                         `json:"is_close,omitempty"`
}


type AuditDocumentRequestConversationListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	CommentedAt string
}

type AuditDocumentRequestConversationDocumentDTO struct {
	ID           uint   `json:"id"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	Size         int    `json:"size"`
	Token        string `json:"token"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
}
