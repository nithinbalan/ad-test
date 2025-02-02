package dtos

import "time"

type RiskSettingDTO struct {
	ID           uint       `json:"id"`
	ReviewRoleID uint       `json:"review_role_id"`
	ReviewRole   RoleDTO    `json:"review_role,omitempty"`
	CreatedByID  uint       `json:"created_by_id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedByID  *uint      `json:"updated_by_id,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

type RoleDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Status  string `json:"status"`
}
