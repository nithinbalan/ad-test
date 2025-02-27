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

type LicenseKeyUpdateFormDTO struct {
	LicenseKey *string `json:"license_key" binding:"required"`
}

type SystemPreferenceDTO struct {
	ID                         int         `json:"id"`
	CurrencyID                 uint        `json:"currency_id"`
	TimezoneID                 uint        `json:"timezone_id"`
	DateFormat                 string      `json:"date_format"`
	TimeFormat                 string      `json:"time_format"`
	DocUploadSize              int         `json:"doc_upload_size"`
	LogoUploadSize             int         `json:"logo_upload_size"`
	IdleTimeout                int         `json:"idle_timeout"`
	IsMaintenanceModeOn        *bool       `json:"is_maintenance_mode_on"`
	LicenseKey                 *string     `json:"license_key"`
	LicenseKeyExpiryDateString *string     `json:"license_key_expiry_date"`
	BalanceLicensedUserCount   *int        `json:"balance_licensed_user_count"`
	TotalLicensedUserCount     *int        `json:"total_licensed_user_count"`
	CreatedByID                uint        `json:"created_by_id"`
	UpdatedByID                *uint       `json:"updated_by_id"`
	Timezone                   TimezoneDTO `json:"timezone"`
	Currency                   CurrencyDTO `json:"currency"`
}

type TimezoneDTO struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Offset      string `json:"offset"`
	DisplayName string `json:"display_name"`
}

type CurrencyDTO struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Code  string `json:"code"`
	Icon  string `json:"icon"`
}

type OrganizationPreferenceResponseDTO struct {
	SystemPreference SystemPreferenceDTO   `json:"system_preference"`
	Integrations     []*IntegrationListDTO `json:"integrations"`
}
