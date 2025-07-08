package dtos

type StrategyProfileFocusAreaListDTO struct {
	ID                          int     `json:"id"`
	FocusAreaID                 *uint   `json:"focus_area_id"`
	FocusAreaTitleEn            *string `json:"focus_area_title_en"`
	FocusAreaTitleAr            *string `json:"focus_area_title_ar"`
	FocusAreaDescriptionEn      *string `json:"focus_area_description_en"`
	FocusAreaDescriptionAr      *string `json:"focus_area_description_ar"`
	StrategyProfileTitle        string  `json:"strategy_profile_title"`
	CreatedUserID               uint    `json:"created_user_id"`
	CreatedUserFirstNameEn      string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string  `json:"created_user_last_name_ar"`
	CreatedUserInitial          string  `json:"created_user_initial"`
	CreatedUserColor            string  `json:"created_user_color"`
	CreatedUserEmail            string  `json:"created_user_email"`
	CreatedUserProfileImagePath string  `json:"created_user_profile_image_path"`
	UpdatedUserID               *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn      *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          *string `json:"updated_user_initial"`
	UpdatedUserColor            *string `json:"updated_user_color"`
	UpdatedUserEmail            *string `json:"updated_user_email"`
	UpdatedUserProfileImagePath *string `json:"updated_user_profile_image_path"`
}

type StrategyProfileFocusAreaFormDTO struct {
	FocusAreaIDs []uint `json:"focus_area_ids" binding:"required,min=1"`
}

type StrategyProfileFocusAreaListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	Export      bool
}
