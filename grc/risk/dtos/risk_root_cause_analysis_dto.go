package dtos

type RiskRootCauseAnalysisListDTO struct {
	ID                          uint    `json:"id"`
	RiskRegisterTitle           string  `json:"risk_register_title" binding:"required"`
	Why                         uint    `json:"why" binding:"required"`
	Question                    string  `json:"question" binding:"required"`
	Answer                      string  `json:"answer,omitempty"`
	CreatedAt                   string  `json:"created_at"`
	UpdatedAt                   *string `json:"updated_at"`
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

type RiskRootCauseAnalysisListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
}

type RiskRootCauseAnalysisFormDTO struct {
	Question string `json:"question" binding:"required"`
	Answer   string `json:"answer" binding:"required"`
}
