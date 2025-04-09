package models

type ProcessDependency struct {
	ID                         uint    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	ProcessID                  uint    `gorm:"not null" json:"process_id"`
	ProcessDependencySubTypeID uint    `gorm:"not null" json:"process_dependency_sub_type_id"`
	DependentProcessID         *uint   `json:"dependent_process_id"`
	OrganizationUnitID         *uint   `json:"organization_unit_id"`
	AssetID                    *uint   `json:"asset_id"`
	Title                      string  `gorm:"type:nvarchar(255);not null" json:"title"`
	Description                string  `gorm:"type:nvarchar(MAX)" json:"description"`
	CreatedBYID                uint    `gorm:"not null" json:"created_by_id"`
	CreatedAt                  string  `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedBYID                *uint   `json:"updated_by_id"`
	UpdatedAt                  *string `gorm:"type:datetime" json:"updated_at"`

	// Relationships
	Process           *Process                  `gorm:"foreignKey:ProcessID" json:"process,omitempty"`
	DependencySubType *ProcessDependencySubType `gorm:"foreignKey:ProcessDependencySubTypeID" json:"process_dependency_sub_type,omitempty"`
	DependentProcess  *Process                  `gorm:"foreignKey:DependentProcessID" json:"dependent_process,omitempty"`
	OrganizationUnit  *OrganizationUnit         `gorm:"foreignKey:OrganizationUnitID" json:"organization_unit,omitempty"`
	Asset             *AssetRegister            `gorm:"foreignKey:AssetID" json:"asset,omitempty"`
	CreatedByUser     *User                     `gorm:"foreignKey:CreatedBYID" json:"created_by,omitempty"`
	UpdatedByUser     *User                     `gorm:"foreignKey:UpdatedBYID" json:"updated_by,omitempty"`
}

func (ProcessDependency) TableName() string {
	return "process_dependencies"
}
