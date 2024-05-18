package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

// ProcessListDTO defines the data transfer object for listing processes
type ProcessListDTO struct {
	ID                       uint       `json:"id"`
	ProcessGroupTitleEn      string     `json:"process_group_title_en"`
	ProcessGroupTitleAr      string     `json:"process_group_title_ar"`
	OwnerUnitTitleEn         *string    `json:"owner_unit_title_en"`
	OwnerUnitTitleAr         *string    `json:"owner_unit_title_ar"`
	OwnerFirstNameEn         *string    `json:"owner_first_name_en"`
	OwnerFirstNameAr         *string    `json:"owner_first_name_ar"`
	OwnerLastNameEn          *string    `json:"owner_last_name_en"`
	OwnerLastNameAr          *string    `json:"owner_last_name_ar"`
	LifecycleStatusTitleEn   string     `json:"lifecycle_status_title_en"`
	LifecycleStatusTitleAr   string     `json:"lifecycle_status_title_ar"`
	LifecycleStatusColorCode string     `json:"lifecycle_status_color_code"`
	ReferenceNumber          string     `json:"reference_number"`
	TitleEn                  string     `json:"title_en"`
	TitleAr                  string     `json:"title_ar"`
	Description              *string    `json:"description"`
	CreatedUserFirstNameEn   string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr   string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn    string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr    string     `json:"created_user_last_name_ar"`
	UpdatedUserFirstNameEn   *string    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr   *string    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn    *string    `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr    *string    `json:"updated_user_last_name_ar"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                *time.Time `json:"updated_at"`
}

// ProcessListQueryParamsDTO defines the filters and search parameters for querying processes
type ProcessListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}

// ProcessListResponseDTO encapsulates the paginated response for a list of processes
type ProcessListResponseDTO struct {
	Items      []*ProcessListDTO            `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

// ProcessFormDTO defines the structure for creating or updating a process
type ProcessFormDTO struct {
	ProcessGroupID    uint    `json:"process_group_id" binding:"required"`
	OwnerUnitID       uint    `json:"owner_unit_id" binding:"required"`
	OwnerID           uint    `json:"owner_id" binding:"required"`
	LifecycleStatusID uint    `json:"lifecycle_status_id" binding:"required"`
	TitleEn           string  `json:"title_en" binding:"required"`
	TitleAr           string  `json:"title_ar" binding:"required"`
	Description       *string `json:"description" binding:"omitempty"`
}

// ProcessDetailDTO defines the detailed view of a process, typically used in detailed screens or reports
type ProcessDetailDTO struct {
	ID                       uint       `json:"id"`
	ProcessGroupTitleEn      string     `json:"process_group_title_en"`
	ProcessGroupTitleAr      string     `json:"process_group_title_ar"`
	OwnerUnitTitleEn         *string    `json:"owner_unit_title_en"`
	OwnerUnitTitleAr         *string    `json:"owner_unit_title_ar"`
	OwnerFirstNameEn         *string    `json:"owner_first_name_en"`
	OwnerFirstNameAr         *string    `json:"owner_first_name_ar"`
	OwnerLastNameEn          *string    `json:"owner_last_name_en"`
	OwnerLastNameAr          *string    `json:"owner_last_name_ar"`
	LifecycleStatusTitleEn   string     `json:"lifecycle_status_title_en"`
	LifecycleStatusTitleAr   string     `json:"lifecycle_status_title_ar"`
	LifecycleStatusColorCode string     `json:"lifecycle_status_color_code"`
	ReferenceNumber          string     `json:"reference_number"`
	TitleEn                  string     `json:"title_en"`
	TitleAr                  string     `json:"title_ar"`
	Description              *string    `json:"description"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                *time.Time `json:"updated_at"`
}
