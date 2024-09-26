package models

type WorkflowStatus struct {
	ID        uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn   string `json:"title_en" binding:"required" gorm:"type:nvarchar(100);not null"`
	TitleAr   string `json:"title_ar" binding:"required" gorm:"type:nvarchar(100);not null"`
	ColorCode string `json:"color_code" binding:"required" gorm:"type:varchar(25);not null"`
	Type      string `json:"type" binding:"required,oneof=submit revert reject approve publish" gorm:"type:varchar(10);not null;check:type IN ('submit', 'revert', 'reject', 'approve', 'publish')"`
	Status    string `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);not null;default:'active';check:status IN ('active', 'inactive')"`
}

// TableName overrides the default table name used by GORM
func (WorkflowStatus) TableName() string {
	return "workflow_statuses"
}
