package models

type WorkflowLevelType struct {
	ID      uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Type    string `json:"type" binding:"required" gorm:"type:varchar(20);not null;check:type IN ('individual', 'multiple-users', 'role', 'hod')"`
	TitleEn string `json:"title_en" binding:"required" gorm:"type:varchar(45);not null"`
	TitleAr string `json:"title_ar" binding:"required" gorm:"type:varchar(45);not null"`
	Status  string `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(20);not null;check:status IN ('active', 'inactive')"`
}

func (WorkflowLevelType) TableName() string {
	return "workflow_level_types"
}
