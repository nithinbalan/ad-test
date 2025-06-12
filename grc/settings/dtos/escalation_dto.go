package dtos

import "time"

type EscalationDTO struct {
	ID       uint   `json:"id"`
	TitleEn  string `json:"title_en"`
	TitleAr  string `json:"title_ar"`
	ModuleID uint   `json:"module_id"`
	Status   string `json:"status"`
}

type CreateEscalationDTO struct {
	TitleEn  string `json:"title_en" binding:"required"`
	TitleAr  string `json:"title_ar" binding:"required"`
	ModuleID uint   `json:"module_id" binding:"required"`
}

type UpdateEscalationDTO struct {
	TitleEn string `json:"title_en" binding:"required"`
	TitleAr string `json:"title_ar" binding:"required"`
	Status  string `json:"status"`
}

type EscalationDetailsDTO struct {
	ID               uint                       `json:"id"`
	TitleEn          string                     `json:"title_en"`
	TitleAr          string                     `json:"title_ar"`
	ModuleID         uint                       `json:"module_id"`
	Status           string                     `json:"status"`
	CreatedAt        time.Time                  `json:"created_at"`
	UpdatedAt        time.Time                  `json:"updated_at"`
	EscalationLevels []EscalationLevelDetailDTO `json:"escalation_levels"`
}

type EscalationListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
	Status  string
	Filters map[string][]string
}

type GroupedEscalationDTO struct {
	ID               uint                       `json:"id"`
	ModuleID         uint                       `json:"module_id"`
	TitleEn          string                     `json:"title_en"`
	TitleAr          string                     `json:"title_ar"`
	Description      string                     `json:"description"`
	Status           string                     `json:"status"`
	CreatedAt        time.Time                  `json:"created_at"`
	UpdatedAt        time.Time                  `json:"updated_at"`
	EscalationLevels []EscalationLevelDetailDTO `json:"escalation_levels"`
}
