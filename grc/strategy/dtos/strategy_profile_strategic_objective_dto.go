package dtos

import (
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type StrategyProfileStrategicObjectiveListDTO struct {
	ID                              int      `json:"id"`
	StrategyProfileTitle            string   `json:"strategy_profile_title"`
	StrategyProfileFocusAreaID      *uint    `json:"strategy_profile_focus_area_id"`
	StrategyProfileFocusAreaTitleEn *string  `json:"strategy_profile_focus_area_title_en"`
	StrategyProfileFocusAreaTitleAr *string  `json:"strategy_profile_focus_area_title_ar"`
	StrategyProfileGoalID           *uint    `json:"strategy_profile_goal_id"`
	StrategyProfileGoalTitle        *string  `json:"strategy_profile_goal_title_en"`
	StrategicObjectiveID            *uint    `json:"strategic_objective_id"`
	StrategicObjectiveTitleEn       *string  `json:"strategic_objective_title_en"`
	StrategicObjectiveTitleAr       *string  `json:"strategic_objective_title_ar"`
	Score                           *float64 `json:"score"`
	CreatedUserID                   uint     `json:"created_user_id"`
	CreatedUserFirstNameEn          string   `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr          string   `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn           string   `json:"created_user_last_name_en"`
	CreatedUserLastNameAr           string   `json:"created_user_last_name_ar"`
	CreatedUserInitial              string   `json:"created_user_initial"`
	CreatedUserColor                string   `json:"created_user_color"`
	CreatedUserEmail                string   `json:"created_user_email"`
	CreatedUserProfileImagePath     string   `json:"created_user_profile_image_path"`
	UpdatedUserID                   *uint    `json:"updated_user_id"`
	UpdatedUserFirstNameEn          *string  `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr          *string  `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn           *string  `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr           *string  `json:"updated_user_last_name_ar"`
	UpdatedUserInitial              *string  `json:"updated_user_initial"`
	UpdatedUserColor                *string  `json:"updated_user_color"`
	UpdatedUserEmail                *string  `json:"updated_user_email"`
	UpdatedUserProfileImagePath     *string  `json:"updated_user_profile_image_path"`
}

type StrategyProfileStrategicObjectiveFormDTO struct {
	StrategyProfileFocusAreaID *uint    `json:"strategy_profile_focus_area_id"`
	StrategicObjectiveIDs      []uint   `json:"strategic_objective_ids" binding:"required,min=1"`
	StrategyProfileGoalIDs     *[]uint  `json:"strategy_profile_goal_ids"`
	Score                      *float64 `json:"score"`
}

type StrategyProfileStrategicObjectiveListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
	Filters  map[string][]string
}

type StrategyProfileStrategicObjectiveDetailDTO struct {
	ID                         uint                               `json:"id"`
	StrategyProfileFocusAreaID *uint                              `json:"strategy_profile_focus_area_id"`
	StrategyProfileFocusArea   *StrategyProfileFocusAreaDetailDTO `json:"strategy_profile_focus_area"`
	StrategicObjectiveID       *uint                              `json:"strategic_objective_id"`
	StrategicObjective         *StrategicObjectiveDetailDTO       `json:"strategic_objective"`
	Score                      *float64                           `json:"score"`
	CreatedByID                uint                               `json:"created_by_id"`
	CreatedByUser              organizationDTOs.UserMinimalDTO    `json:"created_by_user"`
	UpdatedByID                *uint                              `json:"updated_by_id"`
	UpdatedByUser              *organizationDTOs.UserMinimalDTO   `json:"updated_by_user"`
	CreatedAt                  time.Time                          `json:"created_at"`
	UpdatedAt                  *time.Time                         `json:"updated_at"`
}
