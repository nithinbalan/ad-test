package dtos

import "time"

type StakeholderEngagementEvaluationListDTO struct {
	ID                            int     `json:"id"`
	StakeholderEngagementID       uint    `json:"stakeholder_engagement_id"`
	StakeholderEngagementTitle    string  `json:"stakeholder_engagement_title"`
	Comment                       *string `json:"comment"`
	Score                         *int    `json:"score"`
	EvaluatedAt                   string  `json:"evaluated_at"`
	EvaluatedUserID               uint    `json:"evaluated_user_id"`
	EvaluatedUserFirstNameEn      string  `json:"evaluated_user_first_name_en"`
	EvaluatedUserFirstNameAr      string  `json:"evaluated_user_first_name_ar"`
	EvaluatedUserLastNameEn       string  `json:"evaluated_user_last_name_en"`
	EvaluatedUserLastNameAr       *string `json:"evaluated_user_last_name_ar"`
	EvaluatedUserInitial          string  `json:"evaluated_user_initial"`
	EvaluatedUserColor            string  `json:"evaluated_user_color"`
	EvaluatedUserEmail            string  `json:"evaluated_user_email"`
	EvaluatedUserProfileImagePath string  `json:"evaluated_user_profile_image_path"`
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

type StakeholderEngagementEvaluationFormDTO struct {
	Comment       *string `json:"comment"`
	Score         *int    `json:"score" binding:"required"`
	EvaluatedByID uint    `json:"evaluated_by_id"`
	EvaluatedAt   *string `json:"evaluated_at" binding:"omitempty,datetime=2006-01-02"`
	Description   *string `json:"description"`
}

type StakeholderEngagementEvaluationParsedDates struct {
	EvaluatedAt *time.Time
}
