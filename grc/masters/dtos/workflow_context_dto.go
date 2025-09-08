package dtos

type WorkflowContextDTO struct {
	ID       uint   `json:"id"`
	ModuleID uint   `json:"module_id"`
	TitleEn  string `json:"title_en"`
	TitleAr  string `json:"title_ar"`
	Type     string `json:"type"`
}

type WorkflowContextListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
	Status  string
	Filters map[string][]string
}
