package dtos

type ProcessDependencyListDTO struct {
	ID             uint    `json:"id"`
	ProcessID      uint    `json:"process_id"`
	ProcessTitleEn *string `json:"process_title_en"`
	ProcessTitleAr *string `json:"process_title_ar"`

	ProcessDependencySubTypeID      uint    `json:"process_dependency_sub_type_id"`
	ProcessDependencySubTypeTitleEn *string `json:"process_dependency_sub_type_title_en"`
	ProcessDependencySubTypeTitleAr *string `json:"process_dependency_sub_type_title_ar"`
	ProcessDependencyTypeTitleEn    *string `json:"process_dependency_type_title_en"`
	ProcessDependencyTypeTitleAr    *string `json:"process_dependency_type_title_ar"`

	DependentProcessID      *uint   `json:"dependent_process_id"`
	DependentProcessTitleEn *string `json:"dependent_process_title_en"`
	DependentProcessTitleAr *string `json:"dependent_process_title_ar"`

	OrganizationUnitID      *uint   `json:"organization_unit_id"`
	OrganizationUnitTitleEn *string `json:"organization_unit_title_en"`
	OrganizationUnitTitleAr *string `json:"organization_unit_title_ar"`

	AssetID      *uint   `json:"asset_id"`
	AssetTitleEn *string `json:"asset_title_en"`
	AssetTitleAr *string `json:"asset_title_ar"`

	Title       string  `json:"title"`
	Description *string `json:"description"`
}

type ProcessDependencyDetailDTO struct {
	ID                         uint                        `json:"id"`
	ProcessID                  uint                        `json:"process_id"`
	ProcessDependencySubTypeID uint                        `json:"process_dependency_sub_type_id"`
	ProcessDependencySubType   ProcessDependencySubTypeDTO `json:"process_dependency_sub_type"`
	DependentProcessID         *uint                       `json:"dependent_process_id"`
	DependentProcess           *ProcessDetailDTO           `json:"dependent_process"`
	OrganizationUnitID         *uint                       `json:"organization_unit_id"`
	OrganizationUnit           *OrganizationUnit1DTO       `json:"organization_unit"`
	AssetID                    *uint                       `json:"asset_id"`
	Asset                      *AssetMinimalDTO            `json:"asset"`
	Title                      string                      `json:"title"`
	Description                *string                     `json:"description"`
}

type ProcessDependencyFormDTO struct {
	ProcessDependencySubTypeID uint    `json:"process_dependency_sub_type_id" binding:"required"`
	DependentProcessID         *uint   `json:"dependent_process_id"`
	OrganizationUnitID         *uint   `json:"organization_unit_id"`
	AssetID                    *uint   `json:"asset_id"`
	Title                      *string `json:"title" binding:"max=255"`
	Description                *string `json:"description"`
}

type ProcessDependencyListQueryParamsDTO struct {
	Sort    string              `json:"sort"`
	SortBy  string              `json:"sort_by"`
	Keyword string              `json:"keyword"`
	Filters map[string][]string `json:"filters"`
}
