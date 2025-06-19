package dtos

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
	RiskTitle    string `form:"risk_title" binding:"required"`
	Consequences string `form:"consequences" binding:"required"`
	Category     string `form:"category" binding:"required"`
	Source       string `form:"source" binding:"required"`
}

type RiskAIRequest struct {
	RiskId        uint    `json:"risk_id" binding:"required"`
	RegisterTitle string  `json:"register_title" binding:"required"`
	Description   *string `json:"description"`
	Category      string  `json:"category" binding:"required"`
	Source        string  `json:"source"`
}
