package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	settingsDTOs "jethur/grc/settings/dtos"

	"time"
)

// ProcessListDTO defines the data transfer object for listing processes
type ProcessListDTO struct {
	ID                          uint       `json:"id"`
	ProcessGroupTitleEn         string     `json:"process_group_title_en"`
	ProcessGroupTitleAr         string     `json:"process_group_title_ar"`
	OwnerUnitTitleEn            *string    `json:"owner_unit_title_en"`
	OwnerUnitTitleAr            *string    `json:"owner_unit_title_ar"`
	OwnerID                     *uint      `json:"owner_id"`
	OwnerFirstNameEn            *string    `json:"owner_first_name_en"`
	OwnerFirstNameAr            *string    `json:"owner_first_name_ar"`
	OwnerLastNameEn             *string    `json:"owner_last_name_en"`
	OwnerLastNameAr             *string    `json:"owner_last_name_ar"`
	OwnerInitial                string     `json:"owner_initial"`
	OwnerColor                  string     `json:"owner_color"`
	OwnerEmail                  string     `json:"owner_email"`
	OwnerProfileImagePath       string     `json:"owner_profile_image_path"`
	LifecycleStatusTitleEn      string     `json:"lifecycle_status_title_en"`
	LifecycleStatusTitleAr      string     `json:"lifecycle_status_title_ar"`
	LifecycleStatusColorCode    string     `json:"lifecycle_status_color_code"`
	ReferenceNumber             string     `json:"reference_number"`
	TitleEn                     string     `json:"title_en"`
	TitleAr                     string     `json:"title_ar"`
	Description                 *string    `json:"description"`
	CreatedUserID               uint       `json:"created_user_id"`
	CreatedUserFirstNameEn      string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string     `json:"created_user_last_name_ar"`
	CreatedUserInitial          string     `json:"created_user_owner_initial"`
	CreatedUserColor            string     `json:"created_user_owner_color"`
	CreatedUserEmail            string     `json:"created_user_owner_email"`
	CreatedUserProfileImagePath string     `json:"created_user_owner_profile_image_path"`
	UpdatedUserID               *uint      `json:"updated_user_id"`
	UpdatedUserFirstNameEn      *string    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      *string    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       *string    `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       *string    `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          string     `json:"updated_user_owner_initial"`
	UpdatedUserColor            string     `json:"updated_user_owner_color"`
	UpdatedUserEmail            string     `json:"updated_user_owner_email"`
	UpdatedUserProfileImagePath string     `json:"updated_user_owner_profile_image_path"`
	CreatedAt                   time.Time  `json:"created_at"`
	UpdatedAt                   *time.Time `json:"updated_at"`
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
	ProcessGroupID uint    `json:"process_group_id" binding:"required"`
	OwnerID        uint    `json:"owner_id" binding:"required"`
	TitleEn        string  `json:"title_en" binding:"required"`
	TitleAr        string  `json:"title_ar" binding:"required"`
	Description    *string `json:"description" binding:"omitempty"`
	ControlIds     []*uint `json:"control_ids"`
	DocumentIds    []*uint `json:"document_ids"`
}

type ProcessDetailDTO struct {
	ID              uint                              `json:"id"`
	ReferenceNumber string                            `json:"reference_number"`
	TitleEn         string                            `json:"title_en"`
	TitleAr         string                            `json:"title_ar"`
	Description     *string                           `json:"description"`
	ProcessGroup    *ProcessGroupDTO                  `json:"process_group"`
	OwnerUnit       *settingsDTOs.OrganizationUnitDTO `json:"owner_unit"`
	Owner           *UserMinimalDTO                   `json:"owner"`
	CreatedByUser   UserMinimalDTO                    `json:"created_by"`
	UpdatedByUser   *UserMinimalDTO                   `json:"updated_by"`
	LifecycleStatus masterDTOs.LifecycleStatusDTO     `json:"status"`
	CreatedAt       time.Time                         `json:"created_at"`
	UpdatedAt       *time.Time                        `json:"updated_at"`
	Controls        []*ProcessControlDTO              `json:"controls"`
	Documents       []*ProcessAssociatedDocumentDTO   `json:"documents"`
}

type ProcessControlDTO struct {
	ID        *uint            `json:"id"`
	ProcessID uint             `json:"process_id"`
	ControlID uint             `json:"control_id"`
	Control   ControlDetailDTO `json:"control"`
}

type ProcessAssociatedDocumentDTO struct {
	ID         *uint              `json:"id"`
	ProcessID  uint               `json:"process_id"`
	DocumentID uint               `json:"document_id"`
	Document   DocumentDetailsDTO `json:"document"`
}

type DocumentDetailsDTO struct {
	ID            int     `json:"id"`
	TitleEn       string  `json:"title_en"`
	TitleAr       string  `json:"title_ar"`
	DescriptionEn *string `json:"description_en"`
	DescriptionAr *string `json:"description_ar"`
	IsLocked      bool    `json:"is_locked"`
	Acknowledged  *bool   `json:"acknowledged,omitempty"`
}
