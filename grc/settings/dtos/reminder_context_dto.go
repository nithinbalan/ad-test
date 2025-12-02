package dtos

type ReminderContextDTO struct {
	ID       uint   `json:"id"`
	ModuleID uint   `json:"module_id"`
	Type     string `json:"type"`
	TitleEn  string `json:"title_en"`
	TitleAr  string `json:"title_ar"`
	Status   string `json:"status"`
}

type ReminderContextListQueryParamsDTO struct {
	Sort    string              `form:"sort" json:"sort"`
	SortBy  string              `form:"sortBy" json:"sortBy"`
	Keyword string              `form:"keyword" json:"keyword"`
	Status  string              `form:"status" json:"status"`
	Filters map[string][]string `form:"filters" json:"filters"`
}
