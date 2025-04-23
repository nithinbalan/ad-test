package models

type ProcessCriticalityLevel struct {
	ID        uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn   string `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr   string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	Type      string `json:"type" binding:"required,max=20,oneof=low medium high" gorm:"type:varchar(20);not null"`
	ColorCode string `json:"color_code" binding:"required,max=25" gorm:"type:varchar(25);not null"`
	Status    string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

// TableName overrides the table name used by ProcessCriticalityLevel to `process_criticality_level`
func (ProcessCriticalityLevel) TableName() string {
	return "process_criticality_levels"
}
