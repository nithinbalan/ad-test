package dtos

type ModuleListQueryParamsDTO struct {
	SortBy            string `json:"sort_by"`
	Sort              string `json:"sort"`
	Keyword           string `json:"keyword,omitempty"`
	QuickFilter       []string
	IsActiveSubModule bool `json:"is_active_sub_module,omitempty"`
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
	Parent     *ModuleDetailDTO   `json:"parent"`
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

type ModuleReorderDTO struct {
	ID         uint   `json:"id"`
	Identifier string `json:"identifier"`
	Order      int    `json:"order"`
}

type ModuleFormDTO struct {
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Icon    string `json:"icon"`
}

type SubModuleFormDTO struct {
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
}
