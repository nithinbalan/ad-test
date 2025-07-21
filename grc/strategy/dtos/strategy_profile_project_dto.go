package dtos

type StrategyProfileProjectListDTO struct {
	ID                                       int     `json:"id"`
	ProjectID                                *uint   `json:"project_id"`
	ProjectTitleEn                           *string `json:"project_title_en"`
	ProjectTitleAr                           *string `json:"project_title_ar"`
	ProjectDescriptionEn                     *string `json:"project_description_en"`
	ProjectDescriptionAr                     *string `json:"project_description_ar"`
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
