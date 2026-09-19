package models

type ProjectPhase struct {
	ID        uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn   string `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr   string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	Type      string `json:"type" binding:"required,max=50,oneof=initiation planning execution go_live closure" gorm:"type:varchar(50)"`
	ColorCode string `json:"color_code" gorm:"type:varchar(25)"`
	Status    string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

// TableName overrides the table name used by ProjectPhase to `project_phases`
func (ProjectPhase) TableName() string {
	return "project_phases"
}