package dtos

type SystemPreferenceFormDTO struct {
	CurrencyID            uint    `json:"currency_id" binding:"required"`
	DateFormat            *string `json:"date_format"`
	TimeFormat            string  `json:"time_format" binding:"required"`
	DocUploadSize         int     `json:"doc_upload_size" binding:"required"`
	LogoUploadSize        int     `json:"logo_upload_size" binding:"required"`
	IdleTimeout           int     `json:"idle_timeout" binding:"required"`
	IsMaintanceModeOn     *bool   `json:"is_maintance_mode_on"`
}
