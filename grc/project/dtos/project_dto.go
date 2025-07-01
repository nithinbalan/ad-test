package dtos

import (
	"time"

	organizationDTOs "jethur/grc/organization/dtos"
)

type ProjectListDTO struct {
	ID                              uint   `json:"id"`
	TitleEn                         string `json:"title_en"`
	TitleAr                         string `json:"title_ar"`
	DescriptionEn                   string `json:"description_en"`
	DescriptionAr                   string `json:"description_ar"`
	ProjectTypeTitleEn              string `json:"project_type_title_en"`
	ProjectTypeTitleAr              string `json:"project_type_title_ar"`
	ProjectStatusTitleEn            string `json:"project_status_title_en"`
	ProjectStatusTitleAr            string `json:"project_status_title_ar"`
	StartDate                       string `json:"start_date"`
	EndDate                         string `json:"end_date"`
	ResponsibleUserID               uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string `json:"responsible_user_initial"`
	ResponsibleUserColor            string `json:"responsible_user_color"`
	ResponsibleUserEmail            string `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string `json:"responsible_user_profile_image_path"`
	CreatedUserID                   uint   `json:"created_user_id"`
	CreatedUserFirstNameEn          string `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr          string `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn           string `json:"created_user_last_name_en"`
	CreatedUserLastNameAr           string `json:"created_user_last_name_ar"`
	CreatedUserInitial              string `json:"created_user_initial"`
	CreatedUserColor                string `json:"created_user_color"`
	CreatedUserEmail                string `json:"created_user_email"`
	CreatedUserProfileImagePath     string `json:"created_user_profile_image_path"`
	UpdatedUserID                   *uint  `json:"updated_user_id"`
	UpdatedUserFirstNameEn          string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr          string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn           string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr           string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial              string `json:"updated_user_initial"`
	UpdatedUserColor                string `json:"updated_user_color"`
	UpdatedUserEmail                string `json:"updated_user_email"`
	UpdatedUserProfileImagePath     string `json:"updated_user_profile_image_path"`
}

type ProjectPaginationInfoDTO struct {
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
	TotalPages  int   `json:"total_pages"`
	TotalCount  int64 `json:"total_count"`
}

type ProjectListResponseDTO struct {
	Items      []*ProjectListDTO        `json:"items"`
	Pagination ProjectPaginationInfoDTO `json:"pagination"`
}

type ProjectListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}

type ProjectFormDTO struct {
	TitleEn           string `json:"title_en" binding:"required,max=255"`
	TitleAr           string `json:"title_ar" binding:"required,max=255"`
	DescriptionEn     string `json:"description_en"`
	DescriptionAr     string `json:"description_ar"`
	ResponsibleUserID uint   `json:"responsible_user_id" binding:"required"`
	StartDate         string `json:"start_date" binding:"required,datetime=2006-01-02"`
	EndDate           string `json:"end_date" binding:"required,datetime=2006-01-02"`
	ProjectTypeID     uint   `json:"project_type_id" binding:"required"`
	ProjectStatusID   uint   `json:"project_status_id" binding:"required"`
}

type ProjectDetailDTO struct {
	ID              uint                             `json:"id"`
	TitleEn         string                           `json:"title_en"`
	TitleAr         string                           `json:"title_ar"`
	DescriptionEn   string                           `json:"description_en"`
	DescriptionAr   string                           `json:"description_ar"`
	StartDateString *string                          `json:"start_date"`
	EndDateString   *string                          `json:"end_date"`
	ResponsibleUser organizationDTOs.UserMinimalDTO  `json:"responsible_user"`
	ProjectType     ProjectTypeDTO                   `json:"project_type"`
	ProjectStatus   ProjectStatusDTO                 `json:"project_status"`
	CreatedAt       time.Time                        `json:"created_at"`
	UpdatedAt       *time.Time                       `json:"updated_at"`
	CreatedByUser   organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser   *organizationDTOs.UserMinimalDTO `json:"updated_by"`
}
