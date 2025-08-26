package dtos

import "time"

type ReminderRuleDetailDTO struct {
	ID                 uint                `json:"id"`
	TitleEn            string              `json:"title_en"`
	TitleAr            string              `json:"title_ar"`
	ReminderID         uint                `json:"reminder_id"`
	Reminder           ReminderDTO         `json:"reminder"`
	TriggerConditionID uint                `json:"trigger_condition_id"`
	TriggerCondition   TriggerConditionDTO `json:"trigger_condition"`
	Type               string              `json:"type"`
	ReminderContext    ReminderContextDTO  `json:"reminder_context"`
	NotificationMethod string              `json:"notification_method"`
	DayOffset          int                 `json:"day_offset"`
	CreatedAt          time.Time           `json:"created_at"`
	UpdatedAt          *time.Time          `json:"updated_at"`
	CreatedBYUser      UserMinimalDTO      `json:"created_by"`
	UpdatedBYUser      *UserMinimalDTO     `json:"updated_by"`
}

type ReminderRuleListDTO struct {
	ID                          uint       `json:"id"`
	TitleEn                     string     `json:"title_en"`
	TitleAr                     string     `json:"title_ar"`
	ReminderID                  uint       `json:"reminder_id"`
	ReminderTitleEn             string     `json:"reminder_title_en"`
	ReminderTitleAr             string     `json:"reminder_title_ar"`
	ReminderType                string     `json:"reminder_type"`
	TriggerConditionID          uint       `json:"trigger_condition_id"`
	TriggerConditionTitleEn     string     `json:"trigger_condition_title_en"`
	TriggerConditionTitleAr     string     `json:"trigger_condition_title_ar"`
	TriggerConditionType        string     `json:"trigger_condition_type"`
	NotificationMethod          string     `json:"notification_method"`
	DayOffset                   int        `json:"day_offset"`
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

type ReminderRuleListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
	Status  string
	Primary bool
	Filters map[string][]string
}

type ReminderRuleFormDTO struct {
	ReminderID         uint   `json:"reminder_id" binding:"required"`
	TriggerConditionID uint   `json:"trigger_condition_id" binding:"required"`
	NotificationMethod string `json:"notification_method" binding:"required,oneof=app email"`
	TitleEn            string `json:"title_en" binding:"required"`
	TitleAr            string `json:"title_ar" binding:"required"`
	DayOffset          int    `json:"day_offset"`
	Type               string `json:"type" binding:"required"`
	ReminderContextID  uint   `json:"reminder_context_id" binding:"required"`
}
