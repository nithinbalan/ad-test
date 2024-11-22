package dtos

type ModuleListQueryParamsDTO struct {
	SortBy  string `json:"sort_by"`
	Sort    string `json:"sort"`
	Keyword string `json:"keyword,omitempty"`
}

type ModuleListDTO struct {
	ID         uint               `json:"id"`
	Identifier string             `json:"identifier"`
	TitleEn    string             `json:"title_en"`
	TitleAr    string             `json:"title_ar"`
	Order      int                `json:"order"`
	Icon       string             `json:"icon"`
	Route      string             `json:"route"`
	IsEnable   bool               `json:"is_enable"`
	SubModules []SubModuleListDTO `json:"sub_modules"`
}

type ModuleDetailDTO struct {
	ID         uint               `json:"id"`
	Identifier string             `json:"identifier"`
	TitleEn    string             `json:"title_en"`
	TitleAr    string             `json:"title_ar"`
	Order      int                `json:"order"`
	Icon       string             `json:"icon"`
	Route      string             `json:"route"`
	IsEnable   bool               `json:"is_enable"`
	SubModules []SubModuleListDTO `json:"sub_modules"`
}

type SubModuleListDTO struct {
	ID          uint            `json:"id"`
	Identifier  string          `json:"identifier"`
	TitleEn     string          `json:"title_en"`
	TitleAr     string          `json:"title_ar"`
	Order       int             `json:"order"`
	Icon        string          `json:"icon"`
	Route       string          `json:"route"`
	IsEnable    bool            `json:"is_enable"`
	Permissions map[string]bool `json:"permissions"`
}
