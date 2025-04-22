package dtos

type StrategyProfileSettingListDTO struct {
	ID       int    `json:"id"`
	TitleEn  string `json:"title_en"`
	TitleAr  string `json:"title_ar"`
	IsEnable bool   `json:"is_enable"`
}

type StrategyProfileSettingListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}

type StrategyProfileSettingDetailDTO struct {
	ID       uint   `json:"id"`
	TitleEn  string `json:"title_en"`
	TitleAr  string `json:"title_ar"`
	IsEnable bool   `json:"is_enable"`
}
