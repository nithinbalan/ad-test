package models

type ProcessDependency struct {
	ID                         uint    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	ProcessID                  uint    `gorm:"not null" json:"process_id"`
	ProcessDependencySubTypeID uint    `gorm:"not null" json:"process_dependency_sub_type_id"`
	DependentProcessID         *uint   `json:"dependent_process_id"`
	OrganizationUnitID         *uint   `json:"organization_unit_id"`
	AssetID                    *uint   `json:"asset_id"`
	ThirdPartyRegisterID       *uint   `json:"third_party_register_id"`
	Title                      *string `gorm:"type:nvarchar(255)" json:"title"`
	Description                *string `gorm:"type:nvarchar(MAX)" json:"description"`
	Rto                        *int    `json:"rto" gorm:"default:null"`
	RtoUnitId                  *uint   `json:"rto_uint_id" gorm:"default:null"`
	Mtpd                       *int    `json:"mtpd" gorm:"default:null"`
	MtpdUnitId                 *uint   `json:"mtpd_uint_id" gorm:"default:null"`
	Rpo                        *int    `json:"rpo" gorm:"default:null"`
	RpoUnitId                  *uint   `json:"rpo_uint_id" gorm:"default:null"`
	CreatedBYID                uint    `gorm:"not null" json:"created_by_id"`
	CreatedAt                  string  `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedBYID                *uint   `json:"updated_by_id"`
	UpdatedAt                  *string `gorm:"type:datetime" json:"updated_at"`

	// Relationships
	Process                  *Process                  `gorm:"foreignKey:ProcessID"`
	ProcessDependencySubType *ProcessDependencySubType `gorm:"foreignKey:ProcessDependencySubTypeID"`
	DependentProcess         *Process                  `gorm:"foreignKey:DependentProcessID"`
	OrganizationUnit         *OrganizationUnit         `gorm:"foreignKey:OrganizationUnitID"`
	Asset                    *AssetRegister            `gorm:"foreignKey:AssetID"`
	ThirdPartyRegister       *ThirdPartyRegister       `gorm:"foreignKey:ThirdPartyRegisterID"`
	RtoUnit                  *TimeUnit                 `gorm:"foreignKey:RtoUnitId"`
	MtpdUnit                 *TimeUnit                 `gorm:"foreignKey:MtpdUnitId"`
	RpoUnit                  *TimeUnit                 `gorm:"foreignKey:RpoUnitId"`
	CreatedByUser            *User                     `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser            *User                     `gorm:"foreignKey:UpdatedBYID"`
}

func (ProcessDependency) TableName() string {
	return "process_dependencies"
}
