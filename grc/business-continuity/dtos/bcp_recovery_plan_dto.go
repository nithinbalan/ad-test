package dtos

import (
	"time"

	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
)

type BcpRecoveryPlanListDTO struct {
	ID                              uint    `json:"id"`
	Title                           string  `json:"title"`
	BcpPhaseID                      uint    `json:"bcp_phase_id"`
	BcpPhaseTitleEn                 string  `json:"bcp_phase_title_en"`
	BcpPhaseTitleAr                 string  `json:"bcp_phase_title_ar"`
	BcpProcessTitleEn               string  `json:"bcp_process_title_en"`
	BcpProcessTitleAr               string  `json:"bcp_process_title_ar"`
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
	TargetCompletion                string  `json:"target_completion"`
	TargetCompletionUnitTitleEn     *string `json:"target_completion_unit_title_en"`
	TargetCompletionUnitTitleAr     *string `json:"target_completion_unit_title_ar"`
	CreatedAt                       string  `json:"created_at"`
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

type BcpRecoveryPlanListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}

type BcpRecoveryPlanFormDTO struct {
	BcpPhaseID             *uint   `json:"bcp_phase_id"`
	BcpProcessID           *uint   `json:"bcp_process_id"`
	Title                  string  `json:"title" binding:"required"`
	SeverityLevelID        *uint   `json:"severity_level_id" binding:"omitempty"`
	ResponsibleUserID      uint    `json:"responsible_user_id" binding:"required"`
	TargetCompletion       *string `json:"target_completion"`
	TargetCompletionUnitID *uint   `json:"target_completion_unit_id"`
}

type BcpRecoveryPlanDetailDTO struct {
	ID                     int                              `json:"id"`
	Title                  string                           `json:"title"`
	BusinessContinuityPlan BusinessContinuityPlanDetailDTO  `json:"business_continuity_plan"`
	BcpPhaseID             uint                             `json:"bcp_phase_id"`
	BcpPhase               BcpPhaseDTO                      `json:"bcp_phase"`
	BcpProcessID           uint                             `json:"bcp_process_id"`
	BcpProcess             BcpProcessDetailsDTO             `json:"bcp_process"`
	SeverityLevel          *masterDTOs.SeverityLevelDTO     `json:"severity_level"`
	ResponsibleUserID      uint                             `json:"responsible_user_id" binding:"required" gorm:"not null"`
	ResponsibleUser        organizationDTOs.UserMinimalDTO  `json:"responsible_user"`
	ResponsibleUnitID      uint                             `json:"responsible_unit_id"`
	ResponsibleUnit        settingsDTOs.OrganizationUnitDTO `json:"responsible_unit"`
	TargetCompletion       *string                          `json:"target_completion"`
	TargetCompletionUnitID *uint                            `json:"target_completion_time_unit_id"`
	TargetCompletionUnit   *masterDTOs.TimeUnitDTO          `json:"target_completion_time_unit"`
	CreatedAt              time.Time                        `json:"created_at"`
	UpdatedAt              *time.Time                       `json:"updated_at"`
	CreatedByUser          organizationDTOs.UserMinimalDTO  `json:"created_by_user"`
	UpdatedByUser          *organizationDTOs.UserMinimalDTO `json:"updated_by_user"`
}
