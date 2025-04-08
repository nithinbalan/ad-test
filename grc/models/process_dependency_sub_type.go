package models

type ProcessDependencySubType struct {
	ID                      uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ProcessDependencyTypeID uint   `gorm:"column:process_dependency_type_id" json:"process_dependency_type_id"`
	TitleEn                 string `gorm:"type:nvarchar(255);not null" json:"title_en"`
	TitleAr                 string `gorm:"type:nvarchar(255);not null" json:"title_ar"`
	DescriptionEn           string `gorm:"type:nvarchar(MAX)" json:"description_en"`
	DescriptionAr           string `gorm:"type:nvarchar(MAX)" json:"description_ar"`
	Status                  string `gorm:"type:varchar(10);default:'active'" json:"status"`

	ProcessDependencyType *ProcessDependencyType `gorm:"foreignKey:ProcessDependencyTypeID" json:"process_dependency_type,omitempty"`
}

func (ProcessDependencySubType) TableName() string {
	return "process_dependency_sub_types"
}
