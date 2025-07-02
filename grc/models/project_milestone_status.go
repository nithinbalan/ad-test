package models

type ProjectMilestoneStatus struct {
	ID            uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn       string `json:"title_en" binding:"required,max=45" gorm:"type:varchar(45);not null"`
	TitleAr       string `json:"title_ar" binding:"required,max=45" gorm:"type:nvarchar(45);not null"`
	DescriptionEn string `json:"description_en" gorm:"type:varchar(max)"`
	DescriptionAr string `json:"description_ar" gorm:"type:nvarchar(max)"`
	ColorCode     string `json:"color_code" binding:"required,max=25" gorm:"type:varchar(25);not null"`
	Type          string `json:"type" binding:"required,max=10,oneof=active archived expired" gorm:"type:varchar(10);not null"`
	Status        string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

// TableName overrides the table name used by ProjectMilestoneStatus to `project_milestone_statuses`
func (ProjectMilestoneStatus) TableName() string {
	return "project_milestone_statuses"
}
