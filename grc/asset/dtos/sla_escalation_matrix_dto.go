package dtos

type SLAEscalationMatrixListDTO struct {
	ID uint `json:"id"`

	SeverityID             uint   `json:"severity_id"`
	SeverityLevelTitleEn   string `json:"severity_level_title_en"`
	SeverityLevelTitleAr   string `json:"severity_level_title_ar"`
	SeverityLevelType      string `json:"severity_level_type"`
	SeverityLevelColorCode string `json:"severity_level_color_code"`
	SeverityLevelWeight    int    `json:"severity_level_weight"`

	RemediationDays int `json:"remediation_days"`

	IsEmail bool `json:"is_email"`
	IsAlert bool `json:"is_alert"`

	RoleID      uint   `json:"role_id"`
	RoleTitleEn string `json:"role_title_en"`
	RoleTitleAr string `json:"role_title_ar"`

	Status string `json:"status"`
}
