package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"

	"time"
)

type ComplianceActionPlanListDTO struct {
	ID                              uint   `json:"id"`
	ComplianceClauseNo              string `json:"compliance_clause_no"`
	ComplianceClauseTitle           string `json:"compliance_clause_title"`
	ComplianceRequirementTitle      string `json:"compliance_requirement_title"`
	Title                           string `json:"title"`
	Description                     string `json:"description"`
	SeverityLevelTitleEn            string `json:"severity_level_title_en"`
	SeverityLevelTitleAr            string `json:"severity_level_title_ar"`
	SeverityLevelColorCode          string `json:"severity_level_color_code"`
	ResponsibleUserFirstNameEn      string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string `json:"responsible_user_initial"`
	ResponsibleUserColor            string `json:"responsible_user_color"`
	ResponsibleUserEmail            string `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string `json:"responsible_user_profile_image_path"`
	TargetDate                      string `json:"target_date"`
	Percentage                      uint   `json:"percentage"`
	Comments                        string `json:"comments"`
	ActionPlanStatusTitleEn         string `json:"action_plan_status_title_en"`
	ActionPlanStatusTitleAr         string `json:"action_plan_status_title_ar"`
	ActionPlanStatusColorCode       string `json:"action_plan_status_color_code"`
	CreatedUserFirstNameEn          string `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr          string `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn           string `json:"created_user_last_name_en"`
	CreatedUserLastNameAr           string `json:"created_user_last_name_ar"`
	CreatedUserInitial              string `json:"created_user_initial"`
	CreatedUserColor                string `json:"created_user_color"`
	CreatedUserEmail                string `json:"created_user_email"`
	CreatedUserProfileImagePath     string `json:"created_user_profile_image_path"`
	UpdatedUserFirstNameEn          string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr          string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn           string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr           string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial              string `json:"updated_user_initial"`
	UpdatedUserColor                string `json:"updated_user_color"`
	UpdatedUserEmail                string `json:"updated_user_email"`
	UpdatedUserProfileImagePath     string `json:"updated_user_profile_image_path"`
}

type ComplianceActionPlanListQueryParamsDTO struct {
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

type ComplianceActionPlanListResponseDTO struct {
	Items      []*ComplianceActionPlanListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO   `json:"pagination"`
}

type ComplianceActionPlanFormDTO struct {
	ComplianceRequirementID *uint   `json:"compliance_requirement_id"`
	Title                   string  `json:"title" binding:"required"`
	Description             *string `json:"description" binding:"omitempty"`
	SeverityLevelID         *uint   `json:"severity_level_id" binding:"omitempty"`
	ResponsibleUserID       uint    `json:"responsible_user_id" binding:"required"`
	TargetDate              *string `json:"target_date" binding:"omitempty,datetime=2006-01-02"`
	Percentage              *uint   `json:"percentage" binding:"omitempty"`
	Comments                *string `json:"comments" binding:"omitempty"`
	ActionPlanStatusID      uint    `json:"action_plan_status_id" binding:"required"`
}

type ComplianceActionPlanParsedDates struct {
	TargetDate *time.Time
}

type ComplianceRequirementMinimalDTO struct {
	ID               int                 `json:"id"`
	ComplianceClause ComplianceClauseDTO `json:"compliance_clause"`
	Title            string              `json:"title"`
}

type ComplianceActionPlanDetailDTO struct {
	ID                     int                                  `json:"id"`
	ComplianceRequirement  ComplianceRequirementMinimalDTO      `json:"compliance_requirement"`
	Title                  string                               `json:"title"`
	Description            string                               `json:"description"`
	SeverityLevel          *masterDTOs.SeverityLevelDTO         `json:"severity_level"`
	ResponsibleUser        organizationDTOs.UserMinimalDTO      `json:"responsible_user"`
	ActionPlanStatusUpdate *ComplianceActionPlanStatusUpdateDTO `json:"action_plan_status"`
	Percentage             *uint                                `json:"percentage"`
	Comments               *string                              `json:"comments"`
	TargetDate             *time.Time                           `json:"target_date"`
	CreatedAt              time.Time                            `json:"created_at"`
	UpdatedAt              *time.Time                           `json:"updated_at"`
	CreatedByUser          organizationDTOs.UserMinimalDTO      `json:"created_by_user"`
	UpdatedByUser          *organizationDTOs.UserMinimalDTO     `json:"updated_by_user"`
}

type ComplianceActionPlanStatusUpdateDTO struct {
	ID               uint                                           `json:"id"`
	ActionPlanStatus masterDTOs.ActionPlanStatusDTO                 `json:"action_plan_status"`
	Comment          string                                         `json:"comment"`
	Percentage       *uint                                          `json:"percentage"`
	CreatedAt        time.Time                                      `json:"created_at"`
	UpdatedAt        *time.Time                                     `json:"updated_at"`
	CreatedByUser    organizationDTOs.UserMinimalDTO                `json:"created_by_user"`
	UpdatedByUser    *organizationDTOs.UserMinimalDTO               `json:"updated_by_user"`
	Documents        []*ComplianceActionPlanStatusUpdateDocumentDTO `json:"documents"`
}

type ComplianceActionPlanStatusUpdateDocumentDTO struct {
	ID           uint   `json:"id"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	Size         int    `json:"size"`
	Token        string `json:"token"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
}
