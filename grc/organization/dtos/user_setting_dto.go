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
