package dtos

type SystemPreferenceFormDTO struct {
	CurrencyID          uint   `json:"currency_id" binding:"required"`
	TimezoneID          uint   `json:"timezone_id" binding:"required"`
	DateFormat          string `json:"date_format" binding:"required"`
	TimeFormat          string `json:"time_format" binding:"required"`
	DocUploadSize       int    `json:"doc_upload_size" binding:"required"`
	LogoUploadSize      int    `json:"logo_upload_size" binding:"required"`
	IdleTimeout         int    `json:"idle_timeout" binding:"required"`
	IsMaintenanceModeOn *bool  `json:"is_maintenance_mode_on"`
}

type SystemPreferenceDTO struct {
	CurrencyID          uint   `json:"currency_id" binding:"required"`
	TimezoneID          uint   `json:"timezone_id" binding:"required"`
	DateFormat          string `json:"date_format" binding:"required"`
	TimeFormat          string `json:"time_format" binding:"required"`
	DocUploadSize       int    `json:"doc_upload_size" binding:"required"`
	LogoUploadSize      int    `json:"logo_upload_size" binding:"required"`
	IdleTimeout         int    `json:"idle_timeout" binding:"required"`
	IsMaintenanceModeOn *bool  `json:"is_maintenance_mode_on"`
}

type OrganizationPreferenceResponseDTO struct {
	SystemPreference SystemPreferenceDTO   `json:"system_preference"`
	Integrations     []*IntegrationListDTO `json:"integrations"`
}
