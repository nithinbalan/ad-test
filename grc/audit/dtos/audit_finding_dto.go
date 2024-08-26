package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type AuditFindingListDTO struct {
	ID                          uint    `json:"id"`
	AuditID                     uint    `json:"audit_id"`
	AuditableItemID             uint    `json:"auditable_item_id"`
	AuditFindingCategoryID      uint    `json:"audit_finding_category_id"`
	SeverityLevelID             uint    `json:"severity_level_id"`
	Recommendations             *string `json:"recommendations,omitempty"`
	AuditFindingStatusID        uint    `json:"audit_finding_status_id"`
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
	AuditableItemID        uint    `json:"auditable_item_id" binding:"required"`
	AuditFindingCategoryID uint    `json:"audit_finding_category_id" binding:"required"`
	SeverityLevelID        uint    `json:"severity_level_id" binding:"required"`
	Recommendations        *string `json:"recommendations,omitempty"`
}

type AuditFindingDetailDTO struct {
	ID                     uint                             `json:"id"`
	AuditID                uint                             `json:"audit_id"`
	AuditableItemID        uint                             `json:"auditable_item_id"`
	AuditFindingCategoryID uint                             `json:"audit_finding_category_id"`
	SeverityLevelID        uint                             `json:"severity_level_id"`
	Recommendations        *string                          `json:"recommendations,omitempty"`
	AuditFindingStatusID   uint                             `json:"audit_finding_status_id"`
	CreatedAt              time.Time                        `json:"created_at"`
	UpdatedAt              *time.Time                       `json:"updated_at"`
	CreatedByUser          organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser          *organizationDTOs.UserMinimalDTO `json:"updated_by"`
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
