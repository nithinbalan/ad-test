package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"

	"time"
)

type AuditPlanToolListDTO struct {
	ID                          uint    `json:"id"`
	AuditPlanID                 uint    `json:"audit_plan_id"`
	ToolTypeID                  uint    `json:"tool_type_id"`
	ToolTypeTitleEn             string  `json:"tool_type_title_en"`
	ToolTypeTitleAr             string  `json:"tool_type_title_ar"`
	Title                       string  `json:"title"`
	Purpose                     *string `json:"purpose"`
	AuditPlanTitle              string  `json:"audit_plan_title"`
	CreatedUserFirstNameEn      string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string  `json:"created_user_last_name_ar"`
	CreatedUserInitial          string  `json:"created_user_initial"`
	CreatedUserColor            string  `json:"created_user_color"`
	CreatedUserEmail            string  `json:"created_user_email"`
	CreatedUserProfileImagePath string  `json:"created_user_profile_image_path"`
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

type AuditPlanToolFormDTO struct {
	ToolTypeID uint    `json:"tool_type_id" binding:"required"`
	Title      string  `json:"title" binding:"required,max=255"`
	Purpose    *string `json:"purpose"`
}

type AuditPlanToolDetailDTO struct {
	ID            uint                             `json:"id"`
	AuditPlanID   uint                             `json:"audit_plan_id"`
	ToolType      ToolTypeDTO                      `json:"tool_type"`
	AuditPlan     AuditPlanDetailDTO               `json:"audit_plan"`
	Title         string                           `json:"title"`
	Purpose       *string                          `json:"purpose"`
	CreatedAt     time.Time                        `json:"created_at"`
	UpdatedAt     *time.Time                       `json:"updated_at"`
	CreatedByUser organizationDTOs.UserMinimalDTO  `json:"created_by_id"`
	UpdatedByUser *organizationDTOs.UserMinimalDTO `json:"updated_by_id"`
}

type AuditPlanToolListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}

type AuditPlanToolListResponseDTO struct {
	Items      []*AuditPlanToolListDTO      `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type ToolTypeDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
}
