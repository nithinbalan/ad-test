package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"time"
)

type StakeholderActionPlanListDTO struct {
	ID                              uint       `json:"id"`
	StakeholderID                   uint       `json:"stakeholder_id"`
	SeverityLevelID                 *uint      `json:"severity_level_id"`
	SeverityLevelTitleEn            *string    `json:"severity_level_title_en"`
	SeverityLevelTitleAr            *string    `json:"severity_level_title_ar"`
	SeverityLevelColorCode          *string    `json:"severity_level_color_code"`
	Title                           string     `json:"title"`
	ResponsibleUserID               *uint      `json:"responsible_user_id"`
	ResponsibleUnitID               *uint      `json:"responsible_unit_id"`
	ResponsibleUnitTitleEn          *string    `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr          *string    `json:"responsible_unit_title_ar"`
	ResponsibleUserFirstNameEn      *string    `json:"responsible_user_first_name_en"`
	ResponsibleUserLastNameEn       *string    `json:"responsible_user_last_name_en"`
	ResponsibleUserFirstNameAr      *string    `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameAr       *string    `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          *string    `json:"responsible_user_initial"`
	ResponsibleUserEmail            *string    `json:"responsible_user_email"`
	ResponsibleUserColor            *string    `json:"responsible_user_color"`
	ResponsibleUserProfileImagePath *string    `json:"responsible_user_profile_image_path"`
	TargetDate                      *string    `json:"target_date"`
	Requirement                     *string    `json:"requirement"`
	Description                     *string    `json:"description"`
	ActionPlanStatusID              *uint      `json:"action_plan_status_id"`
	ActionPlanStatusTitleEn         *string    `json:"action_plan_status_title_en"`
	ActionPlanStatusTitleAr         *string    `json:"action_plan_status_title_ar"`
	ActionPlanStatusType            *string    `json:"action_plan_status_type"`
	ActionPlanStatusColorCode       *string    `json:"action_plan_status_color_code"`
	AssignedOwnerID                 *uint      `json:"assigned_owner_id"`
	AssignedOwnerFirstNameEn        *string    `json:"assigned_owner_first_name_en"`
	AssignedOwnerFirstNameAr        *string    `json:"assigned_owner_first_name_ar"`
	AssignedOwnerLastNameEn         *string    `json:"assigned_owner_last_name_en"`
	AssignedOwnerLastNameAr         *string    `json:"assigned_owner_last_name_ar"`
	AssignedOwnerInitial            *string    `json:"assigned_owner_initial"`
	AssignedOwnerColor              *string    `json:"assigned_owner_color"`
	AssignedOwnerEmail              *string    `json:"assigned_owner_email"`
	AssignedOwnerProfileImagePath   *string    `json:"assigned_owner_profile_image_path"`
	CreatedByID                     uint       `json:"created_by_id"`
	CreatedUserFirstNameEn          string     `json:"created_user_first_name_en"`
	CreatedUserLastNameEn           string     `json:"created_user_last_name_en"`
	CreatedUserProfileImagePath     string     `json:"created_user_profile_image_path"`
	CreatedUserEmail                string     `json:"created_user_email"`
	CreatedUserInitial              string     `json:"created_user_initial"`
	CreatedUserColor                string     `json:"created_user_color"`
	CreatedAt                       time.Time  `json:"created_at"`
	UpdatedByID                     *uint      `json:"updated_by_id"`
	UpdatedUserFirstNameEn          *string    `json:"updated_user_first_name_en"`
	UpdatedUserLastNameEn           *string    `json:"updated_user_last_name_en"`
	UpdatedUserProfileImagePath     *string    `json:"updated_user_profile_image_path"`
	UpdatedUserEmail                *string    `json:"updated_user_email"`
	UpdatedUserInitial              *string    `json:"updated_user_initial"`
	UpdatedUserColor                *string    `json:"updated_user_color"`
	UpdatedAt                       *time.Time `json:"updated_at"`
}

type StakeholderActionPlanDetailDTO struct {
	ID                 uint                                `json:"id"`
	StakeholderID      uint                                `json:"stakeholder_id"`
	Stakeholder        StakeholderDetailDTO                `json:"stakeholder"`
	SeverityLevelID    *uint                               `json:"severity_level_id"`
	SeverityLevel      *masterDTOs.SeverityLevelDTO        `json:"severity_level"`
	Title              string                              `json:"title"`
	ResponsibleUserID  *uint                               `json:"responsible_user_id"`
	ResponsibleUnitID  *uint                               `json:"responsible_unit_id"`
	ResponsibleUnit    *settingsDTOs.OrganizationUnitDTO   `json:"responsible_unit"`
	ResponsibleUser    *organizationDTOs.UserMinimalDTO    `json:"responsible_user"`
	AssignedOwnerID    *uint                               `json:"assigned_owner_id"`
	AssignedOwner      *organizationDTOs.UserMinimalDTO    `json:"assigned_owner"`
	TargetDateString   *string                             `json:"target_date"`
	ActionPlanStatusID *uint                               `json:"action_plan_status_id"`
	ActionPlanStatus   masterDTOs.ActionPlanStatusDTO      `json:"action_plan_status"`
	Requirement        *string                             `json:"requirement"`
	Description        *string                             `json:"description"`
	CreatedBYID        uint                                `json:"created_by_id"`
	CreatedByUser      organizationDTOs.UserMinimalDTO     `json:"created_by_user"`
	CreatedAt          time.Time                           `json:"created_at"`
	UpdatedBYID        *uint                               `json:"updated_by_id"`
	UpdatedByUser      *organizationDTOs.UserMinimalDTO    `json:"updated_by_user"`
	UpdatedAt          *time.Time                          `json:"updated_at"`
	Documents          []*StakeholderActionPlanDocumentDTO `json:"documents"`
}

type StakeholderActionPlanFormDTO struct {
	Title             string                        `json:"title" binding:"required"`
	Requirement       *string                       `json:"requirement"`
	Description       *string                       `json:"description"`
	SeverityLevelID   *uint                         `json:"severity_level_id"`
	ResponsibleUserID *uint                         `json:"responsible_user_id" binding:"required"`
	TargetDate        string                        `json:"target_date" binding:"required"`
	AssignedOwnerID   *uint                         `json:"assigned_owner_id"`
	Documents         []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
}

type StakeholderActionPlanQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	Export      bool
}

type StakeholderActionPlanParsedDates struct {
	TargetDate *time.Time
}

type StakeholderActionPlanDocumentDTO struct {
	ID           *uint   `json:"id,omitempty"`
	URL          *string `json:"url"`
	Name         string  `json:"file_name"`
	OriginalName string  `json:"original_name"`
	Token        string  `json:"token"`
	Type         string  `json:"type"`
	Size         int     `json:"size"`
}
