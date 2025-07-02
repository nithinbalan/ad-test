package dtos

import "jethur/grc/common/dtos"

type StrategicObjectiveKPIListDTO struct {
	ID                              int     `json:"id"`
	Title                           string  `json:"title"`
	Description                     *string `json:"description"`
	KPICategoryID                   uint    `json:"kpi_category_id"`
	KPICategoryTitleEn              string  `json:"kpi_category_title_en"`
	KPICategoryTitleAr              string  `json:"kpi_category_title_ar"`
	ResponsibleUserID               *uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      *string `json:"responsible_user_first_name_en"`
	ResponsibleUserLastNameEn       *string `json:"responsible_user_last_name_en"`
	ResponsibleUserFirstNameAr      *string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameAr       *string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          *string `json:"responsible_user_initial"`
	ResponsibleUserEmail            *string `json:"responsible_user_email"`
	ResponsibleUserColor            *string `json:"responsible_user_color"`
	ResponsibleUserProfileImagePath *string `json:"responsible_user_profile_image_path"`
	KPIStatusID                     uint    `json:"kpi_status_id"`
	KPIStatusTitleEn                string  `json:"kpi_status_title_en"`
	KPIStatusTitleAr                string  `json:"kpi_status_title_ar"`
	KPIStatusColorCode              string  `json:"kpi_status_color_code"`
	KPIStatusType                   string  `json:"kpi_status_type"`
	CreatedUserID                   uint    `json:"created_user_id"`
	CreatedUserFirstNameEn          string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr          string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn           string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr           string  `json:"created_user_last_name_ar"`
	CreatedUserInitial              string  `json:"created_user_initial"`
	CreatedUserColor                string  `json:"created_user_color"`
	CreatedUserEmail                string  `json:"created_user_email"`
	CreatedUserProfileImagePath     string  `json:"created_user_profile_image_path"`
	UpdatedUserID                   *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn          *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr          *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn           *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr           *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial              *string `json:"updated_user_initial"`
	UpdatedUserColor                *string `json:"updated_user_color"`
	UpdatedUserEmail                *string `json:"updated_user_email"`
	UpdatedUserProfileImagePath     *string `json:"updated_user_profile_image_path"`
}

type StrategicObjectiveKPIListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	Export      bool
}
type StrategicObjectiveKPIListResponseDTO struct {
	Items      []*StrategicObjectiveKPIListDTO `json:"items"`
	Pagination dtos.PaginationInfoDTO          `json:"pagination"`
}
type StrategicObjectiveKPIFormDTO struct {
	KPIRegisterIDs []uint `json:"kpi_register_ids" binding:"required,min=1,dive,required"`
}
