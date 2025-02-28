package dtos

import "time"

type RiskCommentListDTO struct {
	ID                      uint      `json:"id"`
	Comments                string    `json:"comment"`
	RiskRegisterID          uint      `json:"risk_register_id"`
	CreatedAt               time.Time `json:"created_at"`
	CreatedByID             uint      `json:"created_by_id"`
	CreatedUserFirstNameEn  string    `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr  string    `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn   string    `json:"created_user_last_name_en"`
	CreatedUserLastNameAr   string    `json:"created_user_last_name_ar"`
	CreatedUserInitial      string    `json:"created_user_initial"`
	CreatedUserColor        string    `json:"created_user_color"`
	CreatedUserEmail        string    `json:"created_user_email"`
	CreatedUserProfileImage string    `json:"created_user_profile_image"`
}

type RiskCommentDetailDTO struct {
	ID                      uint      `json:"id"`
	Comment                 string    `json:"comment"`
	RiskRegisterID          uint      `json:"risk_register_id"`
	CreatedAt               time.Time `json:"created_at"`
	CreatedByID             uint      `json:"created_by_id"`
	CreatedUserFirstNameEn  string    `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr  string    `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn   string    `json:"created_user_last_name_en"`
	CreatedUserLastNameAr   string    `json:"created_user_last_name_ar"`
	CreatedUserInitial      string    `json:"created_user_initial"`
	CreatedUserColor        string    `json:"created_user_color"`
	CreatedUserEmail        string    `json:"created_user_email"`
	CreatedUserProfileImage string    `json:"created_user_profile_image"`
}

type RiskCommentFormDTO struct {
	RiskTimelineTypeID *uint  `json:"risk_timeline_type_id"`
	Comment            string `json:"comment" binding:"required"`
}
