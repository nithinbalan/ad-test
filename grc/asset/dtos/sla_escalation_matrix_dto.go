package dtos


type SLAEscalationMatrixListDTO struct {
	ID              uint      `json:"id"`

	SeverityID      uint      `json:"severity_id"`
	SeverityTitleEn string    `json:"severity_title_en"`
	SeverityTitleAr string    `json:"severity_title_ar"`

	RemediationDays int       `json:"remediation_days"`

	IsEmail         bool      `json:"is_email"`
	IsAlert         bool      `json:"is_alert"`

	RoleID          uint      `json:"role_id"`
	RoleTitleEn     string    `json:"role_title_en"`
	RoleTitleAr     string    `json:"role_title_ar"`

	Status          string    `json:"status"`
}