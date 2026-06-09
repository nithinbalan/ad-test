package models

type WorkflowLevelType struct {
	ID               uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Type             string `json:"type" binding:"required" gorm:"type:varchar(20);not null;check:type IN ('individual', 'multiple-users', 'role', 'hod', 'custom-roles')"`
	TitleEn          string `json:"title_en" binding:"required" gorm:"type:varchar(45);not null"`
	TitleAr          string `json:"title_ar" binding:"required" gorm:"type:varchar(45);not null"`
	Status           string `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(20);not null;check:status IN ('active', 'inactive')"`
	ModuleIdentifier string `json:"module_identifier" gorm:"type:nvarchar(255);not null;default:'common'"`
}

func (WorkflowLevelType) TableName() string {
	return "workflow_level_types"
}
