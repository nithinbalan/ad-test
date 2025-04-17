package dtos

import (
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type StrategyProfileInitiativeListDTO struct {
	ID                              int     `json:"id"`
	Title                           string  `json:"title"`
	Description                     *string `json:"description"`
	StrategyProfileTitle            string  `json:"strategy_profile_title"`
	StrategyProfileFocusAreaTitle   *string `json:"strategy_profile_focus_area_title"`
	StrategyProfileGoalTitle        *string `json:"strategy_profile_goal_title"`
	StrategyprofileObjectiveTitle   *string `json:"strategy_profile_objective_title"`
	ResponsibleUserID               *uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      *string `json:"responsible_user_first_name_en"`
	ResponsibleUserLastNameEn       *string `json:"responsible_user_last_name_en"`
	ResponsibleUserFirstNameAr      *string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameAr       *string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          *string `json:"responsible_user_initial"`
	ResponsibleUserEmail            *string `json:"responsible_user_email"`
	ResponsibleUserColor            *string `json:"responsible_user_color"`
	ResponsibleUserProfileImagePath *string `json:"responsible_user_profile_image_path"`
	StartDate                       *string `json:"start_date" gorm:"column:start_date"`
	EndDate                         *string `json:"end_date" gorm:"column:end_date"`
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

type StrategyProfileInitiativeFormDTO struct {
	StrategyProfileFocusAreaID *uint   `json:"strategy_profile_focus_area_id"`
	StrategyProfileGoalID      *uint   `json:"strategy_profile_goal_id"`
	StrategyProfileObjectiveID *uint   `json:"strategy_profile_objective_id"`
	Title                      string  `json:"title" binding:"required"`
	Description                *string `json:"description"`
	ResponsibleUserID          *uint   `json:"responsible_user_id"`
	StartDate                  string  `json:"start_date" binding:"required"`
	EndDate                    string  `json:"end_date" binding:"required"`
}

type StrategyProfileInitiativeListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	StartDate   string
	EndDate     string
	Export      bool
}

type StrategyProfileInitiativeDetailDTO struct {
	ID                         uint                              `json:"id"`
	Title                      string                            `json:"title"`
	Description                *string                           `json:"description"`
	StrategyProfileID          uint                              `json:"strategy_profile_id"`
	StrategyProfile            StrategyProfileDetailDTO          `json:"strategy_profile"`
	StrategyProfileFocusAreaID *uint                             `json:"strategy_profile_focus_area_id"`
	StrategyProfileFocusArea   StrategyProfileFocusAreaDetailDTO `json:"strategy_profile_focus_area"`
	StrategyProfileGoalID      *uint                             `json:"strategy_profile_goal_id"`
	StrategyProfileGoal        StrategyProfileGoalDetailDTO      `json:"strategy_profile_goal"`
	StrategyProfileObjectiveID *uint                             `json:"strategy_profile_objective_id"`
	StrategyProfileObjective   StrategyProfileObjectiveDetailDTO `json:"strategy_profile_objective"`
	StartDateString            string                            `json:"start_date"`
	EndDateString              string                            `json:"end_date"`
	ResponsibleUserID          *uint                             `json:"responsible_user_id"`
	ResponsibleUser            *organizationDTOs.UserMinimalDTO  `json:"responsible_user"`
	CreatedAt                  time.Time                         `json:"created_at"`
	UpdatedAt                  *time.Time                        `json:"updated_at"`
	CreatedByUser              organizationDTOs.UserMinimalDTO   `json:"created_by"`
	UpdatedByUser              *organizationDTOs.UserMinimalDTO  `json:"updated_by"`
}

type StrategyProfileInitiativeParsedDates struct {
	StartDate *time.Time
	EndDate   *time.Time
}
