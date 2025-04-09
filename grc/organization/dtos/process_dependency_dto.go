package dtos

type ProcessDependencyDTO struct {
	ID                         uint    `json:"id"`
	ProcessID                  uint    `json:"process_id"`
	ProcessDependencySubTypeID uint    `json:"process_dependency_sub_type_id"`
	DependentProcessID         *uint   `json:"dependent_process_id,omitempty"`
	OrganizationUnitID         *uint   `json:"organization_unit_id,omitempty"`
	AssetID                    *uint   `json:"asset_id,omitempty"`
	Title                      string  `json:"title"`
	Description                *string `json:"description"`
	ProcessTitleEn             *string `json:"process_title_en,omitempty"`
	ProcessDependencySubTypeEn *string `json:"process_dependency_sub_type_title_en,omitempty"`
	ProcessDependencySubTypeAr *string `json:"process_dependency_sub_type_title_ar,omitempty"`
}

type ProcessDependencyFormDTO struct {
	ProcessDependencySubTypeID uint    `json:"process_dependency_sub_type_id" binding:"required"`
	DependentProcessID         *uint   `json:"dependent_process_id"`
	OrganizationUnitID         *uint   `json:"organization_unit_id"`
	AssetID                    *uint   `json:"asset_id"`
	Title                      string  `json:"title" binding:"required,max=255"`
	Description                *string `json:"description"`
}

type ProcessDependencyListQueryParamsDTO struct {
	Sort    string              `json:"sort"`
	SortBy  string              `json:"sort_by"`
	Keyword string              `json:"keyword"`
	Filters map[string][]string `json:"filters"`
}
