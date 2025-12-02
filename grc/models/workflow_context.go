package models

type WorkflowContext struct {
	ID       uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ModuleID uint   `json:"module_id"`
	TitleEn  string `json:"title_en" binding:"required,max=255" gorm:"type:nvarchar(255);not null"`
	TitleAr  string `json:"title_ar" binding:"required,max=255" gorm:"type:nvarchar(255);not null"`
	Type     string `json:"type" binding:"required" gorm:"type:varchar(10);not null"`
}

// TableName overrides the default table name
func (WorkflowContext) TableName() string {
	return "workflow_contexts"
}
