package dtos

type QualysSettingsResponseDTO struct {
	ID           uint   `json:"id"`
	PlatformURL  string `json:"platform_url"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	LookbackDays int    `json:"lookback_days"`
	Status       string `json:"status"`
}

type UpdateQualysSettingsDTO struct {
	PlatformURL  string `json:"platform_url" binding:"required"`
	Username     string `json:"username" binding:"required"`
	Password     string `json:"password" binding:"required"`
	LookbackDays int    `json:"lookback_days" binding:"required"`
}

type TestQualysConnectionDTO struct {
	PlatformURL string `json:"platform_url"`
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type QualysConnectionError struct {
	ErrorMsg string `json:"error_msg"`
}