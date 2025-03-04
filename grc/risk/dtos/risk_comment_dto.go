package dtos

import "time"

type RiskCommentListDTO struct {
	ID                      uint   `json:"id"`
	Comments                string `json:"comment"`
	RiskRegisterID          uint   `json:"risk_register_id"`
	CreatedAt               string `json:"created_at"`
	CreatedByID             uint   `json:"created_by_id"`
	CreatedUserFirstNameEn  string `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr  string `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn   string `json:"created_user_last_name_en"`
	CreatedUserLastNameAr   string `json:"created_user_last_name_ar"`
	CreatedUserInitial      string `json:"created_user_initial"`
	CreatedUserColor        string `json:"created_user_color"`
	CreatedUserEmail        string `json:"created_user_email"`
	CreatedUserProfileImage string `json:"created_user_profile_image"`
	RiskStatusTitleEn       string `json:"risk_status_title_en"`
	RiskStatusTitleAr       string `json:"risk_status_title_ar"`
	RiskStatusColorCode     string `json:"risk_status_color_code"`
	RiskTimelineTitleEn     string `json:"risk_timeline_title_en"`
	RiskTimelineTitleAr     string `json:"risk_timeline_title_ar"`
	RiskTimelineColorCode   string `json:"risk_timeline_color_code"`
	RiskTimelineType        string `json:"risk_timeline_type"`
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
