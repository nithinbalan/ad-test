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

type UpdateTenableConnectionDTO struct {
	BaseURL   string `json:"base_url" binding:"required"`
	AccessKey string `json:"access_key" binding:"required"`
	SecretKey string `json:"secret_key" binding:"required"`
}

type UpdateTenableSyncSettingsDTO struct {
	SyncFrequencyHours int `json:"sync_frequency_hours" binding:"required,min=1"`
	LookbackDays       int `json:"lookback_days" binding:"required,min=1"`
}

type TestTenableConnectionDTO struct {
	BaseURL   string `json:"base_url"`
	AccessKey string `json:"access_key" binding:"required"`
	SecretKey string `json:"secret_key" binding:"required"`
}
