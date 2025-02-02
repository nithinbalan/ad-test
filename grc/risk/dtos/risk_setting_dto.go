package dtos

import "time"

type RiskSettingDTO struct {
	ID            uint            `json:"id"`
	ReviewRoleID  uint            `json:"review_role_id"`
	ReviewRole    RoleDTO         `json:"review_role"`
	CreatedByID   uint            `json:"created_by_id"`
	CreatedByUser UserMinimalDTO  `json:"created_by"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedByUser *UserMinimalDTO `json:"updated_by"`
	UpdatedByID   *uint           `json:"updated_by_id"`
	UpdatedAt     *time.Time      `json:"updated_at"`
}

type RoleDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Status  string `json:"status"`
}
