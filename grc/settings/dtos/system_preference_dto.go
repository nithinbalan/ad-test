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
	ThemeID             uint   `json:"theme_id" binding:"required"`
}

type LicenseKeyUpdateFormDTO struct {
	LicenseKey *string `json:"license_key" binding:"required"`
}

type SystemPreferenceDTO struct {
	ID                  int            `json:"id"`
	CurrencyID          uint           `json:"currency_id"`
	TimezoneID          uint           `json:"timezone_id"`
	DateFormat          string         `json:"date_format"`
	TimeFormat          string         `json:"time_format"`
	DocUploadSize       int            `json:"doc_upload_size"`
	LogoUploadSize      int            `json:"logo_upload_size"`
	IdleTimeout         int            `json:"idle_timeout"`
	IsMaintenanceModeOn *bool          `json:"is_maintenance_mode_on"`
	IsAIEnable          *bool          `json:"is_ai_enable"`
	CreatedByID         uint           `json:"created_by_id"`
	UpdatedByID         *uint          `json:"updated_by_id"`
	Timezone            TimezoneDTO    `json:"timezone"`
	Currency            CurrencyDTO    `json:"currency"`
	LicenseKeyDetails   *LicenseKeyDTO `json:"license_key_details"`
	Theme               ThemeDTO       `json:"theme"`
}

type LicenseKeyDTO struct {
	LicenseKey                 *string  `json:"license_key"`
	LicenseKeyExpiryDateString *string  `json:"license_key_expiry_date"`
	BalanceLicensedUserCount   *int     `json:"balance_licensed_user_count"`
	TotalLicensedUserCount     *int     `json:"total_licensed_user_count"`
	EnabledModules             []string `json:"enabled_modules"`
	EnabledAddons              []string `json:"enabled_addons"`
	ActiveUsersCount           *int     `json:"active_user_count"`
	PassiveUsersCount          *int     `json:"passive_user_count"`
	BalanceActiveUsersCount    *int     `json:"balance_active_user_count"`
	BalancePassiveUsersCount   *int     `json:"balance_passive_user_count"`
}

type ThemeDTO struct {
	ID             uint       `json:"id"`
	TitleEn        string     `json:"title_en"`
	TitleAr        string     `json:"title_ar"`
	DescriptionEn  *string    `json:"description_en"`
	DescriptionAr  *string    `json:"description_ar"`
	PrimaryColor   string     `json:"primary_color"`
	SecondaryColor string     `json:"secondary_color"`
	LinkColor      string     `json:"link_color"`
	SidebarBgColor string     `json:"sidebar_bg_color"`
	IsSystem       bool       `json:"is_system"`
	CreatedByID    uint       `json:"created_by_id"`
	UpdatedByID    *uint      `json:"updated_by_id"`
	CreatedAt      string     `json:"created_at"`
	UpdatedAt      *string    `json:"updated_at"`
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
