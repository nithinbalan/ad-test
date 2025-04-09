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
