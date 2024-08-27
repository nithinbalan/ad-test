package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	mastersDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type AuditWorkpaperListDTO struct {
	ID                          uint    `json:"id"`
	AuditID                     uint    `json:"audit_id"`
	SourceID                    uint    `json:"source_id"`
	SourceTitleEn               string  `json:"source_title_en"`
	SourceTitleAr               string  `json:"source_title_ar"`
	Title                       string  `json:"title"`
	Description                 *string `json:"description"`
	URL                         *string `json:"url"`
	Name                        *string `json:"name"`
	Size                        *int    `json:"size"`
	Token                       *string `json:"token"`
	OriginalName                *string `json:"original_name"`
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
	UpdatedUserFirstNameEn      string  `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      string  `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       string  `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       string  `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          string  `json:"updated_user_initial"`
	UpdatedUserColor            string  `json:"updated_user_color"`
	UpdatedUserEmail            string  `json:"updated_user_email"`
	UpdatedUserProfileImagePath string  `json:"updated_user_profile_image_path"`
	CreatedAt                   string  `json:"created_at"`
	UpdatedAt                   *string `json:"updated_at"`
}

type AuditWorkpaperFormDTO struct {
	SourceID    uint                          `json:"source_id" binding:"required"`
	Title       string                        `json:"title" binding:"required,max=255"`
	Description *string                       `json:"description"`
	Documents   []*commonDTOs.DocumentFormDTO `json:"documents" binding:"required"`
}

type AuditWorkpaperDetailDTO struct {
	ID            uint                             `json:"id"`
	AuditID       uint                             `json:"audit_id"`
	SourceID      uint                             `json:"source_d"`
	Audit         AuditDetailDTO                   `json:"audit"`
	Source        mastersDTOs.SourceDTO            `json:"source"`
	Title         string                           `json:"title"`
	Description   *string                          `json:"description"`
	URL           *string                          `json:"url"`
	Name          *string                          `json:"name"`
	Size          *int                             `json:"size"`
	Token         *string                          `json:"token"`
	OriginalName  *string                          `json:"original_name"`
	CreatedAt     time.Time                        `json:"created_at"`
	UpdatedAt     *time.Time                       `json:"updated_at"`
	CreatedByUser organizationDTOs.UserMinimalDTO  `json:"created_by_user"`
	UpdatedByUser *organizationDTOs.UserMinimalDTO `json:"updated_by_user"`
}

type AuditWorkpaperListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}
type AuditWorkpaperListResponseDTO struct {
	Items      []*AuditWorkpaperListDTO     `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}
