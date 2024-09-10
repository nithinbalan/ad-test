package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"

	"time"
)

type AuditFindingListDTO struct {
	ID                          uint    `json:"id"`
	AuditID                     uint    `json:"audit_id"`
	AuditTitle                  string  `json:"audit_title"`
	ReferenceNumber             string  `json:"reference_number"`
	Title                       string  `json:"title"`
	Description                 *string `json:"description"`
	Recommendations             *string `json:"recommendations"`
	AuditableItemTitleEn        string  `json:"auditable_item_title_en"`
	AuditableItemTitleAr        string  `json:"auditable_item_title_ar"`
	AuditFindingCategoryTitleEn string  `json:"audit_finding_category_title_en"`
	AuditFindingCategoryTitleAr string  `json:"audit_finding_category_title_ar"`
	SeverityLevelTitleEn        string  `json:"severity_level_title_en"`
	SeverityLevelTitleAr        string  `json:"severity_level_title_ar"`
	SeverityLevelColorCode      string  `json:"severity_level_color_code"`
	AuditFindingStatusTitleEn   string  `json:"audit_finding_status_title_en"`
	AuditFindingStatusTitleAr   string  `json:"audit_finding_status_title_ar"`
	AuditFindingStatusColorCode string  `json:"audit_finding_status_color_code"`
	CreatedAt                   string  `json:"created_at"`
	UpdatedAt                   *string `json:"updated_at"`
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

type AuditFindingFormDTO struct {
	AuditID                uint                           `json:"audit_id" binding:"required"`
	Title                  string                         `json:"title" binding:"required,max=255"`
	Description            *string                        `json:"description"`
	AuditableItemID        uint                           `json:"auditable_item_id" binding:"required"`
	AuditFindingCategoryID uint                           `json:"audit_finding_category_id" binding:"required"`
	SeverityLevelID        uint                           `json:"severity_level_id" binding:"required"`
	Recommendations        *string                        `json:"recommendations,omitempty"`
	Documents              []*AuditFindingDocumentFormDTO `json:"documents" binding:"required"`
}

type AuditFindingDocumentFormDTO struct {
	ID            *uint  `json:"id,omitempty"`
	URL           string `json:"url"`
	FileName      string `json:"file_name"`
	OriginalName  string `json:"original_name"`
	Type          string `json:"type"`
	Size          int    `json:"size"`
	IsNewFile     *bool  `json:"is_new_file,omitempty"`
	IsDeletedFile *bool  `json:"is_deleted_file,omitempty"`
}

type AuditFindingDetailDTO struct {
	ID                     uint                             `json:"id"`
	ReferenceNumber        string                           `json:"reference_number"`
	Title                  string                           `json:"title"`
	Description            *string                          `json:"description"`
	AuditID                uint                             `json:"audit_id"`
	AuditableItemID        uint                             `json:"auditable_item_id"`
	AuditFindingCategoryID uint                             `json:"audit_finding_category_id"`
	SeverityLevelID        uint                             `json:"severity_level_id"`
	Recommendations        *string                          `json:"recommendations,omitempty"`
	AuditFindingStatusID   uint                             `json:"audit_finding_status_id"`
	AuditableItem          AuditableItemDetailDTO           `json:"auditable_item"`
	AuditFindingCategory   AuditFindingCategoryDTO          `json:"audit_finding_category"`
	SeverityLevel          masterDTOs.SeverityLevelDTO      `json:"severity_level"`
	AuditFindingStatus     AuditFindingStatusDTO            `json:"audit_finding_status"`
	CreatedAt              time.Time                        `json:"created_at"`
	UpdatedAt              *time.Time                       `json:"updated_at"`
	CreatedByUser          organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser          *organizationDTOs.UserMinimalDTO `json:"updated_by"`
	Documents              []*AuditFindingDocumentDTO       `json:"documents"`
}

type AuditFindingDocumentDTO struct {
	ID           *uint   `json:"id"`
	URL          *string `json:"url"`
	Name         string  `json:"name"`
	Token        string  `json:"token"`
	OriginalName string  `json:"original_name"`
	Type         string  `json:"type"`
	Size         int     `json:"size"`
}

type AuditFindingListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}

type AuditFindingListResponseDTO struct {
	Items      []*AuditFindingListDTO       `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}
