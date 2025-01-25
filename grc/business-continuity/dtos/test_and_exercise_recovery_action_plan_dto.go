package dtos

import (
	"time"

	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
)

type TestAndExerciseRecoveryActionPlanListDTO struct {
	ID                              uint       `json:"id"`
	TestAndExerciseID               uint       `json:"test_and_exercise_id"`
	BcpPhaseID                      uint       `json:"bcp_phase_id"`
	BcpPhaseTitleEn                 string     `json:"bcp_phase_title_en"`
	BcpPhaseTitleAr                 string     `json:"bcp_phase_title_ar"`
	ProcessID                       uint       `json:"process_id"`
	ProcessTitleEn                  string     `json:"process_title_en"`
	ProcessTitleAr                  string     `json:"process_title_ar"`
	SeverityLevelID                 *uint      `json:"severity_level_id"`
	SeverityLevelTitleEn            *string    `json:"severity_level_title_en"`
	SeverityLevelTitleAr            *string    `json:"severity_level_title_ar"`
	SeverityLevelColorCode          *string    `json:"severity_level_color_code"`
	Title                           string     `json:"title"`
	ResponsibleUserID               uint       `json:"responsible_user_id"`
	ResponsibleUnitID               uint       `json:"responsible_unit_id"`
	ResponsibleUnitTitleEn          string     `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr          string     `json:"responsible_unit_title_ar"`
	ResponsibleUserFirstNameEn      string     `json:"responsible_user_first_name_en"`
	ResponsibleUserLastNameEn       string     `json:"responsible_user_last_name_en"`
	ResponsibleUserFirstNameAr      string     `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameAr       string     `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string     `json:"responsible_user_initial"`
	ResponsibleUserEmail            string     `json:"responsible_user_email"`
	ResponsibleUserColor            string     `json:"responsible_user_color"`
	ResponsibleUserProfileImagePath string     `json:"responsible_user_profile_image_path"`
	TargetDate                      *string    `json:"target_date"`
	ActionPlanStatusID              *uint      `json:"action_plan_status_id"`
	ActionPlanStatusTitleEn         *string    `json:"action_plan_status_title_en"`
	ActionPlanStatusTitleAr         *string    `json:"action_plan_status_title_ar"`
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

type TestAndExerciseRecoveryActionPlanDetailDTO struct {
	ID                 uint                              `json:"id"`
	TestAndExerciseID  uint                              `json:"test_and_exercise_id"`
	BcpPhaseID         uint                              `json:"bcp_phase_id"`
	BcpPhase           BcpPhaseDTO                       `json:"bcp_phase"`
	ProcessID          uint                              `json:"process_id"`
	Process            organizationDTOs.ProcessDetailDTO `json:"process"`
	SeverityLevelID    *uint                             `json:"severity_level_id"`
	SeverityLevel      *masterDTOs.SeverityLevelDTO      `json:"severity_level"`
	Title              string                            `json:"title"`
	ResponsibleUserID  uint                              `json:"responsible_user_id"`
	ResponsibleUnitID  uint                              `json:"responsible_unit_id"`
	ResponsibleUnit    settingsDTOs.OrganizationUnitDTO  `json:"responsible_unit"`
	ResponsibleUser    organizationDTOs.UserMinimalDTO   `json:"responsible_user"`
	TargetDateString   *string                           `json:"target_date"`
	ActionPlanStatusID *uint                             `json:"action_plan_status_id"`
	ActionPlanStatus   masterDTOs.ActionPlanStatusDTO    `json:"action_plan_status"`
	CreatedBYID        uint                              `json:"created_by_id"`
	CreatedByUser      organizationDTOs.UserMinimalDTO   `json:"created_by_user"`
	CreatedAt          time.Time                         `json:"created_at"`
	UpdatedBYID        *uint                             `json:"updated_by_id"`
	UpdatedByUser      *organizationDTOs.UserMinimalDTO  `json:"updated_by_user"`
	UpdatedAt          *time.Time                        `json:"updated_at"`
}

type TestAndExerciseRecoveryActionPlanFormDTO struct {
	Title             string `json:"title"`
	BcpPhaseID        uint   `json:"bcp_phase_id" binding:"required"`
	ProcessID         uint   `json:"process_id" binding:"required"`
	SeverityLevelID   *uint  `json:"severity_level_id"`
	ResponsibleUserID uint   `json:"responsible_user_id" binding:"required"`
	TargetDate        string `json:"target_date" binding:"required"`
}

type TestAndExerciseRecoveryActionPlanQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	Export      bool
}

type TestAndExerciseRecoveryActionPlanParsedDates struct {
	TargetDate *time.Time
}
