package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type TaskListDTO struct {
	ID                              uint    `json:"id"`
	ReferenceNumber                 string  `json:"reference_number"`
	ModuleID                        uint    `json:"module_id"`
	ComplianceActionPlanID          uint    `json:"compliance_action_plan_id"`
	SLAActionPlanID                 uint    `json:"sla_action_plan_id"`
	IncidentActionPlanID            uint    `json:"incident_action_plan_id"`
	ControlActionPlanID             uint    `json:"control_action_plan_id"`
	SeverityLevelID                 uint    `json:"severity_level_id"`
	TaskStatusID                    uint    `json:"task_status_id"`
	ModuleTitleEn                   string  `json:"module_title_en"`
	ModuleTitleAr                   string  `json:"module_title_ar"`
	Title                           string  `json:"title"`
	Description                     string  `json:"description"`
	SeverityLevelTitleEn            string  `json:"severity_level_title_en"`
	SeverityLevelTitleAr            string  `json:"severity_level_title_ar"`
	SeverityLevelColorCode          string  `json:"severity_level_color_code"`
	ResponsibleUserID               uint    `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string  `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string  `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string  `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string  `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string  `json:"responsible_user_initial"`
	ResponsibleUserColor            string  `json:"responsible_user_color"`
	ResponsibleUserEmail            string  `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string  `json:"responsible_user_profile_image_path"`
	TargetDate                      string  `json:"target_date"`
	URLParam                        string  `json:"url_param"`
	TaskStatusTitleEn               string  `json:"task_status_title_en"`
	TaskStatusTitleAr               string  `json:"task_status_title_ar"`
	TaskStatusColorCode             string  `json:"task_status_color_code"`
	CreatedAt                       string  `json:"created_at"`
	UpdatedAt                       *string `json:"updated_at"`
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
	UpdatedUserFirstNameEn          string  `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr          string  `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn           string  `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr           string  `json:"updated_user_last_name_ar"`
	UpdatedUserInitial              string  `json:"updated_user_initial"`
	UpdatedUserColor                string  `json:"updated_user_color"`
	UpdatedUserEmail                string  `json:"updated_user_email"`
	UpdatedUserProfileImagePath     string  `json:"updated_user_profile_image_path"`
}

type TaskListQueryParamsDTO struct {
	Sort           string
	SortBy         string
	Keyword        string
	Status         string
	Page           int
	PageSize       int
	Filters        map[string][]string
	QuickFilter    []string
	TargetDateFrom string
	TargetDateTo   string
}

type TaskListResponseDTO struct {
	Items      []*TaskListDTO               `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type TaskParsedDates struct {
	TargetDate *time.Time
}
