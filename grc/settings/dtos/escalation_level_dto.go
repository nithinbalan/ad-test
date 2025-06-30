package dtos

import "time"

type EscalationLevelDetailDTO struct {
	ID                 uint            `json:"id"`
	EscalationID       uint            `json:"escalation_id"`
	Escalation         EscalationDTO   `json:"escalation"`
	EscalationTo       string          `json:"escalation_to"`
	NotificationMethod string          `json:"notification_method"`
	Designation        *DesignationDTO `json:"designation"`
	DayOffset          int             `json:"day_offset"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          *time.Time      `json:"updated_at"`
	CreatedBYUser      UserMinimalDTO  `json:"created_by"`
	UpdatedBYUser      *UserMinimalDTO `json:"updated_by"`
}

type DesignationDTO struct {
	ID      int    `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Status  string `json:"status"`
}

type EscalationLevelListDTO struct {
	ID                          uint       `json:"id"`
	EscalationID                uint       `json:"escalation_id"`
	EscalationTitleEn           string     `json:"escalation_title_en"`
	EscalationTitleAr           string     `json:"escalation_title_ar"`
	NotificationMethod          string     `json:"notification_method"`
	EscalationTo                string     `json:"escalation_to"`
	DayOffset                   int        `json:"day_offset"`
	DesignaionTitleEn           *string    `json:"designation_title_en"`
	DesignaionTitleAr           *string    `json:"designation_title_ar"`
	CreatedUserID               uint       `json:"created_user_id"`
	CreatedUserFirstNameEn      string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string     `json:"created_user_last_name_ar"`
	CreatedUserInitial          string     `json:"created_user_initial"`
	CreatedUserColor            string     `json:"created_user_color"`
	CreatedUserEmail            string     `json:"created_user_email"`
	CreatedUserProfileImagePath string     `json:"created_user_profile_image_path"`
	UpdatedUserID               *uint      `json:"updated_user_id"`
	UpdatedUserFirstNameEn      *string    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      *string    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       *string    `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       *string    `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          *string    `json:"updated_user_initial"`
	UpdatedUserColor            *string    `json:"updated_user_color"`
	UpdatedUserEmail            *string    `json:"updated_user_email"`
	UpdatedUserProfileImagePath *string    `json:"updated_user_profile_image_path"`
	CreatedAt                   time.Time  `json:"created_at"`
	UpdatedAt                   *time.Time `json:"updated_at"`
}

type EscalationLevelListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
	Status  string
	Filters map[string][]string
}

type EscalationLevelFormDTO struct {
	EscalationID       uint   `json:"escalation_id" binding:"required"`
	NotificationMethod string `json:"notification_method" binding:"required,oneof=app email"`
	DayOffset          int    `json:"day_offset"`
	EscalationTo       string `json:"escalation_to" binding:"required"`
	DesignationID      *uint  `json:"designation_id"`
}
