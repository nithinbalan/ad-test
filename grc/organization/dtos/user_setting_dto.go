package dtos

type UserSettingDTO struct {
	UserID         uint   `json:"user_id"`
	TimezoneID     uint   `json:"timezone_id"`
	Language       string `json:"language"`
	Timeout        int    `json:"timeout"`
	IsNotification bool   `json:"is_notification"`
}

type UserSettingFormDTO struct {
	TimezoneID uint   `json:"timezone_id" binding:"required"`
	Language   string `json:"language" binding:"required"`
	Timeout    int    `json:"timeout" binding:"required"`
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
	Notifications *[]NotificationListDTO `json:"notifications"`
}
