package dtos

type StrategyProfileProjectMilestoneListDTO struct {
	ID                            int     `json:"id"`
	StrategyProfileID             uint    `json:"strategy_profile_id"`
	StrategyProfileTitle          string  `json:"strategy_profile_title"`
	StrategyProfileProjectID      *uint   `json:"strategy_profile_project_id"`
	StrategyProfileProjectTitleEn *string `json:"strategy_profile_project_title_en"`
	StrategyProfileProjectTitleAr *string `json:"strategy_profile_project_title_ar"`
	ProjectMilestoneID            *uint   `json:"project_milestone_id"`
	ProjectMilestoneTitleEn       *string `json:"project_milestone_title_en"`
	ProjectMilestoneTitleAr       *string `json:"project_milestone_title_ar"`
	CreatedUserID                 uint    `json:"created_user_id"`
	CreatedUserFirstNameEn        string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr        string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn         string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr         string  `json:"created_user_last_name_ar"`
	CreatedUserInitial            string  `json:"created_user_initial"`
	CreatedUserColor              string  `json:"created_user_color"`
	CreatedUserEmail              string  `json:"created_user_email"`
	CreatedUserProfileImagePath   string  `json:"created_user_profile_image_path"`
	UpdatedUserID                 *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn        *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr        *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn         *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr         *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial            *string `json:"updated_user_initial"`
	UpdatedUserColor              *string `json:"updated_user_color"`
	UpdatedUserEmail              *string `json:"updated_user_email"`
	UpdatedUserProfileImagePath   *string `json:"updated_user_profile_image_path"`
}

type StrategyProfileProjectMilestoneFormDTO struct {
	StrategyProfileProjectID *uint  `json:"strategy_profile_project_id"`
	ProjectMilestoneIDs      []uint `json:"project_milestone_ids" binding:"required,min=1"`
}

type StrategyProfileProjectMilestoneListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
	Filters  map[string][]string
}
