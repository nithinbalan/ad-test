package dtos

import (
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type AssetRiskResponseFormDTO struct {
	RiskResponseStrategyID uint `json:"risk_response_strategy_id" binding:"required"`
}

type AssetRiskResponseStrategyDTO struct {
	ID            uint   `json:"id"`
	TitleEn       string `json:"title_en"`
	TitleAr       string `json:"title_ar"`
	DescriptionEn string `json:"description_en"`
	DescriptionAr string `json:"description_ar"`
	Status        string `json:"status"`
}

type AssetRiskResponseDTO struct {
	ID                   uint                             `json:"id"`
	RiskRegisterID       uint                             `json:"risk_register_id"`
	RiskResponseStrategy *AssetRiskResponseStrategyDTO    `json:"risk_response_strategy"`
	CreatedByID          *uint                            `json:"created_by_id"`
	CreatedBy            *organizationDTOs.UserMinimalDTO `json:"created_by"`
	UpdatedByID          *uint                            `json:"updated_by_id"`
	UpdatedBy            *organizationDTOs.UserMinimalDTO `json:"updated_by"`
	CreatedAt            time.Time                        `json:"created_at"`
	UpdatedAt            *time.Time                       `json:"updated_at"`
}
