package dtos

import "time"

type RiskAutofillRequest struct {
	RiskTitle      string            `json:"risk_title"`
	RiskCategories []RiskCategoryDTO `json:"risk_categories"`
	RiskSources    []RiskSourceDTO   `json:"risk_sources"`
}

type RiskCategoryDTO struct {
	CategoryID uint   `json:"category_id"`
	Category   string `json:"category"`
}

type RiskSourceDTO struct {
	SourceID uint   `json:"source_id"`
	Source   string `json:"source"`
}

type RCAItem struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type RiskRCARequest struct {
	RiskTitle    string `json:"risk_title" binding:"required"`
	Consequences string `json:"consequences"`
	Category     string `json:"category" binding:"required"`
	Source       string `json:"source"`
}

type RiskAIRequest struct {
	RiskId        uint    `json:"risk_id" binding:"required"`
	RegisterTitle string  `json:"register_title" binding:"required"`
	Description   *string `json:"description"`
	Category      string  `json:"category" binding:"required"`
	Source        string  `json:"source"`
}

type RiskAIChatRequest struct {
	Content   string `json:"content" binding:"required"`
	Role      string `json:"role" binding:"required"`
	ContextID string `json:"context_id"`
}

type RiskAIChatDTO struct {
	ID        uint      `json:"id"`
	ContextID string    `json:"context_id"`
	Message   string    `json:"message"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
