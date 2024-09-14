package models

type RolePermissionType struct {
	ID      uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Type    string `json:"type" binding:"required" gorm:"type:varchar(20);not null;check:type IN ('manage', 'modify', 'view')"`
	TitleEn string `json:"title_en" binding:"required" gorm:"type:varchar(100);not null"`
	TitleAr string `json:"title_ar" binding:"required" gorm:"type:varchar(100);not null"`
	Status  string `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(20);not null;check:status IN ('active', 'inactive')"`
}

func (RolePermissionType) TableName() string {
	return "workflow_level_types"
}
