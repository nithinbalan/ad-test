package dtos

import (
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"time"
)

type StrategyProfileActionPlanListDTO struct {
	ID                              int     `json:"id"`
	StrategyProfileTitle            string  `json:"strategy_profile_title"`
	StrategyProfileFocusAreaTitle   *string `json:"strategy_profile_focus_area_title"`
	StrategyProfileGoalTitle        *string `json:"strategy_profile_goal_title"`
	StrategyprofileObjectiveTitle   *string `json:"strategy_profile_objective_title"`
	StrategyprofileInitiativeTitle  *string `json:"strategy_profile_initiative_title"`
	StrategyprofileProjectTitle     *string `json:"strategy_profile_project_title"`
	Title                           string  `json:"title"`
	ResponsibleUserID               *uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      *string `json:"responsible_user_first_name_en"`
	ResponsibleUserLastNameEn       *string `json:"responsible_user_last_name_en"`
	ResponsibleUserFirstNameAr      *string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameAr       *string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          *string `json:"responsible_user_initial"`
	ResponsibleUserEmail            *string `json:"responsible_user_email"`
	ResponsibleUserColor            *string `json:"responsible_user_color"`
	ResponsibleUserProfileImagePath *string `json:"responsible_user_profile_image_path"`
	TargetDate                      *string `json:"target_date" gorm:"column:start_date"`
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

type StrategyProfileActionPlanFormDTO struct {
	StrategyProfileFocusAreaID  *uint   `json:"strategy_profile_focus_area_id"`
	StrategyProfileGoalID       *uint   `json:"strategy_profile_goal_id"`
	StrategyProfileObjectiveID  *uint   `json:"strategy_profile_objective_id"`
	StrategyProfileInitiativeID *uint   `json:"strategy_profile_initiative_id"`
	StrategyProfileProjectID    *uint   `json:"strategy_profile_project_id"`
	StrategyProfileMilestoneID  *uint   `json:"strategy_profile_milestone_id"`
	Title                       string  `json:"title" binding:"required"`
	Description                 *string `json:"description"`
	ResponsibleUserID           *uint   `json:"responsible_user_id"`
	SeverityLevelID             *uint   `json:"severity_level_id"`
	TargetDate                  string  `json:"target_date"`
	ReviewUserID                *uint   `json:"review_user_id"`
}

type StrategyProfileActionPlanListQueryParamsDTO struct {
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

type StrategyProfileActionPlanDetailDTO struct {
	ID                          uint                                         `json:"id"`
	StrategyProfileID           uint                                         `json:"strategy_profile_id"`
	StrategyProfile             StrategyProfileDetailDTO                     `json:"strategy_profile"`
	StrategyProfileFocusAreaID  *uint                                        `json:"strategy_profile_focus_area_id"`
	StrategyProfileFocusArea    StrategyProfileFocusAreaDetailDTO            `json:"strategy_profile_focus_area"`
	StrategyProfileGoalID       *uint                                        `json:"strategy_profile_goal_id"`
	StrategyProfileGoal         StrategyProfileGoalDetailDTO                 `json:"strategy_profile_goal"`
	StrategyProfileObjectiveID  *uint                                        `json:"strategy_profile_objective_id"`
	StrategyProfileObjective    StrategyProfileObjectiveDetailDTO            `json:"strategy_profile_objective"`
	StrategyProfileInitiativeID *uint                                        `json:"strategy_profile_initiative_id"`
	StrategyProfileInitiative   StrategyProfileInitiativeDetailDTO           `json:"strategy_profile_initiative"`
	StrategyProfileProjectID    *uint                                        `json:"strategy_profile_project_id"`
	StrategyProfileProject      StrategyProfileProjectDetailDTO              `json:"strategy_profile_project"`
	StrategyProfileMilestoneID  *uint                                        `json:"strategy_profile_milestone_id"`
	StrategyProfileMilestone    StrategyProfileMilestoneDetailDTO            `json:"strategy_profile_milestone"`
	SeverityLevelID             *uint                                        `json:"severity_level_id"`
	SeverityLevel               *masterDTOs.SeverityLevelDTO                 `json:"severity_level"`
	Title                       string                                       `json:"title"`
	Description                 *string                                      `json:"description"`
	ResponsibleUserID           *uint                                        `json:"responsible_user_id"`
	ResponsibleUnitID           *uint                                        `json:"responsible_unit_id"`
	ResponsibleUnit             settingsDTOs.OrganizationUnitDTO             `json:"responsible_unit"`
	ResponsibleUser             organizationDTOs.UserMinimalDTO              `json:"responsible_user"`
	ReviewUserID                *uint                                        `json:"review_user_id"`
	ReviewUser                  *organizationDTOs.UserMinimalDTO             `json:"review_user"`
	TargetDateString            string                                       `json:"target_date"`
	ActionPlanStatusID          *uint                                        `json:"action_plan_status_id"`
	ActionPlanStatus            masterDTOs.ActionPlanStatusDTO               `json:"action_plan_status"`
	CreatedAt                   time.Time                                    `json:"created_at"`
	UpdatedAt                   *time.Time                                   `json:"updated_at"`
	CreatedByUser               organizationDTOs.UserMinimalDTO              `json:"created_by"`
	UpdatedByUser               *organizationDTOs.UserMinimalDTO             `json:"updated_by"`
	ActionPlanStatusUpdates     []*StrategyProfileActionPlanStatusUpdateDTO  `json:"action_plan_status_updates"`
	ActionPlanComments          []*StrategyProfileActionPlanCommentDetailDTO `json:"action_plan_comments"`
}

type StrategyProfileActionPlanStatusUpdateDTO struct {
	ID               uint                                                `json:"id"`
	ActionPlanStatus masterDTOs.ActionPlanStatusDTO                      `json:"action_plan_status"`
	Comment          string                                              `json:"comment"`
	Percentage       *uint                                               `json:"percentage"`
	CreatedAt        time.Time                                           `json:"created_at"`
	UpdatedAt        *time.Time                                          `json:"updated_at"`
	CreatedByUser    organizationDTOs.UserMinimalDTO                     `json:"created_by_user"`
	UpdatedByUser    *organizationDTOs.UserMinimalDTO                    `json:"updated_by_user"`
	Documents        []*StrategyProfileActionPlanStatusUpdateDocumentDTO `json:"documents"`
}

type StrategyProfileActionPlanParsedDates struct {
	TargetDate *time.Time
}
