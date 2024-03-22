package models

type LifecycleStatus struct {
	ID        uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn   string `json:"title_en" binding:"required,max=45" gorm:"type:varchar(45);not null"`
	TitleAr   string `json:"title_ar" binding:"required,max=45" gorm:"type:nvarchar(45);not null"`
	ColorCode string `json:"color_code" binding:"required,max=25" gorm:"type:varchar(25);not null"`
	Type      string `json:"type" binding:"required,max=10,oneof=active archived expired" gorm:"type:varchar(10);not null"`
	Status    string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

// TableName overrides the table name used by LifecycleStatus to `lifecycle_statuses`
func (LifecycleStatus) TableName() string {
	return "lifecycle_statuses"
}
