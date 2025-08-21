package dtos

type BcpRiskForm struct {
	RiskRegisterIDs []uint `json:"risk_register_ids" binding:"required"`
}

type BcpRiskListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Value       int
	QuickFilter []string
}
