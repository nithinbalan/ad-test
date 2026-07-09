package dtos

import "time"

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

	CreatedByID            uint      `json:"created_by_id"`
	CreatedByFirstNameEn   string    `json:"created_by_first_name_en"`
	CreatedByFirstNameAr   string    `json:"created_by_first_name_ar"`
	CreatedByLastNameEn    string    `json:"created_by_last_name_en"`
	CreatedByLastNameAr    string    `json:"created_by_last_name_ar"`
	CreatedByEmail         string    `json:"created_by_email"`
	CreatedByInitial       string    `json:"created_by_initial"`
	CreatedByColor         string    `json:"created_by_color"`
	CreatedByProfileImage  string    `json:"created_by_profile_image_path"`

	CreatedAt time.Time `json:"created_at"`
}