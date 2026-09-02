package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type ProjectMilestoneListDTO struct {
	ID                              uint    `json:"id"`
	ProjectID                       uint    `json:"project_id"`
	TitleEn                         string  `json:"title_en"`
	TitleAr                         string  `json:"title_ar"`
	DescriptionEn                   *string `json:"description_en"`
	DescriptionAr                   *string `json:"description_ar"`
	PlannedDate                     *string `json:"planned_date"`
	ActualDate                      *string `json:"actual_date"`
	ProjectMilestoneStatusID        uint    `json:"project_milestone_status_id"`
	ProjectMilestoneStatusTitleEn   string  `json:"project_milestone_status_title_en"`
	ProjectMilestoneStatusTitleAr   string  `json:"project_milestone_status_title_ar"`
	ProjectMilestoneStatusType      string  `json:"project_milestone_status_type"`
	ProjectMilestoneStatusColorCode string  `json:"project_milestone_status_color_code"`
	ProjectPhaseID                  *uint   `json:"project_phase_id"`
	ProjectPhaseTitleEn             *string `json:"project_phase_title_en"`
	ProjectPhaseTitleAr             *string `json:"project_phase_title_ar"`
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

type ProjectMilestoneListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	DateFrom    string
	DateTo      string
}

type ProjectMilestoneListResponseDTO struct {
	Items      []*ProjectMilestoneListDTO   `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type ProjectMilestoneDetailDTO struct {
	ID                     uint                                 `json:"id"`
	Project                ProjectDetailDTO                     `json:"project"`
	TitleEn                string                               `json:"title_en"`
	TitleAr                string                               `json:"title_ar"`
	DescriptionEn          *string                              `json:"description_en"`
	DescriptionAr          *string                              `json:"description_ar"`
	PlannedDateString      *string                              `json:"planned_date"`
	ActualDateString       *string                              `json:"actual_date"`
	ProjectMilestoneStatus masterDTOs.ProjectMilestoneStatusDTO `json:"project_milestone_status"`
	ProjectPhaseID         *uint                                `json:"project_phase_id"`
	ProjectPhase           *ProjectPhaseDTO                     `json:"project_phase"`
	CreatedByUser          organizationDTOs.UserMinimalDTO      `json:"created_by"`
	UpdatedByUser          *organizationDTOs.UserMinimalDTO     `json:"updated_by"`
	CreatedAt              time.Time                            `json:"created_at"`
	UpdatedAt              *time.Time                           `json:"updated_at"`
}

type ProjectMilestoneFormDTO struct {
	TitleEn                  string  `json:"title_en" binding:"required"`
	TitleAr                  string  `json:"title_ar" binding:"required"`
	DescriptionEn            *string `json:"description_en"`
	DescriptionAr            *string `json:"description_ar"`
	PlannedDate              *string `json:"planned_date"`
	ActualDate               *string `json:"actual_date"`
	ProjectMilestoneStatusID uint    `json:"project_milestone_status_id" binding:"required"`
	ProjectPhaseID           *uint   `json:"project_phase_id"`
}

type ProjectMilestoneStatusUpdateDTO struct {
	ProjectMilestoneStatusID uint    `json:"project_milestone_status_id" binding:"required"`
	ActualDate               *string `json:"actual_date"`
}
