package dtos

import (
	organizationDTOs "jethur/grc/organization/dtos"
	projectDTOs "jethur/grc/project/dtos"
	"time"
)

type StrategyProfileProjectListDTO struct {
	ID                                       int     `json:"id"`
	ProjectID                                *uint   `json:"project_id"`
	ProjectTitleEn                           *string `json:"project_title_en"`
	ProjectTitleAr                           *string `json:"project_title_ar"`
	ProjectDescriptionEn                     *string `json:"project_description_en"`
	ProjectDescriptionAr                     *string `json:"project_description_ar"`
	ProjectProjectStatusID                   *uint   `json:"project_project_status_id"`
	ProjectProjectStatusTitleEn              *string `json:"project_project_status_title_en"`
	ProjectProjectStatusTitleAr              *string `json:"project_project_status_title_ar"`
	ProjectStatusColorCode                   *string `json:"project_status_color_code"`
	StrategyProfileTitle                     string  `json:"strategy_profile_title"`
	StrategyProfileStrategicObjectiveID      *uint   `json:"strategy_profile_strategic_objective_id"`
	StrategyProfileStrategicObjectiveTitleEn *string `json:"strategy_profile_strategic_objective_title_en"`
	StrategyProfileStrategicObjectiveTitleAr *string `json:"strategy_profile_strategic_objective_title_ar"`
	StrategyProfileStrategicIntiativeID      *uint   `json:"strategy_profile_strategic_intiative_id"`
	StrategyProfileStrategicIntiativeTitleEn *string `json:"strategy_profile_strategic_intiative_title_en"`
	StrategyProfileStrategicIntiativeTitleAr *string `json:"strategy_profile_strategic_intiative_title_ar"`
	CreatedUserID                            uint    `json:"created_user_id"`
	CreatedUserFirstNameEn                   string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr                   string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn                    string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr                    string  `json:"created_user_last_name_ar"`
	CreatedUserInitial                       string  `json:"created_user_initial"`
	CreatedUserColor                         string  `json:"created_user_color"`
	CreatedUserEmail                         string  `json:"created_user_email"`
	CreatedUserProfileImagePath              string  `json:"created_user_profile_image_path"`
	UpdatedUserID                            *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn                   *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr                   *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn                    *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr                    *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial                       *string `json:"updated_user_initial"`
	UpdatedUserColor                         *string `json:"updated_user_color"`
	UpdatedUserEmail                         *string `json:"updated_user_email"`
	UpdatedUserProfileImagePath              *string `json:"updated_user_profile_image_path"`
}

type StrategyProfileProjectFormDTO struct {
	ProjectIDs                          []uint `json:"project_ids" binding:"required,min=1"`
	StrategyProfileStrategicObjectiveID *uint  `json:"strategy_profile_strategic_objective_id"`
	StrategyProfileStrategicIntiativeID *uint  `json:"strategy_profile_strategic_intiative_id"`
}

type StrategyProfileProjectListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
	Filters  map[string][]string
}

type StrategyProfileProjectDetailDTO struct {
	ID                                  uint                                        `json:"id"`
	StrategyProfileStrategicObjectiveID *uint                                       `json:"strategy_profile_strategic_objective_id"`
	StrategyProfileStrategicObjective   *StrategyProfileStrategicObjectiveDetailDTO `json:"strategy_profile_strategic_objective"`
	StrategyProfileStrategicIntiativeID *uint                                       `json:"strategy_profile_strategic_intiative_id"`
	StrategyProfileStrategicIntiative   *StrategyProfileStrategicIntiativeDetailDTO `json:"strategy_profile_strategic_intiative"`
	ProjectID                           *uint                                       `json:"project_id"`
	Project                             *projectDTOs.ProjectDetailDTO               `json:"project"`
	CreatedByID                         uint                                        `json:"created_by_id"`
	CreatedByUser                       organizationDTOs.UserMinimalDTO             `json:"created_by_user"`
	UpdatedByID                         *uint                                       `json:"updated_by_id"`
	UpdatedByUser                       *organizationDTOs.UserMinimalDTO            `json:"updated_by_user"`
	CreatedAt                           time.Time                                   `json:"created_at"`
	UpdatedAt                           *time.Time                                  `json:"updated_at"`
}
