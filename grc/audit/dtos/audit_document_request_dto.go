package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"time"
)

type AuditDocumentRequestListDTO struct {
	ID                                  uint    `json:"id"`
	AuditID                             uint    `json:"audit_id"`
	AuditDocumentRequestStatusID        uint    `json:"audit_document_request_status_id"`
	AuditDocumentRequestStatusTitleEn   string  `json:"audit_document_request_status_title_en"`
	AuditDocumentRequestStatusTitleAr   string  `json:"audit_document_request_status_title_ar"`
	AuditDocumentRequestStatusType      string  `json:"audit_document_request_status_type"`
	AuditDocumentRequestStatusColorCode string  `json:"audit_document_request_status_color_code"`
	ResponsibleUnitID                   uint    `json:"responsible_unit_id"`
	ResponsibleUnitTitleEn              string  `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr              string  `json:"responsible_unit_title_ar"`
	ResponsibleUserID                   uint    `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn          string  `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr          string  `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn           string  `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr           string  `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial              string  `json:"responsible_user_initial"`
	ResponsibleUserColor                string  `json:"responsible_user_color"`
	ResponsibleUserEmail                string  `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath     string  `json:"responsible_user_profile_image_path"`
	RequestedUserID                     uint    `json:"requested_user_id"`
	RequestedUserFirstNameEn            string  `json:"requested_user_first_name_en"`
	RequestedUserFirstNameAr            string  `json:"requested_user_first_name_ar"`
	RequestedUserLastNameEn             string  `json:"requested_user_last_name_en"`
	RequestedUserLastNameAr             string  `json:"requested_user_last_name_ar"`
	RequestedUserInitial                string  `json:"requested_user_initial"`
	RequestedUserColor                  string  `json:"requested_user_color"`
	RequestedUserEmail                  string  `json:"requested_user_email"`
	RequestedUserProfileImagePath       string  `json:"requested_user_profile_image_path"`
	Title                               string  `json:"title"`
	Description                         string  `json:"description"`
	RequestedAt                         *string `json:"requested_at"`
	CreatedAt                           string  `json:"created_at"`
	UpdatedAt                           *string `json:"updated_at"`
	CreatedUserID                       uint    `json:"created_user_id"`
	CreatedUserFirstNameEn              string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr              string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn               string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr               string  `json:"created_user_last_name_ar"`
	CreatedUserInitial                  string  `json:"created_user_initial"`
	CreatedUserColor                    string  `json:"created_user_color"`
	CreatedUserEmail                    string  `json:"created_user_email"`
	CreatedUserProfileImagePath         string  `json:"created_user_profile_image_path"`
	UpdatedUserID                       *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn              *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr              *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn               *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr               *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial                  *string `json:"updated_user_initial"`
	UpdatedUserColor                    *string `json:"updated_user_color"`
	UpdatedUserEmail                    *string `json:"updated_user_email"`
	UpdatedUserProfileImagePath         *string `json:"updated_user_profile_image_path"`
}

type AuditDocumentRequestFormDTO struct {
	Title            string  `json:"title" binding:"required,max=255"`
	Description      *string `json:"description"`
	AuditableItemIds []*uint `json:"auditable_item_ids"`
}

type AuditDocumentRequestDetailDTO struct {
	ID                         uint                                          `json:"id"`
	Audit                      AuditDetailDTO                                `json:"audit"`
	AuditDocumentRequestStatus AuditDocumentRequestStatusDTO                 `json:"audit_document_request_status"`
	Title                      string                                        `json:"title"`
	Description                *string                                       `json:"description"`
	AuditableItems             []*AuditDocumentRequestAuditableItemDetailDTO `json:"auditable_items"`
	RequestedAtString          string                                        `json:"requested_at"`
	CreatedAt                  time.Time                                     `json:"created_at"`
	UpdatedAt                  *time.Time                                    `json:"updated_at,omitempty"`
	ResponsibleUser            *organizationDTOs.UserMinimalDTO              `json:"responsible_user"`
	ResponsibleUnit            settingsDTOs.OrganizationUnitDTO              `json:"responsible_unit"`
	RequestedByUser            *organizationDTOs.UserMinimalDTO              `json:"requested_by_user"`
	CreatedByUser              organizationDTOs.UserMinimalDTO               `json:"created_by_user"`
	UpdatedByUser              *organizationDTOs.UserMinimalDTO              `json:"updated_by_user,omitempty"`
	ConversationCount          int                                           `json:"conversation_count"`
}

type AuditDocumentRequestListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	RequestedAt string
}

type AuditDocumentRequestListResponseDTO struct {
	Items      []*AuditDocumentRequestListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO   `json:"pagination"`
}

type AuditDocumentRequestParsedDates struct {
	RequestedAt *time.Time
}

type AuditDocumentRequestAuditableItemDTO struct {
	AuditableItemID        uint `json:"auditable_item_id"`
	AuditDocumentRequestID uint `json:"audit_document_request_id"`
	CreatedBYID            uint `json:"created_by_id"`
}

type ResponsibleUserGroupDTO struct {
	ResponsibleUnitID uint                                    `json:"responsible_unit_id"`
	ResponsibleUserID uint                                    `json:"responsible_user_id"`
	AuditableItems    []*AuditDocumentRequestAuditableItemDTO `json:"auditable_items"`
}

type AuditDocumentRequestAuditableItemDetailDTO struct {
	ID                     *uint                  `json:"id"`
	AuditDocumentRequestID uint                   `json:"audit_document_request_id"`
	AuditableItemID        uint                   `json:"auditable_item_id"`
	AuditableItem          AuditableItemDetailDTO `json:"auditable_item"`
}
