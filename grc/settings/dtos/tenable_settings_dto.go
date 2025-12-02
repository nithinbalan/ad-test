package dtos

type TenableSettingsResponseDTO struct {
	ID                 uint   `json:"id"`
	BaseURL            string `json:"base_url"`
	AccessKey          string `json:"access_key"`
	SecretKey          string `json:"secret_key"`
	SyncFrequencyHours int    `json:"sync_frequency_hours"`
	LookbackDays       int    `json:"lookback_days"`
	Status             string `json:"status"`
}

type UpdateTenableSettingsDTO struct {
	BaseURL            string `json:"base_url" binding:"required"`
	AccessKey          string `json:"access_key" binding:"required"`
	SecretKey          string `json:"secret_key" binding:"required"`
	SyncFrequencyHours int    `json:"sync_frequency_hours" binding:"required"`
	LookbackDays       int    `json:"lookback_days" binding:"required"`
}

type TestTenableConnectionDTO struct {
	BaseURL   string `json:"base_url"`
	AccessKey string `json:"access_key" binding:"required"`
	SecretKey string `json:"secret_key" binding:"required"`
}

type TestConnectionResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type TenableTestErrorResponse struct {
	ErrorMsg string `json:"error_msg"`
}
