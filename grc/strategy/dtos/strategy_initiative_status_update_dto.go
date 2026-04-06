package dtos

type StrategyInitiativeStatusUpdateListDTO struct {
	ID                                 uint    `json:"id"`
	StrategyInitiativeID               uint    `json:"strategy_initiative_id"`
	StrategyInitiativeTitle            string  `json:"strategic_initiative_title"`
	StrategicInitiativeStatusTitleEn   string  `json:"strategic_initiative_status_title_en"`
	StrategicInitiativeStatusTitleAr   string  `json:"strategic_initiative_status_title_ar"`
	StrategicInitiativeStatusColorCode string  `json:"strategic_initiative_status_color_code"`
	StrategicInitiativeStatusType      string  `json:"strategic_initiative_status_type"`
	Comment                            string  `json:"comment"`
	Percentage                         string  `json:"percentage"`
	CreatedAt                          string  `json:"created_at"`
	UpdatedAt                          string  `json:"updated_at"`
	CreatedUserID                      uint    `json:"created_user_id"`
	CreatedUserFirstNameEn             string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr             string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn              string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr              *string `json:"created_user_last_name_ar"`
	CreatedUserInitial                 string  `json:"created_user_initial"`
	CreatedUserColor                   string  `json:"created_user_color"`
	CreatedUserEmail                   string  `json:"created_user_email"`
	CreatedUserProfileImagePath        string  `json:"created_user_profile_image_path"`
	UpdatedUserID                      *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn             string  `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr             string  `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn              string  `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr              string  `json:"updated_user_last_name_ar"`
	UpdatedUserInitial                 string  `json:"updated_user_initial"`
	UpdatedUserColor                   string  `json:"updated_user_color"`
	UpdatedUserEmail                   string  `json:"updated_user_email"`
	UpdatedUserProfileImagePath        string  `json:"updated_user_profile_image_path"`
}

type StrategyInitiativeStatusUpdateFormDTO struct {
	StrategyInitiativeStatusID uint   `json:"strategy_initiative_status_id" binding:"required"`
	Comment                    string `json:"comment"`
	Percentage                 uint   `json:"percentage" binding:"required"`
}
