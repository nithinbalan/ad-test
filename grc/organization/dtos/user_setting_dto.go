package dtos

import "time"

type UserSettingDTO struct {
	UserID         uint   `json:"user_id"`
	TimezoneID     uint   `json:"timezone_id"`
	Language       string `json:"language"`
	Timeout        int    `json:"timeout"`
	IsNotification bool   `json:"is_notification"`
	IsMfa          bool   `json:"is_mfa"`
}

type UserSettingFormDTO struct {
	TimezoneID     uint   `json:"timezone_id" binding:"required"`
	Language       string `json:"language" binding:"required"`
	Timeout        int    `json:"timeout" binding:"required"`
	IsNotification bool   `json:"is_notification" binding:"omitempty"`
}

// notification list
type NotificationListDTO struct {
	ID         uint   `json:"id"`
	ModulesID  uint   `json:"modules_id"`
	Identifier string `json:"identifier"`
	TitleEn    string `json:"title_en"`
	TitleAr    string `json:"title_ar"`
	IsEnable   bool   `json:"is_enable"`
}

// notification list module wise
type NotificationSettingListResponseDTO struct {
	ModuleID      uint                   `json:"module_id"`
	ModuleTitleEn string                 `json:"module_title_en"`
	ModuleTitleAr string                 `json:"module_title_ar"`
	Order         int                    `json:"order"`
	Notifications *[]NotificationListDTO `json:"notifications"`
}

type UserNotificationListDTO struct {
	ID                   uint      `json:"id"`
	UserID               uint      `json:"user_id"`
	ModuleID             uint      `json:"module_id"`
	ModuleIdentifier     string    `json:"module_identifier"`
	ModuleTitleEn        string    `json:"module_title_en"`
	ModuleTitleAr        string    `json:"module_title_ar"`
	NotificationType     string    `json:"notification_type"`
	ReferenceNumber      string    `json:"reference_number"`
	HeadingEn            string    `json:"heading_en"`
	HeadingAr            string    `json:"heading_ar"`
	TitleEn              string    `json:"title_en"`
	TitleAr              string    `json:"title_ar"`
	URLParam             string    `json:"url_param"`
	QueryParam           string    `json:"query_param"`
	IsNew                bool      `json:"is_new"`
	IsRead               bool      `json:"is_read"`
	IsSystemNotification bool      `json:"is_system_notification"`
	CreatedAt            time.Time `json:"created_at"`
	CreatedByFirstNameEn string    `json:"created_by_first_name_en"`
	CreatedByFirstNameAr string    `json:"created_by_first_name_ar"`
	CreatedByLastNameEn  string    `json:"created_by_last_name_en"`
	CreatedByLastNameAr  string    `json:"created_by_last_name_ar"`
}
