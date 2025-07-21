package dtos

import (
	"time"

	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
)

type ProjectActionPlanListDTO struct {
	ID                              uint    `json:"id"`
	Title                           string  `json:"title"`
	Description                     *string `json:"description"`
	ProjectID                       *uint   `json:"project_id"`
	ProjectTitleEn                  *string `json:"project_title_en"`
	ProjectTitleAr                  *string `json:"project_title_ar"`
	ProjectMilestoneID              *uint   `json:"project_milestone_id"`
	ProjectMilestoneTitleEn         *string `json:"project_milestone_title_en"`
	ProjectMilestoneTitleAr         *string `json:"project_milestone_title_ar"`
	ProjectMilestoneStartDate       *string `json:"project_milestone_start_date"`
	ProjectMilestoneEndDate         *string `json:"project_milestone_end_date"`
	SeverityLevelTitleEn            *string `json:"severity_level_title_en"`
	SeverityLevelTitleAr            *string `json:"severity_level_title_ar"`
	SeverityLevelColorCode          *string `json:"severity_level_color_code"`
	Comments                        *string `json:"comments"`
	ResponsibleUserID               uint    `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string  `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string  `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string  `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string  `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string  `json:"responsible_user_initial"`
	ResponsibleUserColor            string  `json:"responsible_user_color"`
	ResponsibleUserEmail            string  `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string  `json:"responsible_user_profile_image_path"`
	ReviewUserID                    *uint   `json:"review_user_id"`
	ReviewUserFirstNameEn           *string `json:"review_user_first_name_en"`
	ReviewUserFirstNameAr           *string `json:"review_user_first_name_ar"`
	ReviewUserLastNameEn            *string `json:"review_user_last_name_en"`
	ReviewUserLastNameAr            *string `json:"review_user_last_name_ar"`
	ReviewUserInitial               *string `json:"review_user_initial"`
	ReviewUserColor                 *string `json:"review_user_color"`
	ReviewUserEmail                 *string `json:"review_user_email"`
	ReviewUserProfileImagePath      *string `json:"review_user_profile_image_path"`
	TargetDate                      *string `json:"target_date"`
	Percentage                      int     `json:"percentage"`
	ActionPlanStatusTitleEn         string  `json:"action_plan_status_title_en"`
	ActionPlanStatusTitleAr         string  `json:"action_plan_status_title_ar"`
	ActionPlanStatusColorCode       string  `json:"action_plan_status_color_code"`
	ActionPlanStatusType            string  `json:"action_plan_status_type"`
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

type ProjectActionPlanListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}

type ProjectActionPlanListResponseDTO struct {
	Items      []*ProjectActionPlanListDTO  `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type ProjectActionPlanFormDTO struct {
	ProjectMilestoneID *uint   `json:"project_milestone_id"`
	Title              string  `json:"title" binding:"required"`
	Description        *string `json:"description" binding:"omitempty"`
	ResponsibleUserID  uint    `json:"responsible_user_id" binding:"required"`
	ReviewUserID       *uint   `json:"review_user_id" binding:"omitempty"`
	SeverityLevelID    *uint   `json:"severity_level_id" binding:"omitempty"`
	Comments           *string `json:"comments" binding:"omitempty"`
	TargetDate         *string `json:"target_date"`
	Percentage         int     `json:"percentage" binding:"omitempty"`
	ActionPlanStatusID uint    `json:"action_plan_status_id"`
}

type ProjectActionPlanParsedDates struct {
	TargetDate *time.Time
}

type ProjectActionPlanDetailDTO struct {
	ID                      int                                 `json:"id"`
	Project                 *ProjectDetailDTO                   `json:"project"`
	ProjectMilestone        *ProjectMilestoneDetailDTO          `json:"project_milestone"`
	Title                   string                              `json:"title"`
	Description             *string                             `json:"description"`
	ResponsibleUserID       uint                                `json:"responsible_user_id" binding:"required" gorm:"not null"`
	ResponsibleUser         organizationDTOs.UserMinimalDTO     `json:"responsible_user"`
	ResponsibleUserUnitID   *uint                               `json:"responsible_user_unit_id"`
	ResponsibleUserUnit     *settingsDTOs.OrganizationUnitDTO   `json:"responsible_user_unit"`
	ActionPlanStatus        masterDTOs.ActionPlanStatusDTO      `json:"action_plan_status"`
	Percentage              *int                                `json:"percentage"`
	SeverityLevel           *masterDTOs.SeverityLevelDTO        `json:"severity_level"`
	ReviewUser              *organizationDTOs.UserMinimalDTO    `json:"review_user"`
	ReviewUserUnitID        *uint                               `json:"review_user_unit_id"`
	ActionPlanStatusUpdates []*ProjectActionPlanStatusUpdateDTO `json:"action_plan_status_updates"`
	ReviewUserUnit          *settingsDTOs.OrganizationUnitDTO   `json:"review_user_unit"`
	TargetDateString        *string                             `json:"target_date"`
	Comments                *string                             `json:"comments"`
	CreatedAt               time.Time                           `json:"created_at"`
	UpdatedAt               *time.Time                          `json:"updated_at"`
	CreatedByUser           organizationDTOs.UserMinimalDTO     `json:"created_by_user"`
	UpdatedByUser           *organizationDTOs.UserMinimalDTO    `json:"updated_by_user"`
}

type ProjectActionPlanStatusUpdateDTO struct {
	ID               uint                                        `json:"id"`
	ActionPlanStatus masterDTOs.ActionPlanStatusDTO              `json:"action_plan_status"`
	Comment          string                                      `json:"comment"`
	Percentage       *int                                        `json:"percentage"`
	CreatedAt        time.Time                                   `json:"created_at"`
	UpdatedAt        *time.Time                                  `json:"updated_at"`
	CreatedByUser    organizationDTOs.UserMinimalDTO             `json:"created_by_user"`
	UpdatedByUser    *organizationDTOs.UserMinimalDTO            `json:"updated_by_user"`
	Documents        []*ProjectActionPlanStatusUpdateDocumentDTO `json:"documents"`
}

type ProjectActionPlanStatusUpdateDocumentDTO struct {
	ID           uint   `json:"id"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	Size         int    `json:"size"`
	Token        string `json:"token"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
}
