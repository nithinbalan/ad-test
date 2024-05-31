package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type TaskListDTO struct {
	ID                         uint   `json:"id"`
	ReferenceNumber            string `json:"reference_number"`
	ModuleID                   uint   `json:"module_id"`
	ComplianceActionPlanID     uint   `json:"compliance_action_plan_id"`
	SLAActionPlanID            uint   `json:"sla_action_plan_id"`
	IncidentActionPlanID       uint   `json:"incident_action_plan_id"`
	ControlActionPlanID        uint   `json:"control_action_plan_id"`
	SeverityLevelID            uint   `json:"severity_level_id"`
	ResponsibleUserID          uint   `json:"responsible_user_id"`
	TaskStatusID               uint   `json:"task_status_id"`
	ModuleTitleEn              string `json:"module_title_en"`
	ModuleTitleAr              string `json:"module_title_ar"`
	Title                      string `json:"title"`
	Description                string `json:"description"`
	SeverityLevelTitleEn       string `json:"severity_level_title_en"`
	SeverityLevelTitleAr       string `json:"severity_level_title_ar"`
	SeverityLevelColorCode     string `json:"severity_level_color_code"`
	ResponsibleUserFirstNameEn string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn  string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr  string `json:"responsible_user_last_name_ar"`
	TargetDate                 string `json:"target_date"`
	URL                        string `json:"url"`
	TaskStatusTitleEn          string `json:"task_status_title_en"`
	TaskStatusTitleAr          string `json:"task_status_title_ar"`
	TaskStatusColorCode        string `json:"task_status_color_code"`
	CreatedUserFirstNameEn     string `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr     string `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn      string `json:"created_user_last_name_en"`
	CreatedUserLastNameAr      string `json:"created_user_last_name_ar"`
	UpdatedUserFirstNameEn     string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr     string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn      string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr      string `json:"updated_user_last_name_ar"`
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

type UserTaskStatusCountDTO struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

type TaskParsedDates struct {
	TargetDate *time.Time
}
