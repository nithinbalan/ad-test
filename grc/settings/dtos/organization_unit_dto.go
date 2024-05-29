package dtos

type OrganizationUnitDTO struct {
	ID         int                  `json:"id"`
	TitleEn    string               `json:"title_en"`
	TitleAr    string               `json:"title_ar"`
	Status     string               `json:"status"`
	ColorCode  string               `json:"color_code"`
	ParentUnit *OrganizationUnitDTO `json:"parent_unit"`
}

type OrganizationUnitListDTO struct {
	ID                            uint   `json:"id"`
	OrganizationLevelID           string `json:"organization_level_id"`
	OrganizationLevelTitleEn      string `json:"organization_level_title_en"`
	OrganizationLevelTitleAr      string `json:"organization_level_title_ar"`
	ParentOrganizationUnitID      string `json:"parent_organization_unit_id"`
	ParentOrganizationUnitTitleEn string `json:"parent_organization_unit_title_en"`
	ParentOrganizationUnitTitleAr string `json:"parent_organization_unit_title_ar"`
	ColorCode                     string `json:"color_code"`
	TitleEn                       string `json:"title_en"`
	TitleAr                       string `json:"title_ar"`
	Status                        string `json:"status"`
}

type OrganizationUnitListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
	Status  string
	Primary bool
	Filters map[string][]string
}

type OrganizationUnitFormDTO struct {
	OrganizationLevelID uint   `json:"organization_level_id" binding:"required"`
	OrganizationUnitID  *uint  `json:"organization_unit_id"`
	TitleEn             string `json:"title_en" binding:"required,max=100"`
	TitleAr             string `json:"title_ar" binding:"required,max=100"`
	Status              string `json:"status" binding:"required"`
}
