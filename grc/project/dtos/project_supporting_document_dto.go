package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type ProjectSupportingDocumentListDTO struct {
	ID                          uint                                    `json:"id"`
	ProjectID                   uint                                    `json:"project_id"`
	TitleEn                     *string                                 `json:"title_en"`
	TitleAr                     *string                                 `json:"title_ar"`
	Documents                   []*ProjectSupportingDocumentDocumentDTO `json:"documents" gorm:"-"`
	CreatedUserID               *uint                                   `json:"created_user_id"`
	CreatedUserFirstNameEn      *string                                 `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      *string                                 `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       *string                                 `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       *string                                 `json:"created_user_last_name_ar"`
	CreatedUserInitial          *string                                 `json:"created_user_initial"`
	CreatedUserColor            *string                                 `json:"created_user_color"`
	CreatedUserEmail            *string                                 `json:"created_user_email"`
	CreatedUserProfileImagePath *string                                 `json:"created_user_profile_image_path"`
	UpdatedUserID               *uint                                   `json:"updated_user_id"`
	UpdatedUserFirstNameEn      *string                                 `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      *string                                 `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       *string                                 `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       *string                                 `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          *string                                 `json:"updated_user_initial"`
	UpdatedUserColor            *string                                 `json:"updated_user_color"`
	UpdatedUserEmail            *string                                 `json:"updated_user_email"`
	UpdatedUserProfileImagePath *string                                 `json:"updated_user_profile_image_path"`
	CreatedAt                   string                                  `json:"created_at"`
	UpdatedAt                   *string                                 `json:"updated_at"`
}

type ProjectSupportingDocumentDetailDTO struct {
	ID            uint                                    `json:"id"`
	ProjectID     uint                                    `json:"project_id"`
	TitleEn       *string                                 `json:"title_en"`
	TitleAr       *string                                 `json:"title_ar"`
	CreatedByUser *organizationDTOs.UserMinimalDTO        `json:"created_by_user"`
	UpdatedByUser *organizationDTOs.UserMinimalDTO        `json:"updated_by_user"`
	CreatedAt     time.Time                               `json:"created_at"`
	UpdatedAt     *time.Time                              `json:"updated_at"`
	Documents     []*ProjectSupportingDocumentDocumentDTO `json:"documents"`
}

type ProjectSupportingDocumentDocumentDTO struct {
	ID           uint    `json:"id"`
	URL          string  `json:"url"`
	Name         string  `json:"name"`
	OriginalName *string `json:"original_name"`
	Size         int     `json:"size"`
	Token        string  `json:"token"`
	Type         string  `json:"type"`
}

type ProjectSupportingDocumentFormDTO struct {
	TitleEn   *string                       `json:"title_en"`
	TitleAr   *string                       `json:"title_ar"`
	Documents []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
}

type ProjectSupportingDocumentListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	DateFrom    string
	DateTo      string
}

type ProjectSupportingDocumentListResponseDTO struct {
	Items       []*ProjectSupportingDocumentListDTO `json:"data"`
	Pagination commonDTOs.PaginationInfoDTO        `json:"pagination"`
}
