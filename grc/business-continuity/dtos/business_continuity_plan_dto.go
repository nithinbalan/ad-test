package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type BusinessContinuityPlanListDTO struct {
	ID                          uint    `json:"id"`
	ReferenceNumber             string  `json:"reference_number"`
	BusinessFunctionID          uint    `json:"business_function_id"`
	BusinessFunctionTitle       string  `json:"business_function_title"`
	Version                     *int    `json:"version"`
	Purpose                     *string `json:"purpose"`
	Scope                       *string `json:"scope"`
	Objectives                  *string `json:"objectives"`
	Assumptions                 *string `json:"assumptions"`
	CreatedAt                   string  `json:"created_at"`
	UpdatedAt                   string  `json:"updated_at"`
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
}

type BusinessContinuityPlanFormDTO struct {
	BusinessFunctionID uint    `json:"business_function_id"`
	Version            *int    `json:"version"`
	Purpose            *string `json:"purpose"`
	Scope              *string `json:"scope"`
	Objectives         *string `json:"objectives"`
	Assumptions        *string `json:"assumptions"`
}

type BusinessContinuityPlanListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	Export      bool
}

type BusinessContinuityPlanListResponseDTO struct {
	Items      []*BusinessContinuityPlanListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO     `json:"pagination"`
}

type BusinessContinuityPlanDetailDTO struct {
	ID                 uint                            `json:"id"`
	ReferenceNumber    string                          `json:"reference_number"`
	BusinessFunctionID uint                            `json:"business_function_id"`
	BusinessFunction   BusinessFunctionDetailDTO       `json:"business_function"`
	Version            int                             `json:"version"`
	Purpose            *string                         `json:"purpose"`
	Scope              *string                         `json:"scope"`
	Objectives         *string                         `json:"objectives"`
	Assumptions        *string                         `json:"assumptions"`
	CreatedAt          time.Time                       `json:"created_at"`
	UpdatedAt          *time.Time                      `json:"updated_at"`
	CreatedByUser      organizationDTOs.UserMinimalDTO `json:"created_by"`
	UpdatedByUser      organizationDTOs.UserMinimalDTO `json:"updated_by"`
}
